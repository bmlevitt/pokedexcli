package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
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
	// Use the utility function to validate the Pokemon parameter and check if it exists
	apiName, nameInfo, _, _, err := GetPokemonIfExists(cfg, params)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "evolve", err) {
			return err
		}
		return nil
	}

	// Get the evolution chain for the Pokemon
	evolutionChain, err := cfg.pokeapiClient.GetEvolutionChainBySpecies(apiName)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "evolve", err) {
			return err
		}
		return nil
	}

	// Find the Pokemon in the evolution chain and its possible evolutions
	evolutions, err := findEvolutionsFor(apiName, evolutionChain.Chain)
	if err != nil {
		// Create a specific error for this case
		evolveErr := errorhandling.NewInvalidInputError(
			fmt.Sprintf("%s cannot evolve (not found in evolution chain)", nameInfo.Formatted), err)

		// Use standardized error handling
		if HandleCommandError(cfg, "evolve", evolveErr) {
			return evolveErr
		}
		return nil
	}

	if len(evolutions) == 0 {
		// Create a specific error for this case
		evolveErr := errorhandling.NewInvalidInputError(
			fmt.Sprintf("%s cannot evolve any further", nameInfo.Formatted), nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "evolve", evolveErr) {
			return evolveErr
		}
		return nil
	}

	// Handle evolution choice
	var selectedEvolution pokeapi.ChainLink
	if len(evolutions) == 1 {
		// Only one possible evolution
		selectedEvolution = evolutions[0]
	} else {
		// Multiple possible evolutions, need to choose
		if len(params) > 1 {
			// User provided a selection parameter
			selection := params[1]
			selectionNum, err := strconv.Atoi(selection)
			if err == nil && selectionNum > 0 && selectionNum <= len(evolutions) {
				// Valid numeric selection
				selectedEvolution = evolutions[selectionNum-1]
			} else {
				// Try to match by name
				found := false
				for _, evo := range evolutions {
					if strings.EqualFold(ConvertToAPIFormat(selection), evo.Species.Name) {
						selectedEvolution = evo
						found = true
						break
					}
				}
				if !found {
					// Show evolution options
					fmt.Printf("%s can evolve into multiple forms. Choose one:\n", nameInfo.Formatted)
					for i, evolution := range evolutions {
						formattedEvolution := FormatPokemonName(evolution.Species.Name)
						fmt.Printf("%d. %s\n", i+1, formattedEvolution)
					}
					return errorhandling.NewInvalidInputError(
						fmt.Sprintf("Invalid evolution selection: '%s'", selection), nil)
				}
			}
		} else {
			// No selection provided, show options
			fmt.Printf("%s can evolve into multiple forms. Choose one:\n", nameInfo.Formatted)
			for i, evolution := range evolutions {
				formattedEvolution := FormatPokemonName(evolution.Species.Name)
				fmt.Printf("%d. %s\n", i+1, formattedEvolution)
			}
			return errorhandling.NewInvalidInputError(
				fmt.Sprintf("Please specify which evolution to use (e.g., 'evolve %s 1')", nameInfo.APIFormat), nil)
		}
	}

	// Get data for the evolved form
	evolvedName := selectedEvolution.Species.Name
	evolvedFormattedName := FormatPokemonName(evolvedName)
	evolvedData, err := cfg.pokeapiClient.GetPokemonData(evolvedName)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "evolve", err) {
			return err
		}
		return nil
	}

	// Add evolved form to pokedex
	cfg.mutex.Lock()
	// First remove the original pokemon
	delete(cfg.pokedex, apiName)
	// Then add the evolved form
	cfg.pokedex[evolvedName] = evolvedData
	cfg.mutex.Unlock()

	fmt.Printf("Evolving %s into %s...\n", nameInfo.Formatted, evolvedFormattedName)
	fmt.Printf("Congratulations! Your %s evolved into %s!\n", nameInfo.Formatted, evolvedFormattedName)
	fmt.Println("-----")

	// Auto-save after evolving
	if err := UpdatePokedexAndSave(cfg); err != nil {
		// Use standardized error handling but don't return the error
		// since we still want to show the success message
		HandleCommandError(cfg, "evolve", err)
	}

	return nil
}

// findEvolutionsFor searches through an evolution chain to find the possible
// evolutions for a given Pokémon.
//
// Parameters:
//   - pokemonName: The name of the Pokémon to find evolutions for
//   - chainLink: The current link in the evolution chain to search
//
// Returns:
//   - A slice of ChainLink objects representing possible evolutions
//   - An error if the Pokémon cannot be found in the evolution chain
func findEvolutionsFor(pokemonName string, chainLink pokeapi.ChainLink) ([]pokeapi.ChainLink, error) {
	// Check if this is our pokemon
	if chainLink.Species.Name == pokemonName {
		return chainLink.EvolvesTo, nil
	}

	// Check if our pokemon is in the evolves_to array
	for _, evolution := range chainLink.EvolvesTo {
		if evolution.Species.Name == pokemonName {
			return evolution.EvolvesTo, nil
		}

		// Recursively check deeper in the chain
		result, err := findEvolutionsFor(pokemonName, evolution)
		if err == nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("pokemon not found in evolution chain")
}
