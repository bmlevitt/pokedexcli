package main

import "fmt"

// commandPokedex displays a list of all Pokémon the user has caught.
// This command provides a simple inventory view of the user's collection,
// listing the names of all Pokémon currently in their Pokédex.
//
// If the Pokédex is empty (no Pokémon have been caught), a message indicating
// this is displayed instead of an empty list.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - Always returns nil as this command cannot fail under normal circumstances
func commandPokedex(cfg *config, params []string) error {
	// Acquire a read lock before accessing the pokedex
	cfg.mutex.RLock()
	pokedexEmpty := len(cfg.pokedex) == 0
	cfg.mutex.RUnlock()

	if pokedexEmpty {
		fmt.Println("You have not caught any Pokémon yet")
	} else {
		fmt.Println("Your Pokédex:")

		// Acquire a read lock for iteration
		cfg.mutex.RLock()
		for key := range cfg.pokedex {
			formattedName := FormatPokemonName(key)
			fmt.Printf(" - %s\n", formattedName)
		}
		cfg.mutex.RUnlock()

		fmt.Println("-----")
	}
	return nil
}
