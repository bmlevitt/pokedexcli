package main

import (
	"errors"
	"fmt"
	"strconv"
)

// commandExplore retrieves and displays Pokémon that can be found at a specific location.
// This command is a key part of the exploration gameplay, allowing users to discover
// which Pokémon they might encounter at a given location area before attempting to catch them.
//
// The function takes a location name or a number as a parameter:
// - If a name is provided, it directly searches for that location.
// - If a number is provided, it looks up the corresponding location from the most recent map display.
//
// Parameters:
//   - cfg: The application configuration containing the API client
//   - params: Command parameters where params[0] is the location name or number to explore
//
// Returns:
//   - An error if no location name is provided, if the location doesn't exist,
//     if the number is out of range, or if the API request fails
func commandExplore(cfg *config, params []string) error {
	if len(params) == 0 {
		return errors.New("no location name provided")
	}

	// Check if the input is a number
	locationName := params[0]
	locationNumber, err := strconv.Atoi(locationName)

	// If the input is a number, convert it to the corresponding location name
	if err == nil {
		// Check if the location list exists and if the number is in range
		if len(cfg.recentLocations) == 0 {
			return errors.New("no location list available, please run 'map' command first")
		}

		// Check if the number is in range (1-based indexing)
		if locationNumber < 1 || locationNumber > len(cfg.recentLocations) {
			return fmt.Errorf("location number %d is out of range", locationNumber)
		}

		// Convert from 1-based user input to 0-based array index
		locationName = cfg.recentLocations[locationNumber-1].Name
		fmt.Printf("Exploring %s...\n", locationName)
	} else {
		// If the input is a string, print a message for consistency
		fmt.Printf("Exploring %s...\n", locationName)
	}

	resp, err := cfg.pokeapiClient.ExploreLocation(locationName)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
