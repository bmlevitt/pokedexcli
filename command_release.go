package main

import (
	"errors"
	"fmt"
)

// commandRelease removes a Pokémon from the user's Pokédex.
// This command allows the user to release a caught Pokémon back into the wild,
// removing it from their collection. This action cannot be undone - to get the
// Pokémon back, the user would need to catch it again.
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
	fmt.Printf("%s was released back into the wild!\n", pokemonName)

	return nil
}
