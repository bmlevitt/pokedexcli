package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// commandFlavor retrieves and displays flavor text for a Pokémon in the user's Pokédex.
// Flavor text is descriptive content about the Pokémon from various game versions.
// The function retrieves the species data, filters for English text entries,
// and randomly selects one to display alongside the Pokémon's category (genus).
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon name
//
// Returns:
//   - An error if the Pokémon is not in the Pokédex, if no flavor text is available,
//     or if there's an issue with the API request
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
// and normalizes spacing for better readability in the terminal.
//
// Parameters:
//   - text: The raw flavor text string to clean
//
// Returns:
//   - A string with consistent spacing and no control characters
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
