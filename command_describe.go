package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// commandDescribe displays detailed Pokédex information about a Pokémon.
// This command shows flavor text entries (Pokédex descriptions) for a Pokémon,
// including its genus (e.g., "Mouse Pokémon") and a randomly selected
// description from the games.
//
// The command can only be used with Pokémon that are currently in the user's Pokédex.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon name to describe
//
// Returns:
//   - An error if no Pokémon name is provided, if the Pokémon is not in the Pokédex,
//     or if there's an issue with the API request
func commandDescribe(cfg *config, params []string) error {
	// Use the utility function to validate the Pokemon parameter and check if it exists
	apiName, nameInfo, _, _, err := GetPokemonIfExists(cfg, params)
	if err != nil {
		return err
	}

	// Fetch species data for the pokemon
	speciesData, err := cfg.pokeapiClient.GetPokemonSpecies(apiName)
	if err != nil {
		return fmt.Errorf("error fetching species data: %v", err)
	}

	// Find the English genus
	var genus string
	for _, genusEntry := range speciesData.Genera {
		if genusEntry.Language.Name == "en" {
			genus = genusEntry.Genus
			break
		}
	}

	// Find English flavor text entries
	var englishEntries []int
	for i, entry := range speciesData.FlavorTextEntries {
		if entry.Language.Name == "en" {
			englishEntries = append(englishEntries, i)
		}
	}

	// Display the information
	if len(englishEntries) > 0 {
		// Select a random entry index
		randomIndex := englishEntries[rand.Intn(len(englishEntries))]
		selectedEntry := speciesData.FlavorTextEntries[randomIndex]

		// Clean up the flavor text (remove newlines and extra spaces)
		flavorText := strings.ReplaceAll(selectedEntry.FlavorText, "\n", " ")
		flavorText = strings.ReplaceAll(flavorText, "\f", " ")
		flavorText = strings.Join(strings.Fields(flavorText), " ")

		// Display the Pokémon name and genus
		if genus != "" {
			fmt.Printf("%s, the %s\n", nameInfo.Formatted, genus)
		} else {
			fmt.Printf("%s\n", nameInfo.Formatted)
		}

		// Display the flavor text
		fmt.Printf("- %s", flavorText)

		// Format the game name
		formattedGameName := FormatLocationName(selectedEntry.Version.Name)

		// Display the source game
		if formattedGameName != "" {
			fmt.Printf(" (From Pokémon %s)\n", formattedGameName)
		} else {
			fmt.Println()
		}

		fmt.Println("-----")
	} else {
		fmt.Printf("No Pokédex entries found for %s\n", nameInfo.Formatted)
		fmt.Println("-----")
	}

	return nil
}
