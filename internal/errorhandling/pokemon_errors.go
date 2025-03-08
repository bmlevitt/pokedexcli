// This file provides specialized error types and functions for Pokemon-related errors.
// It defines helper functions for creating consistent error messages
// for different types of Pokemon-related error scenarios.
package errorhandling

import (
	"fmt"
	"strings"
)

// Resource types for error messages
const (
	ResourcePokemon          = "Pokémon"
	ResourceLocation         = "location"
	ResourcePokemonSpecies   = "Pokémon species"
	ResourceEvolutionChain   = "evolution chain"
	ResourcePokemonMove      = "Pokémon move"
	ResourcePokemonAbility   = "Pokémon ability"
	ResourcePokemonEncounter = "Pokémon encounter"
)

// PokemonNotFoundError creates a specific error for when a Pokémon is not found.
// It includes a more helpful message with suggestions based on common mistakes.
//
// Parameters:
//   - pokemonName: The name of the Pokémon that wasn't found
//   - err: The original error that caused the not found condition
//
// Returns:
//   - An AppError with detailed information about the not found condition
func PokemonNotFoundError(pokemonName string, err error) *AppError {
	message := fmt.Sprintf("The Pokémon '%s' was not found. Please check the spelling and try again.", pokemonName)

	// Add suggestions based on common issues
	if strings.Contains(pokemonName, " ") {
		message += " Note: Pokémon names should not contain spaces (use '-' instead)."
	}
	if strings.Contains(pokemonName, ".") || strings.Contains(pokemonName, "'") {
		message += " Note: Special characters like '.' and ''' are not used in API Pokémon names."
	}

	return &AppError{
		Type:       NotFound,
		StatusCode: 404,
		Message:    message,
		Err:        err,
		Context: map[string]string{
			"resourceType": ResourcePokemon,
			"resourceName": pokemonName,
		},
	}
}

// LocationNotFoundError creates a specific error for when a location is not found.
// This provides a user-friendly message specific to location resources.
//
// Parameters:
//   - locationName: The name of the location that wasn't found
//   - err: The original error that caused the not found condition
//
// Returns:
//   - An AppError with detailed information about the not found condition
func LocationNotFoundError(locationName string, err error) *AppError {
	message := fmt.Sprintf("The location '%s' was not found. Please run the 'map' command to see available locations.", locationName)

	return &AppError{
		Type:       NotFound,
		StatusCode: 404,
		Message:    message,
		Err:        err,
		Context: map[string]string{
			"resourceType": ResourceLocation,
			"resourceName": locationName,
		},
	}
}

// EvolutionNotFoundError creates a specific error for when a Pokémon's evolution data is not found.
// This is useful when trying to evolve a Pokémon and the evolution chain cannot be retrieved.
//
// Parameters:
//   - pokemonName: The name of the Pokémon whose evolution data wasn't found
//   - err: The original error that caused the not found condition
//
// Returns:
//   - An AppError with detailed information about the not found condition
func EvolutionNotFoundError(pokemonName string, err error) *AppError {
	message := fmt.Sprintf("No evolution information found for '%s'. This Pokémon might not have any evolutions.", pokemonName)

	return &AppError{
		Type:       NotFound,
		StatusCode: 404,
		Message:    message,
		Err:        err,
		Context: map[string]string{
			"resourceType": ResourceEvolutionChain,
			"resourceName": pokemonName,
		},
	}
}

// FormatResourceNotFoundError creates a generic resource not found error with a
// standardized format for the error message.
//
// Parameters:
//   - resourceType: The type of resource that wasn't found (e.g., "Pokémon", "location")
//   - resourceName: The name of the specific resource that wasn't found
//   - err: The original error that caused the not found condition
//
// Returns:
//   - An AppError with detailed information about the not found condition
func FormatResourceNotFoundError(resourceType, resourceName string, err error) *AppError {
	switch resourceType {
	case ResourcePokemon:
		return PokemonNotFoundError(resourceName, err)
	case ResourceLocation:
		return LocationNotFoundError(resourceName, err)
	case ResourceEvolutionChain:
		return EvolutionNotFoundError(resourceName, err)
	default:
		return NewNotFoundError(resourceType, resourceName, err)
	}
}

// PokemonNotInPokedexError creates an error for when a user tries to interact with
// a Pokémon that they haven't caught yet (not in their Pokédex).
//
// Parameters:
//   - pokemonName: The name of the Pokémon that's not in the user's Pokédex
//
// Returns:
//   - An AppError with a user-friendly message about the Pokémon not being in the Pokédex
func PokemonNotInPokedexError(pokemonName string) *AppError {
	message := fmt.Sprintf("Pokémon '%s' is not in your Pokédex. Try catching it first!", pokemonName)

	return &AppError{
		Type:       NotFound,
		StatusCode: 404,
		Message:    message,
		Context: map[string]string{
			"resourceType": "Pokémon in Pokédex",
			"resourceName": pokemonName,
		},
	}
}

// InvalidPokemonNameError creates an error for when a user provides an invalid Pokémon name.
// This is used when the name format is incorrect or when the Pokémon doesn't exist.
//
// Parameters:
//   - pokemonName: The invalid Pokémon name provided by the user
//
// Returns:
//   - An AppError with a user-friendly message about the invalid Pokémon name
func InvalidPokemonNameError(pokemonName string) *AppError {
	message := fmt.Sprintf("'%s' is not a valid Pokémon name. Please check your spelling - this Pokémon doesn't exist.", pokemonName)

	return &AppError{
		Type:       InvalidInput,
		StatusCode: 400,
		Message:    message,
		Context: map[string]string{
			"resourceType": "Pokémon name",
			"resourceName": pokemonName,
		},
	}
}
