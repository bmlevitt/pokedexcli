package main

import (
	"errors"
	"fmt"
	"math/rand"
)

// commandCatch attempts to catch a specified Pokémon and add it to the user's Pokédex.
// The function simulates the catching mechanic from the Pokémon games by using
// the Pokémon's capture rate from the API to determine catch probability.
//
// The catch probability is calculated by comparing a random number (0-255)
// against the Pokémon's capture rate. If the random number is less than the
// capture rate, the catch is successful.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon name to catch
//
// Returns:
//   - An error if no Pokémon name is provided or if there's an issue with the API request
func commandCatch(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}

	// Process the Pokémon name input
	nameInfo := FormatPokemonInput(params[0])

	// Fetch pokemon capture rate
	resp, err := cfg.pokeapiClient.GetPokemonCaptureRate(nameInfo.APIFormat)
	if err != nil {
		return err
	}

	captureRate := resp.CaptureRate
	fmt.Printf("Throwing a Pokéball at %s...\n", nameInfo.Formatted)
	randNum := rand.Intn(256)
	caught := randNum < captureRate
	if caught {
		pokeData, err := cfg.pokeapiClient.GetPokemonData(nameInfo.APIFormat)
		if err != nil {
			return err
		}

		// Lock the config before modifying the pokedex
		cfg.mutex.Lock()
		cfg.pokedex[nameInfo.APIFormat] = pokeData
		cfg.mutex.Unlock()

		fmt.Printf("%s was caught!\n", nameInfo.Formatted)

		// Auto-save after catching a Pokémon
		if err := UpdatePokedexAndSave(cfg); err != nil {
			return err
		}
	} else {
		fmt.Printf("%s escaped!\n", nameInfo.Formatted)
	}
	fmt.Println("-----")
	return nil
}
