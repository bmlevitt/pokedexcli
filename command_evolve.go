package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

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

	return nil
}

// findEvolutionsFor recursively searches the evolution chain for the specified Pokemon
// and returns its possible evolutions
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
