package main

import (
	"errors"
	"fmt"
)

// commandMap handles the 'map' command, showing the next page of location areas
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

// commandMapb handles the 'mapb' command, showing the previous page of location areas
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
