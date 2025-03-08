package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	locationNameOrNumber := params[0]
	locationNumber, err := strconv.Atoi(locationNameOrNumber)

	var apiLocationName string // This will store the API-formatted location name

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
		apiLocationName = cfg.recentLocations[locationNumber-1].Name
		formattedLocation := FormatLocationName(apiLocationName)
		fmt.Printf("Exploring %s...\n", formattedLocation)
	} else {
		// If the input is a string, first try to find it in the recently displayed locations
		// This will handle cases where a user copies a formatted location name from the map
		found := false
		inputLower := strings.ToLower(locationNameOrNumber)

		if len(cfg.recentLocations) > 0 {
			// Try to match by comparing the formatted location name
			for _, loc := range cfg.recentLocations {
				formattedLoc := FormatLocationName(loc.Name)
				if strings.EqualFold(formattedLoc, locationNameOrNumber) {
					apiLocationName = loc.Name
					fmt.Printf("Exploring %s...\n", formattedLoc)
					found = true
					break
				}

				// Also try matching by converting to API format
				if strings.EqualFold(loc.Name, ConvertToAPIFormat(locationNameOrNumber)) {
					apiLocationName = loc.Name
					fmt.Printf("Exploring %s...\n", FormatLocationName(loc.Name))
					found = true
					break
				}

				// Try partial matching - if the input is a prefix of a location name
				if strings.HasPrefix(strings.ToLower(formattedLoc), inputLower) {
					apiLocationName = loc.Name
					fmt.Printf("Exploring %s...\n", formattedLoc)
					found = true
					break
				}
			}
		}

		// If not found in recent locations, use the API format conversion
		if !found {
			apiLocationName = ConvertToAPIFormat(locationNameOrNumber)
			formattedLocation := FormatLocationName(apiLocationName)
			fmt.Printf("Exploring %s...\n", formattedLocation)
		}
	}

	resp, err := cfg.pokeapiClient.ExploreLocation(apiLocationName)
	if err != nil {
		// If we get an error and have recent locations, try to provide suggestions
		if len(cfg.recentLocations) > 0 {
			fmt.Println("Location not found. Did you mean one of these?")
			for i, loc := range cfg.recentLocations {
				formattedLoc := FormatLocationName(loc.Name)
				fmt.Printf("%d. %s\n", i+1, formattedLoc)
			}
			return fmt.Errorf("please select a valid location")
		}
		return err
	}

	fmt.Println("Found Pokémon:")
	for _, encounter := range resp.PokemonEncounters {
		formattedName := FormatPokemonName(encounter.Pokemon.Name)
		fmt.Printf(" - %s\n", formattedName)
	}
	fmt.Println("-----")
	return nil
}
