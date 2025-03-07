package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandDescribe(cfg *config, params []string) error {
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

	// Fetch species data which contains descriptions
	speciesData, err := cfg.pokeapiClient.GetPokemonSpecies(pokemonName)
	if err != nil {
		return fmt.Errorf("error fetching description: %v", err)
	}

	// Filter descriptions in English
	var englishDescriptions []string
	for _, desc := range speciesData.FormDescriptions {
		if desc.Language.Name == "en" {
			englishDescriptions = append(englishDescriptions, desc.Description)
		}
	}

	// If no form descriptions are available, just show the genus
	if len(englishDescriptions) == 0 {
		// Get genus information as a fallback
		genus := "Pokémon" // Default
		for _, genusEntry := range speciesData.Genera {
			if genusEntry.Language.Name == "en" {
				genus = genusEntry.Genus
				break
			}
		}

		fmt.Printf("%s is a %s.\n", capitalizeFirstLetter(pokemonName), genus)
		return nil
	}

	// Get a random description if multiple are available
	var selectedDescription string
	if len(englishDescriptions) > 1 {
		randomIndex := rand.Intn(len(englishDescriptions))
		selectedDescription = englishDescriptions[randomIndex]
	} else {
		selectedDescription = englishDescriptions[0]
	}

	// Print the description
	fmt.Printf("%s: %s\n", capitalizeFirstLetter(pokemonName), selectedDescription)

	return nil
}
