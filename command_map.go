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
	// Get the URL to use - either the next page URL or the base URL
	var locationsResp, err = cfg.pokeapiClient.ListLocationAreas(nil)
	if err != nil {
		return err
	}

	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous

	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results

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
	if cfg.nextLocationURL == nil {
		return errors.New("no more pages available")
	}

	// Get the next page of locations
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous

	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results

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
	if cfg.prevLocationURL == nil {
		return errors.New("already at the first page")
	}

	// Get the previous page of locations
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationURL)
	if err != nil {
		return err
	}

	// Update pagination URLs
	cfg.nextLocationURL = locationsResp.Next
	cfg.prevLocationURL = locationsResp.Previous

	// Store the location results for reference by other commands
	cfg.recentLocations = locationsResp.Results

	// Display the location areas
	for i, loc := range locationsResp.Results {
		formattedLocation := FormatLocationName(loc.Name)
		fmt.Printf("%d. %s\n", i+1, formattedLocation)
	}
	fmt.Println("-----")

	return nil
}
