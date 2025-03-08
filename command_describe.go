package main

import (
	"errors"
	"fmt"
	"math/rand"
)

// commandDescribe displays form descriptions for a Pokémon in the user's Pokédex.
// This command retrieves species data which contains descriptive information about
// the Pokémon's different forms. It filters for English descriptions and randomly
// selects one to display.
//
// If no form descriptions are available in English, the function falls back to
// displaying the Pokémon's genus (e.g., "Mouse Pokémon").
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon name
//
// Returns:
//   - An error if no Pokémon name is provided, if the Pokémon is not in the Pokédex,
//     or if there's an issue with the API request
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
