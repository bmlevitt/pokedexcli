package main

import (
	"errors"
	"fmt"
)

// commandInspect displays detailed information about a Pokémon in the user's Pokédex.
// This command shows various attributes of a caught Pokémon, including:
//   - Base stats (HP, Attack, Defense, etc.)
//   - Physical attributes (Height and Weight)
//   - Types (Fire, Water, etc.)
//
// The information is only available for Pokémon that have been caught and are
// currently in the user's Pokédex.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to inspect
//
// Returns:
//   - An error if no Pokémon name is provided
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
