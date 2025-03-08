package main

import (
	"fmt"
)

// PokemonNameInfo encapsulates information about a Pokémon's name in different formats
type PokemonNameInfo struct {
	Input     string // The original user input
	APIFormat string // The name in API format (lowercase with hyphens)
	Formatted string // The formatted display name (proper capitalization)
}

// FormatPokemonInput processes a Pokémon name input and returns it in both API format and display format
func FormatPokemonInput(input string) PokemonNameInfo {
	apiFormat := ConvertToAPIFormat(input)
	formatted := FormatPokemonName(apiFormat)
	return PokemonNameInfo{
		Input:     input,
		APIFormat: apiFormat,
		Formatted: formatted,
	}
}

// CheckPokemonExists checks if a Pokémon exists in the Pokédex
// It returns the API name (which might be different from the input due to capitalization),
// whether the Pokémon exists, and the Pokémon data if it exists
func CheckPokemonExists(cfg *config, pokemonName string) (string, bool, interface{}) {
	nameInfo := FormatPokemonInput(pokemonName)

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

// HandlePokemonNotFound returns an error with a standard message when a Pokémon is not found
func HandlePokemonNotFound(pokemonName string) error {
	formatted := FormatPokemonName(pokemonName)
	return fmt.Errorf("%s is not in your Pokédex", formatted)
}

// UpdatePokedexAndSave handles all the auto-save logic after a change to the Pokédex
func UpdatePokedexAndSave(cfg *config) error {
	cfg.changesSinceSync++
	if cfg.changesSinceSync >= cfg.autoSaveInterval {
		if err := autoSaveIfEnabled(cfg); err != nil {
			return fmt.Errorf("error auto-saving: %w", err)
		}
		cfg.changesSinceSync = 0
	}
	return nil
}
