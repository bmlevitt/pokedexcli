package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// Common error variables used across commands
var (
	ErrNoPokemonName    = errorhandling.NewInvalidInputError("No Pokémon name provided", nil)
	ErrNoLocationNumber = errorhandling.NewInvalidInputError("No location number provided", nil)
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

	// Process the Pokémon name and check if it exists in the Pokédex
	nameInfo := FormatPokemonInput(pokemonParam)
	apiName, existsInPokedex, pokemonData := CheckPokemonExists(cfg, nameInfo.APIFormat)

	// Return error if Pokemon doesn't exist in Pokedex
	if !existsInPokedex {
		// Before returning "not in Pokédex" error, verify if it's a valid Pokémon
		// by checking if it exists in the API
		_, err := cfg.pokeapiClient.GetPokemonData(nameInfo.APIFormat)

		if err != nil {
			// If the API returns a NotFound error, it's not a valid Pokémon name
			if errorhandling.IsNotFoundError(err) {
				return apiName, nameInfo, nil, false, errorhandling.InvalidPokemonNameError(nameInfo.Formatted)
			}

			// For other API errors, still return the original not-in-pokedex error
			// but log the API error if in debug mode
			if cfg.debugMode {
				log.Printf("API error while checking if %s exists: %v", nameInfo.APIFormat, err)
			}
		}

		// If we get here, either:
		// 1. The Pokémon exists in the API but not in the user's Pokédex
		// 2. There was a non-404 API error and we're treating it as if the Pokémon might be valid
		return apiName, nameInfo, nil, false, errorhandling.PokemonNotInPokedexError(nameInfo.Formatted)
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
		return pokeapi.PokemonDataResp{}, errorhandling.NewInternalError(
			fmt.Sprintf("Unexpected data type for %s", pokemonName),
			errors.New("type conversion error"))
	}
	return data, nil
}

// HandleCommandError processes errors from command functions in a standardized way.
// It logs detailed error information in debug mode and displays user-friendly messages.
//
// Parameters:
//   - cfg: The application configuration containing debug settings
//   - commandName: The name of the command that encountered the error
//   - err: The error to handle
//
// Returns:
//   - true if the command should return the error for further handling
//   - false if the error has been fully handled and the command should return nil
func HandleCommandError(cfg *config, commandName string, err error) bool {
	if err == nil {
		return false
	}

	// Log detailed error info in debug mode
	if cfg.debugMode {
		log.Printf("ERROR in command '%s': %v", commandName, err)
	}

	// For certain error types, we want to return the error for consistent handling in the REPL
	// This includes invalid input errors and "not found" errors, which should be displayed with
	// their specific user-friendly message
	if errorhandling.IsNotFoundError(err) || errorhandling.IsInvalidInputError(err) {
		return true
	}

	// For other errors, display the user-friendly message but don't propagate the error
	fmt.Printf("Error: %s\n", errorhandling.FormatUserMessage(err))
	fmt.Println("-----")
	return false
}
