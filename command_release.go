package main

import (
	"errors"
	"fmt"
)

func commandRelease(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	pokemonName := params[0]

	// Check if the pokemon exists in the pokedex
	_, exists := cfg.pokedex[pokemonName]
	if !exists {
		return fmt.Errorf("%s is not in your pokedex", pokemonName)
	}

	// Remove the pokemon from the pokedex
	delete(cfg.pokedex, pokemonName)
	fmt.Printf("%s was released back into the wild!\n", pokemonName)

	return nil
}
