package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// commandDescribe provides information about a Pokémon in the user's Pokédex.
// This enhanced command combines form descriptions and flavor text to give comprehensive
// information about the Pokémon.
//
// The command displays:
// 1. The Pokémon's name and genus (e.g., "Pikachu, the Mouse Pokémon")
// 2. A flavor text entry from one of the Pokémon games, with the game name in parentheses
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

	// Fetch species data which contains descriptions and flavor text
	speciesData, err := cfg.pokeapiClient.GetPokemonSpecies(pokemonName)
	if err != nil {
		return fmt.Errorf("error fetching description: %v", err)
	}

	// Get the Pokémon's genus (e.g., "Mouse Pokémon")
	var genus string
	for _, genusEntry := range speciesData.Genera {
		if genusEntry.Language.Name == "en" {
			genus = genusEntry.Genus
			break
		}
	}

	// Get flavor text entries in English
	var englishFlavorTexts []struct {
		text    string
		version string
	}

	for _, entry := range speciesData.FlavorTextEntries {
		if entry.Language.Name == "en" {
			cleanText := cleanText(entry.FlavorText)
			if cleanText != "" {
				englishFlavorTexts = append(englishFlavorTexts, struct {
					text    string
					version string
				}{
					text:    cleanText,
					version: entry.Version.Name,
				})
			}
		}
	}

	// Display the Pokémon info
	capitalizedName := cases.Title(language.English).String(strings.ToLower(pokemonName))
	fmt.Println("-----")
	if genus != "" {
		fmt.Printf("%s, the %s\n", capitalizedName, genus)
	} else {
		fmt.Printf("%s\n", capitalizedName)
	}

	// Display a random flavor text if available
	if len(englishFlavorTexts) > 0 {
		// Select a random flavor text
		selectedEntry := englishFlavorTexts[rand.Intn(len(englishFlavorTexts))]

		// Format the game name to look nicer
		gameName := strings.ReplaceAll(selectedEntry.version, "-", " ")
		gameName = cases.Title(language.English).String(strings.ToLower(gameName))

		// Display the flavor text with the requested format
		fmt.Printf("\"%s\"\n", selectedEntry.text)
		fmt.Printf("(From Pokémon %s)\n", gameName)
		fmt.Println("-----")
	} else {
		fmt.Println("- No description available.")
	}

	return nil
}

// cleanText removes unwanted characters and normalizes spaces in text.
func cleanText(text string) string {
	// Replace line breaks and form feeds with spaces
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\f", " ")

	// Replace any sequence of spaces with a single space
	text = strings.Join(strings.Fields(text), " ")

	return text
}
