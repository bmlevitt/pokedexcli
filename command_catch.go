package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandCatch(cfg *config, params []string) error {

	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	pokemonName := params[0]

	// Fetch pokemon capture rate
	resp, err := cfg.pokeapiClient.GetPokemonCaptureRate(pokemonName)
	if err != nil {
		return err
	}

	captureRate := resp.CaptureRate

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	randNum := rand.Intn(256)
	caught := randNum < captureRate
	if caught {
		pokeData, err := cfg.pokeapiClient.GetPokemonData(pokemonName)
		if err != nil {
			return err
		}
		cfg.pokedex[pokemonName] = pokeData
		fmt.Printf("%s was caught!\n", pokemonName)

	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}
