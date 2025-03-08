package main

import (
	"errors"
	"fmt"
)

// commandExplore retrieves and displays Pokémon that can be found at a specific location.
// This command is a key part of the exploration gameplay, allowing users to discover
// which Pokémon they might encounter at a given location area before attempting to catch them.
//
// The function takes a location name as a parameter, queries the PokeAPI for Pokémon
// encounters at that location, and prints a list of all Pokémon that can be found there.
//
// Parameters:
//   - cfg: The application configuration containing the API client
//   - params: Command parameters where params[0] is the location name to explore
//
// Returns:
//   - An error if no location name is provided or if the API request fails
func commandExplore(cfg *config, params []string) error {

	if len(params) == 0 {
		return errors.New("no location name provided")
	}

	locationName := params[0]
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
