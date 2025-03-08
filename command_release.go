package main

import (
	"errors"
	"fmt"
)

// commandRelease releases a Pokémon from the user's Pokédex.
// This is the opposite of catching a Pokémon - it removes the specified
// Pokémon from the user's collection.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to release
//
// Returns:
//   - An error if no Pokémon name is provided or if the Pokémon is not in the Pokédex
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
	fmt.Printf("%s was released. Bye, %s!\n", pokemonName, pokemonName)

	// Auto-save after releasing a Pokémon
	cfg.changesSinceSync++
	if cfg.changesSinceSync >= cfg.autoSaveInterval {
		if err := autoSaveIfEnabled(cfg); err != nil {
			return fmt.Errorf("error auto-saving: %w", err)
		}
		cfg.changesSinceSync = 0
	}

	return nil
}
