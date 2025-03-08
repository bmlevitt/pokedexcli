package main

import (
	"errors"
	"fmt"
)

// commandMap handles the 'map' command, showing the next page of location areas.
// This command is used to navigate forward through the paginated list of Pokémon
// locations in the game world. It fetches location data from the PokeAPI
// and updates the pagination URLs for navigation.
//
// If this is the first time the command is run, it will fetch the first page.
// Subsequent calls will fetch the next page in the sequence if available.
//
// Parameters:
//   - cfg: The application configuration containing the API client and pagination URLs
//   - params: Command parameters (not used for this command)
//
// Returns:
//   - An error if the API request fails
func commandMap(cfg *config, params []string) error {
	// Get the next page of location areas (or first page if nextLocationURL is nil)
	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	// Print all locations from the response
	fmt.Println("Location areas:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s \n", area.Name)
	}

	// Update the pagination URLs in the config
	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous
	return nil
}

// commandMapb handles the 'mapb' command, showing the previous page of location areas.
// This command complements the 'map' command by enabling backward navigation through
// the paginated list of Pokémon locations. It fetches the previous page of location data
// and updates the pagination URLs accordingly.
//
// Parameters:
//   - cfg: The application configuration containing the API client and pagination URLs
//   - params: Command parameters (not used for this command)
//
// Returns:
//   - An error if on the first page (no previous data) or if the API request fails
func commandMapb(cfg *config, params []string) error {
	// Check if we can go back to the previous page
	if cfg.prevLocationURL == nil {
		return errors.New("you are on the first response - no previous data available")
	}

	// Get the previous page of location areas
	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationURL)
	if err != nil {
		return err
	}

	// Print all locations from the response
	fmt.Println("Location areas:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s \n", area.Name)
	}

	// Update the pagination URLs in the config
	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous
	return nil
}
