package main

import (
	"errors"
	"fmt"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// Common error messages used across commands
var (
	ErrNoPokemonName       = errors.New("no pokemon name provided")
	ErrNoLocationNumber    = errors.New("no location number provided")
	ErrPokemonNotInPokedex = errors.New("pokemon is not in your Pokédex")
)

// ValidatePokemonParam checks if a Pokemon name parameter was provided
// and returns the name if it was, or an error if it wasn't.
//
// Parameters:
//   - params: Command parameters where params[0] should be the Pokémon name
//
// Returns:
//   - The Pokémon name if provided
//   - An error if no Pokémon name is provided
func ValidatePokemonParam(params []string) (string, error) {
	if len(params) == 0 {
		return "", ErrNoPokemonName
	}
	return params[0], nil
}

// GetPokemonIfExists validates the Pokemon parameter, checks if it exists in the Pokedex,
// and returns the relevant information with appropriate error handling.
//
// This function combines several common operations:
// 1. Checking if a Pokemon name was provided as a parameter
// 2. Formatting the Pokemon name properly
// 3. Checking if the Pokemon exists in the Pokedex
// 4. Returning standardized errors if it doesn't
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] should be the Pokémon name
//
// Returns:
//   - apiName: The Pokemon name in API format
//   - nameInfo: Structured information about the Pokemon name in different formats
//   - pokemonData: The Pokemon data if it exists in the Pokedex
//   - exists: Whether the Pokemon exists in the Pokedex
//   - err: An error if no parameter was provided or if the Pokemon doesn't exist
func GetPokemonIfExists(cfg *config, params []string) (string, PokemonNameInfo, interface{}, bool, error) {
	// Check if Pokemon name parameter was provided
	pokemonParam, err := ValidatePokemonParam(params)
	if err != nil {
		return "", PokemonNameInfo{}, nil, false, err
	}

	// Process the Pokémon name and check if it exists
	nameInfo := FormatPokemonInput(pokemonParam)
	apiName, exists, pokemonData := CheckPokemonExists(cfg, nameInfo.APIFormat)

	// Return error if Pokemon doesn't exist in Pokedex
	if !exists {
		return apiName, nameInfo, nil, false, HandlePokemonNotFound(nameInfo.APIFormat)
	}

	return apiName, nameInfo, pokemonData, true, nil
}

// GetTypedPokemonData converts the generic Pokemon data to a typed Pokemon response.
// This is useful when specific Pokemon data fields are needed.
//
// Parameters:
//   - pokemonData: The raw Pokemon data from the Pokedex
//   - pokemonName: The formatted Pokemon name (for error messages)
//
// Returns:
//   - The typed Pokemon data
//   - An error if the conversion fails
func GetTypedPokemonData(pokemonData interface{}, pokemonName string) (pokeapi.PokemonDataResp, error) {
	data, ok := pokemonData.(pokeapi.PokemonDataResp)
	if !ok {
		return pokeapi.PokemonDataResp{}, fmt.Errorf("unexpected data type for %s", pokemonName)
	}
	return data, nil
}
