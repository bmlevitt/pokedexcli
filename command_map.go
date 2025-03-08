// This file implements navigation commands for the Pokédex CLI application.
// It provides functionality for exploring the Pokémon world map by navigating
// through paginated lists of location areas from the PokeAPI.
package main

import (
	"fmt"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// commandMap displays the first page of Pokémon location areas.
// It retrieves data from the PokeAPI and displays a numbered list of up to 20 locations.
// This command serves as the entry point for map exploration before using 'next' and 'prev'.
//
// The function stores the current location list and pagination URLs in the application config
// for subsequent navigation commands.
//
// Parameters:
//   - cfg: The application configuration for storing location data and pagination URLs
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - An error if there's an issue with the API request
func commandMap(cfg *config, params []string) error {
	// Get the URL to use - always use the base URL (nil) for the initial map command
	var locationsResp, err = cfg.pokeapiClient.ListLocationAreas(nil)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "map", err) {
			return err
		}
		return nil
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
// It uses the nextLocationURL stored in the application config to retrieve
// the next set of up to 20 locations from the PokeAPI.
//
// If there are no more pages (nextLocationURL is nil), an error is returned.
// This command can only be used after the 'map' command has been used at least once.
//
// Parameters:
//   - cfg: The application configuration containing the next page URL
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - An error if there are no more pages or if there's an issue with the API request
func commandNext(cfg *config, params []string) error {
	// Check if map has been viewed in this session
	if !cfg.mapViewedThisSession {
		// Create a specific error for this case
		noMapViewedErr := errorhandling.NewInvalidInputError("You need to use the 'map' command first to load locations", nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "next", noMapViewedErr) {
			return noMapViewedErr
		}
		return nil
	}

	// Lock to prevent race conditions when reading/writing shared state
	cfg.mutex.RLock()
	nextURL := cfg.nextLocationURL
	cfg.mutex.RUnlock()

	// Check if there's a next page
	if nextURL == nil {
		// Create a specific error for this case
		noNextErr := errorhandling.NewInvalidInputError("You're on the last page", nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "next", noNextErr) {
			return noNextErr
		}
		return nil
	}

	// Make the API request with the next URL
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationURL)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "next", err) {
			return err
		}
		return nil
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
// It uses the prevLocationURL stored in the application config to retrieve
// the previous set of up to 20 locations from the PokeAPI.
//
// If there are no previous pages (prevLocationURL is nil), an error is returned.
// This command can only be used after navigating forward at least once with 'next'.
//
// Parameters:
//   - cfg: The application configuration containing the previous page URL
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - An error if there are no previous pages or if there's an issue with the API request
func commandPrev(cfg *config, params []string) error {
	// Check if map has been viewed in this session
	if !cfg.mapViewedThisSession {
		// Create a specific error for this case
		noMapViewedErr := errorhandling.NewInvalidInputError("You need to use the 'map' command first to load locations", nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "prev", noMapViewedErr) {
			return noMapViewedErr
		}
		return nil
	}

	// Lock to prevent race conditions when reading/writing shared state
	cfg.mutex.RLock()
	prevURL := cfg.prevLocationURL
	cfg.mutex.RUnlock()

	// Check if there's a previous page
	if prevURL == nil {
		// Create a specific error for this case
		noPrevErr := errorhandling.NewInvalidInputError("You're on the first page", nil)

		// Use standardized error handling
		if HandleCommandError(cfg, "prev", noPrevErr) {
			return noPrevErr
		}
		return nil
	}

	// Make the API request with the previous URL
	locationsResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationURL)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "prev", err) {
			return err
		}
		return nil
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
