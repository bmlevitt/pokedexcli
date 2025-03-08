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
	if len(cfg.pokedex) == 0 {
		fmt.Println("you have not caught any pokemon")
	} else {
		fmt.Println("Your Pokedex:")
		for key := range cfg.pokedex {
			fmt.Printf(" - %s\n", key)
		}
	}
	return nil
}
