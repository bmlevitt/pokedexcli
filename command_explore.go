package main

import (
	"errors"
	"fmt"
)

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
