package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// commandEvolve evolves a Pokémon in the user's Pokédex into its next evolution.
// This command simulates the evolution mechanic from the Pokémon games by replacing
// the original Pokémon in the Pokédex with its evolved form.
//
// If a Pokémon has multiple possible evolutions, the user is prompted to choose
// which evolution they want. They can specify an evolution either by number or name
// as an additional parameter, or select from a menu if no choice is provided.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon to evolve and
//     params[1] (optional) is the evolution selection
//
// Returns:
//   - An error if no Pokémon name is provided, if the Pokémon is not in the Pokédex,
//     if the Pokémon cannot evolve further, if an invalid selection is made,
//     or if there's an issue with the API request
func commandEvolve(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	inputName := params[0]

	// Convert the input name to API format if it's in a formatted style
	apiPokemonName := ConvertToAPIFormat(inputName)
	formattedName := FormatPokemonName(apiPokemonName)

	// Check if the pokemon exists in the pokedex
	_, exists := cfg.pokedex[apiPokemonName]
	if !exists {
		// Check if it's a capitalization issue by trying all keys
		found := false
		for key := range cfg.pokedex {
			if ConvertToAPIFormat(key) == apiPokemonName {
				apiPokemonName = key
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%s is not in your Pokédex", formattedName)
		}
	}

	// Get the evolution chain for the Pokemon
	evolutionChain, err := cfg.pokeapiClient.GetEvolutionChainBySpecies(apiPokemonName)
	if err != nil {
		return fmt.Errorf("error getting evolution data: %v", err)
	}

	// Find the Pokemon in the evolution chain and its possible evolutions
	evolutions, err := findEvolutionsFor(apiPokemonName, evolutionChain.Chain)
	if err != nil {
		return err
	}

	// Check if the Pokemon can evolve
	if len(evolutions) == 0 {
		return fmt.Errorf("%s cannot evolve any further", formattedName)
	}

	// Handle multiple evolution options
	var chosenEvolution string
	if len(evolutions) == 1 {
		// Only one evolution option
		chosenEvolution = evolutions[0].Species.Name
	} else {
		// Multiple evolution options - let user choose
		fmt.Printf("%s can evolve into multiple forms:\n", formattedName)
		for i, evolution := range evolutions {
			formattedEvolution := FormatPokemonName(evolution.Species.Name)
			fmt.Printf(" %d. %s\n", i+1, formattedEvolution)
		}

		// If additional parameters were provided, they might specify which evolution
		if len(params) > 1 {
			// Check if the second parameter is a number (selection)
			selection, err := strconv.Atoi(params[1])
			if err == nil && selection > 0 && selection <= len(evolutions) {
				chosenEvolution = evolutions[selection-1].Species.Name
			} else {
				// Check if it matches an evolution name - convert to API format first
				evolutionInput := ConvertToAPIFormat(params[1])
				for _, evolution := range evolutions {
					if strings.EqualFold(evolution.Species.Name, evolutionInput) {
						chosenEvolution = evolution.Species.Name
						break
					}
				}
			}
		}

		// If no valid selection was made, prompt the user
		if chosenEvolution == "" {
			fmt.Print("Choose a form to evolve into (enter number): ")
			var selection int
			fmt.Scanf("%d", &selection)
			if selection > 0 && selection <= len(evolutions) {
				chosenEvolution = evolutions[selection-1].Species.Name
			} else {
				return errors.New("invalid selection")
			}
		}
	}

	// Evolve the Pokemon
	formattedEvolution := FormatPokemonName(chosenEvolution)
	fmt.Printf("Evolving %s into %s...\n", formattedName, formattedEvolution)

	// Get the evolution's data
	evolvedPokemonData, err := cfg.pokeapiClient.GetPokemonData(chosenEvolution)
	if err != nil {
		return fmt.Errorf("error fetching evolved Pokémon data: %v", err)
	}

	// Remove the original Pokemon from the pokedex
	delete(cfg.pokedex, apiPokemonName)

	// Add the evolved Pokemon to the pokedex
	cfg.pokedex[chosenEvolution] = evolvedPokemonData

	fmt.Printf("Congratulations! Your %s evolved into %s!\n",
		formattedName,
		formattedEvolution)
	fmt.Println("-----")

	// Auto-save after evolving a Pokémon
	cfg.changesSinceSync++
	if cfg.changesSinceSync >= cfg.autoSaveInterval {
		if err := autoSaveIfEnabled(cfg); err != nil {
			return fmt.Errorf("error auto-saving: %w", err)
		}
		cfg.changesSinceSync = 0
	}

	return nil
}

// findEvolutionsFor searches the evolution chain for a specific Pokémon
// and returns a list of its possible evolutions.
//
// Parameters:
//   - pokemonName: The name of the Pokémon to find evolutions for
//   - chainLink: The current node in the evolution chain tree
//
// Returns:
//   - A slice of chain links representing possible evolutions for the Pokémon
//   - An error if the Pokémon cannot be found in the evolution chain
func findEvolutionsFor(pokemonName string, chainLink pokeapi.ChainLink) ([]pokeapi.ChainLink, error) {
	// Check if this link is the Pokemon we're looking for
	if strings.EqualFold(chainLink.Species.Name, pokemonName) {
		return chainLink.EvolvesTo, nil
	}

	// Recursively check each child evolution
	for _, evolution := range chainLink.EvolvesTo {
		if found, err := findEvolutionsFor(pokemonName, evolution); err == nil {
			return found, nil
		}
	}

	return nil, fmt.Errorf("could not find %s in the evolution chain", pokemonName)
}
