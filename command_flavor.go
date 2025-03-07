package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func commandFlavor(cfg *config, params []string) error {
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

	// Fetch species data which contains flavor text
	speciesData, err := cfg.pokeapiClient.GetPokemonSpecies(pokemonName)
	if err != nil {
		return fmt.Errorf("error fetching flavor text: %v", err)
	}

	// Filter flavor texts in English
	var englishFlavorTexts []struct {
		Text    string
		Version string
	}

	for _, entry := range speciesData.FlavorTextEntries {
		if entry.Language.Name == "en" {
			englishFlavorTexts = append(englishFlavorTexts, struct {
				Text    string
				Version string
			}{
				Text:    cleanFlavorText(entry.FlavorText),
				Version: entry.Version.Name,
			})
		}
	}

	if len(englishFlavorTexts) == 0 {
		return fmt.Errorf("no flavor text available for %s", pokemonName)
	}

	// Get a random flavor text
	randomIndex := rand.Intn(len(englishFlavorTexts))
	selectedFlavorText := englishFlavorTexts[randomIndex]

	// Get genus information (e.g., "Mouse Pokémon")
	genus := "Pokémon" // Default
	for _, genusEntry := range speciesData.Genera {
		if genusEntry.Language.Name == "en" {
			genus = genusEntry.Genus
			break
		}
	}

	// Print the flavor text with game version and genus
	fmt.Printf("%s, the %s\n\n", capitalizeFirstLetter(pokemonName), genus)
	fmt.Printf("%s\n\n", selectedFlavorText.Text)
	fmt.Printf("(From Pokémon %s)\n", capitalizeFirstLetter(selectedFlavorText.Version))

	return nil
}

// cleanFlavorText cleans up the flavor text by removing newlines and form feeds
func cleanFlavorText(text string) string {
	// Replace newlines and form feeds with spaces
	cleaned := strings.ReplaceAll(text, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\f", " ")

	// Replace multiple spaces with a single space
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	return cleaned
}
