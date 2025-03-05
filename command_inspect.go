package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	pokemonName := params[0]

	data, exists := cfg.pokedex[pokemonName]
	if exists {
		fmt.Printf("Name: %s\n", pokemonName)
		fmt.Printf("Height: %v\n", data.Height)
		fmt.Printf("Weight: %v\n", data.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range data.Stats {
			fmt.Printf(" - %v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Type:\n")
		for _, typ := range data.Types {
			fmt.Printf(" - %v\n", typ.Type.Name)
		}

	} else {
		fmt.Printf("%s has not been caught yet", pokemonName)
	}
	return nil
}
