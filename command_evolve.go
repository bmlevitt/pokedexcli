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
	pokemonName := params[0]

	// Check if the pokemon exists in the pokedex
	_, exists := cfg.pokedex[pokemonName]
	if !exists {
		return fmt.Errorf("%s is not in your pokedex", pokemonName)
	}

	// Get the evolution chain for the Pokemon
	evolutionChain, err := cfg.pokeapiClient.GetEvolutionChainBySpecies(pokemonName)
	if err != nil {
		return fmt.Errorf("error getting evolution data: %v", err)
	}

	// Find the Pokemon in the evolution chain and its possible evolutions
	evolutions, err := findEvolutionsFor(pokemonName, evolutionChain.Chain)
	if err != nil {
		return err
	}

	// Check if the Pokemon can evolve
	if len(evolutions) == 0 {
		return fmt.Errorf("%s cannot evolve any further", pokemonName)
	}

	// Handle multiple evolution options
	var chosenEvolution string
	if len(evolutions) == 1 {
		// Only one evolution option
		chosenEvolution = evolutions[0].Species.Name
	} else {
		// Multiple evolution options - let user choose
		fmt.Printf("%s can evolve into multiple forms:\n", pokemonName)
		for i, evolution := range evolutions {
			fmt.Printf(" %d. %s\n", i+1, capitalizeFirstLetter(evolution.Species.Name))
		}

		// If additional parameters were provided, they might specify which evolution
		if len(params) > 1 {
			// Check if the second parameter is a number (selection)
			selection, err := strconv.Atoi(params[1])
			if err == nil && selection > 0 && selection <= len(evolutions) {
				chosenEvolution = evolutions[selection-1].Species.Name
			} else {
				// Check if it matches an evolution name
				for _, evolution := range evolutions {
					if strings.EqualFold(evolution.Species.Name, params[1]) {
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
	fmt.Printf("Evolving %s into %s...\n", pokemonName, capitalizeFirstLetter(chosenEvolution))

	// Get the evolution's data
	evolvedPokemonData, err := cfg.pokeapiClient.GetPokemonData(chosenEvolution)
	if err != nil {
		return fmt.Errorf("error fetching evolved Pokemon data: %v", err)
	}

	// Remove the original Pokemon from the pokedex
	delete(cfg.pokedex, pokemonName)

	// Add the evolved Pokemon to the pokedex
	cfg.pokedex[chosenEvolution] = evolvedPokemonData

	fmt.Printf("Congratulations! Your %s evolved into %s!\n",
		capitalizeFirstLetter(pokemonName),
		capitalizeFirstLetter(chosenEvolution))

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

// findEvolutionsFor recursively searches the evolution chain for a specified Pokémon
// and returns its possible evolutions.
//
// This helper function traverses the evolution chain structure to find where the
// target Pokémon exists in the chain, then returns the next evolutions in the sequence.
//
// Parameters:
//   - pokemonName: The name of the Pokémon to find in the evolution chain
//   - chainLink: The current node in the evolution chain to search
//
// Returns:
//   - A slice of ChainLink objects representing possible evolutions
//   - An error if the Pokémon is not found in the evolution chain
func findEvolutionsFor(pokemonName string, chainLink pokeapi.ChainLink) ([]pokeapi.ChainLink, error) {
	// Check if this is the Pokemon we're looking for
	if strings.EqualFold(chainLink.Species.Name, pokemonName) {
		// If this is our Pokemon, return its possible evolutions
		return chainLink.EvolvesTo, nil
	}

	// If not, search in the evolution chain
	for _, evolution := range chainLink.EvolvesTo {
		result, err := findEvolutionsFor(pokemonName, evolution)
		if err == nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("couldn't find %s in the evolution chain", pokemonName)
}
