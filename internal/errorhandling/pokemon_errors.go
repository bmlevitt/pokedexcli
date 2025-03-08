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
// It includes a more helpful message with suggestions.
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

// EvolutionNotFoundError creates a specific error for when an evolution chain is not found.
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

// FormatResourceNotFoundError creates an appropriate error message based on the resource type.
// It provides specialized error messages for known resource types, while delegating to the
// generic NewNotFoundError for unknown types.
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

// PokemonNotInPokedexError creates a specific error for when a Pokémon is not found in the user's Pokédex.
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

// InvalidPokemonNameError creates a specific error for when an invalid Pokémon name is provided.
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
