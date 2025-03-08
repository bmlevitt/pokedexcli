package main

import (
	"errors"
	"fmt"
	"strconv"
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
// displayed in the most recent map command.
//
// Parameters:
//   - cfg: The application configuration containing the API client
//   - params: Command parameters where params[0] is the location number (1-20) to explore
//
// Returns:
//   - An error if no location number is provided, if the number is invalid,
//     if the number is out of range, or if the API request fails
func commandExplore(cfg *config, params []string) error {
	// Validate the location parameter
	locNumStr, err := ValidateLocationParam(params)
	if err != nil {
		return err
	}

	// Parse the location number from input
	locationNumber, err := strconv.Atoi(locNumStr)
	if err != nil {
		return errors.New("invalid location number: please provide a number between 1-20")
	}

	// Check if the location list exists
	if len(cfg.recentLocations) == 0 {
		return errors.New("no location list available, please run 'map' command first")
	}

	// Check if the number is in range (1-based indexing)
	if locationNumber < 1 || locationNumber > len(cfg.recentLocations) {
		return fmt.Errorf("location number %d is out of range", locationNumber)
	}

	// Convert from 1-based user input to 0-based array index
	apiLocationName := cfg.recentLocations[locationNumber-1].Name
	formattedLocation := FormatLocationName(apiLocationName)
	fmt.Printf("Exploring %s...\n", formattedLocation)

	// Make the API request to explore the location
	resp, err := cfg.pokeapiClient.ExploreLocation(apiLocationName)
	if err != nil {
		return err
	}

	// Display the Pokémon found at this location
	fmt.Println("Found Pokémon:")
	for _, encounter := range resp.PokemonEncounters {
		formattedName := FormatPokemonName(encounter.Pokemon.Name)
		fmt.Printf(" - %s\n", formattedName)
	}
	fmt.Println("-----")
	return nil
}
