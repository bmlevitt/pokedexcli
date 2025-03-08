package main

import (
	"errors"
	"fmt"
)

// commandRelease releases a Pokémon from the user's Pokédex.
// This is the opposite of catching a Pokémon - it removes the specified
// Pokémon from the user's collection.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to release
//
// Returns:
//   - An error if no Pokémon name is provided or if the Pokémon is not in the Pokédex
func commandRelease(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	inputName := params[0]

	// Convert the input name to API format if it's in a formatted style
	apiPokemonName := ConvertToAPIFormat(inputName)
	formattedName := FormatPokemonName(apiPokemonName)

	// First check for exact match
	_, exists := cfg.pokedex[apiPokemonName]
	if exists {
		// Remove the pokemon from the pokedex
		delete(cfg.pokedex, apiPokemonName)
		fmt.Printf("%s was released. Bye, %s!\n", formattedName, formattedName)
		fmt.Println("-----")
	} else {
		// Check if it's a capitalization issue by trying all keys
		found := false
		var matchedKey string
		for key := range cfg.pokedex {
			if ConvertToAPIFormat(key) == apiPokemonName {
				matchedKey = key
				found = true
				break
			}
		}

		if found {
			// Remove the pokemon from the pokedex
			delete(cfg.pokedex, matchedKey)
			fmt.Printf("%s was released. Bye, %s!\n", formattedName, formattedName)
			fmt.Println("-----")
		} else {
			return fmt.Errorf("%s is not in your Pokédex", formattedName)
		}
	}

	// Auto-save after releasing a Pokémon
	cfg.changesSinceSync++
	if cfg.changesSinceSync >= cfg.autoSaveInterval {
		if err := autoSaveIfEnabled(cfg); err != nil {
			return fmt.Errorf("error auto-saving: %w", err)
		}
		cfg.changesSinceSync = 0
	}

	return nil
}
