// This file contains utility functions for managing the Pokédex.
// It provides functions for formatting Pokémon names, checking if Pokémon exist
// in the user's collection, and handling Pokédex updates with auto-save functionality.
package main

import (
	"fmt"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// PokemonNameInfo encapsulates information about a Pokémon's name in different formats.
// This struct is used to maintain consistency when a Pokémon name needs to be referenced
// in multiple formats (user input, API format, and display format).
type PokemonNameInfo struct {
	Input     string // The original user input
	APIFormat string // The name in API format (lowercase with hyphens)
	Formatted string // The formatted display name (proper capitalization)
}

// FormatPokemonInput processes a Pokémon name input and returns it in both API format and display format.
// This function ensures consistent handling of Pokémon names throughout the application.
//
// Parameters:
//   - input: The raw Pokémon name as entered by the user
//
// Returns:
//   - A PokemonNameInfo struct containing the name in various formats
func FormatPokemonInput(input string) PokemonNameInfo {
	apiFormat := ConvertToAPIFormat(input)
	formatted := FormatPokemonName(apiFormat)
	return PokemonNameInfo{
		Input:     input,
		APIFormat: apiFormat,
		Formatted: formatted,
	}
}

// CheckPokemonExists checks if a Pokémon exists in the user's Pokédex.
// It handles case-insensitive matching and returns the Pokémon's data if found.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - pokemonName: The name of the Pokémon to check for
//
// Returns:
//   - The API-formatted name of the Pokémon
//   - A boolean indicating whether the Pokémon exists in the Pokédex
//   - The Pokémon's data if it exists, nil otherwise
func CheckPokemonExists(cfg *config, pokemonName string) (string, bool, interface{}) {
	nameInfo := FormatPokemonInput(pokemonName)

	// Acquire a read lock before accessing the pokedex
	cfg.mutex.RLock()
	defer cfg.mutex.RUnlock()

	// Check if the pokemon exists directly
	pokemonData, exists := cfg.pokedex[nameInfo.APIFormat]
	if exists {
		return nameInfo.APIFormat, true, pokemonData
	}

	// Check if it's a capitalization issue by trying all keys
	for key, data := range cfg.pokedex {
		if ConvertToAPIFormat(key) == nameInfo.APIFormat {
			return key, true, data
		}
	}

	return nameInfo.APIFormat, false, nil
}

// HandlePokemonNotInPokedex returns a standardized error when a Pokémon is not found in the Pokédex.
// This ensures consistent error messaging for this common error condition.
//
// Parameters:
//   - pokemonName: The name of the Pokémon that wasn't found
//
// Returns:
//   - A formatted error indicating the Pokémon is not in the Pokédex
func HandlePokemonNotInPokedex(pokemonName string) error {
	return errorhandling.PokemonNotInPokedexError(pokemonName)
}

// UpdatePokedexAndSave handles all the auto-save logic after a change to the Pokédex.
// It increments the change counter and triggers an auto-save if the threshold is reached.
//
// Parameters:
//   - cfg: The application configuration containing auto-save settings
//
// Returns:
//   - An error if the auto-save fails, nil otherwise
func UpdatePokedexAndSave(cfg *config) error {
	// Lock the config before modifying the counter
	cfg.mutex.Lock()
	cfg.changesSinceSync++
	shouldSave := cfg.changesSinceSync >= cfg.autoSaveInterval
	if shouldSave {
		cfg.changesSinceSync = 0
	}
	cfg.mutex.Unlock()

	if shouldSave {
		if err := autoSaveIfEnabled(cfg); err != nil {
			return fmt.Errorf("error auto-saving: %w", err)
		}
	}
	return nil
}
