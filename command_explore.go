// This file implements the explore command for the Pokédex CLI application.
// It allows users to discover which Pokémon can be found at specific locations
// in the Pokémon world, providing a key part of the exploration gameplay.
package main

import (
	"fmt"
	"strconv"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// ValidateLocationParam checks if a location number parameter was provided
// and returns the value if it was, or an error if it wasn't.
//
// Parameters:
//   - params: Command parameters where params[0] should be the location number
//
// Returns:
//   - The location number parameter if provided
//   - An error if no location number is provided
func ValidateLocationParam(params []string) (string, error) {
	if len(params) == 0 {
		return "", ErrNoLocationNumber
	}
	return params[0], nil
}

// commandExplore retrieves and displays Pokémon that can be found at a specific location.
// This command is a key part of the exploration gameplay, allowing users to discover
// which Pokémon they might encounter at a given location area before attempting to catch them.
//
// The function takes a location number as a parameter, which corresponds to the location
// displayed by the map command (1-20). It then fetches a list of Pokémon that can be
// encountered at that location and displays them to the user.
//
// Parameters:
//   - cfg: The application configuration containing the API client and recent locations
//   - params: Command parameters where params[0] is the location number to explore
//
// Returns:
//   - An error if no location number is provided, if the number is invalid,
//     if the map hasn't been viewed yet, or if there's an issue with the API request
func commandExplore(cfg *config, params []string) error {
	// Validate the location parameter
	locNumStr, err := ValidateLocationParam(params)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "explore", err) {
			return err
		}
		return nil
	}

	// Parse the location number from input
	locationNumber, err := strconv.Atoi(locNumStr)
	if err != nil {
		// Create a specific error for this case
		invalidNumErr := errorhandling.NewInvalidInputError("Invalid location number: please provide a number between 1-20", err)

		// Use standardized error handling
		if HandleCommandError(cfg, "explore", invalidNumErr) {
			return invalidNumErr
		}
		return nil
	}

	// Check if the location list exists
	if len(cfg.recentLocations) == 0 {
		// Create a specific error for this case
		noLocationsErr := errorhandling.NewInvalidInputError("No location list available, please run the 'map' command first", nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "explore", noLocationsErr) {
			return noLocationsErr
		}
		return nil
	}

	// Check if the number is in range (1-based indexing)
	if locationNumber < 1 || locationNumber > len(cfg.recentLocations) {
		// Create a specific error for this case
		outOfRangeErr := errorhandling.NewInvalidInputError(
			fmt.Sprintf("Location number %d is out of range (valid range: 1-%d)",
				locationNumber, len(cfg.recentLocations)), nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "explore", outOfRangeErr) {
			return outOfRangeErr
		}
		return nil
	}

	// Convert from 1-based user input to 0-based array index
	apiLocationName := cfg.recentLocations[locationNumber-1].Name
	formattedLocation := FormatLocationName(apiLocationName)
	fmt.Printf("Exploring %s...\n", formattedLocation)

	// Make the API request to explore the location
	resp, err := cfg.pokeapiClient.ExploreLocation(apiLocationName)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "explore", err) {
			return err
		}
		return nil
	}

	// Display the Pokémon found at this location
	if len(resp.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found at this location.")
	} else {
		fmt.Println("Found Pokémon:")
		for i, encounter := range resp.PokemonEncounters {
			formattedName := FormatPokemonName(encounter.Pokemon.Name)
			fmt.Printf("%d. %s\n", i+1, formattedName)
		}
	}
	fmt.Println("-----")
	return nil
}
