package main

import (
	"errors"
	"fmt"
)

// commandMap navigates to the next page of Pokémon location areas.
// It's also available under the 'next' command for more intuitive navigation.
//
// The function retrieves a list of location areas from the PokeAPI,
// displaying up to 20 locations at once. It also updates the pagination
// URLs for subsequent navigation.
//
// Parameters:
//   - cfg: The application configuration containing the pagination URLs
//   - params: Command parameters (unused)
//
// Returns:
//   - An error if there's an issue with the API request or if there are no more pages
func commandMap(cfg *config, params []string) error {
	// Get the URL to use - always use the base URL (nil) for the initial map command
	var locationsResp, err = cfg.pokeapiClient.ListLocationAreas(nil)
	if err != nil {
		return err
	}

	// Update shared state with a lock
	cfg.mutex.Lock()
	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous
	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results
	// Mark that the map has been viewed in this session
	cfg.mapViewedThisSession = true
	cfg.mutex.Unlock()

	// Display the location areas
	for i, loc := range locationsResp.Results {
		formattedLocation := FormatLocationName(loc.Name)
		fmt.Printf("%d. %s\n", i+1, formattedLocation)
	}
	fmt.Println("-----")

	return nil
}

// commandNext navigates to the next page of Pokémon location areas.
// This is the primary implementation for both 'next' and 'map' commands.
//
// The function retrieves a list of location areas from the PokeAPI,
// displaying up to 20 locations at once. It also updates the pagination
// URLs for subsequent navigation.
//
// Parameters:
//   - cfg: The application configuration containing the pagination URLs
//   - params: Command parameters (unused)
//
// Returns:
//   - An error if there's an issue with the API request or if there are no more pages
func commandNext(cfg *config, params []string) error {
	// Check if map has been viewed in this session
	if !cfg.mapViewedThisSession {
		return errors.New("you need to use the 'map' command first to load locations")
	}

	// Check if there's a next page URL
	if cfg.nextLocationURL == nil {
		return errors.New("you're on the last page")
	}

	// Get the next page of locations
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	// Update shared state with a lock
	cfg.mutex.Lock()
	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous
	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results
	cfg.mutex.Unlock()

	// Display the location areas
	for i, loc := range locationsResp.Results {
		formattedLocation := FormatLocationName(loc.Name)
		fmt.Printf("%d. %s\n", i+1, formattedLocation)
	}
	fmt.Println("-----")

	return nil
}

// commandPrev navigates to the previous page of Pokémon location areas.
// This is the primary implementation for both 'prev' and 'mapb' commands.
//
// Parameters:
//   - cfg: The application configuration containing the pagination URLs
//   - params: Command parameters (unused)
//
// Returns:
//   - An error if there's an issue with the API request or if there are no previous pages
func commandPrev(cfg *config, params []string) error {
	// Check if map has been viewed in this session
	if !cfg.mapViewedThisSession {
		return errors.New("you need to use the 'map' command first to load locations")
	}

	// Check if there's a previous page URL
	if cfg.prevLocationURL == nil {
		return errors.New("you're on the first page")
	}

	// Get the previous page of locations
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationURL)
	if err != nil {
		return err
	}

	// Update shared state with a lock
	cfg.mutex.Lock()
	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous
	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results
	cfg.mutex.Unlock()

	// Display the location areas
	for i, loc := range locationsResp.Results {
		formattedLocation := FormatLocationName(loc.Name)
		fmt.Printf("%d. %s\n", i+1, formattedLocation)
	}
	fmt.Println("-----")

	return nil
}
