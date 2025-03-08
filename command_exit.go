package main

import (
	"fmt"
	"os"
)

// commandExit handles the exit command, which gracefully terminates the program.
// Before exiting, it ensures that the user's Pokédex data is saved to disk
// to prevent data loss.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex data
//   - params: Command parameters (unused)
//
// Returns:
//   - Never returns as the program exits
func commandExit(cfg *config, params []string) error {
	// Save the Pokédex data before exiting
	err := savePokedexData(cfg)
	if err != nil {
		fmt.Printf("Warning: Could not save Pokédex data: %v\n", err)
	} else {
		fmt.Println("Pokédex data saved!")
	}

	fmt.Println("Thanks for using the Pokédex! See you next time!")
	os.Exit(0)
	return nil // This line is never reached but keeps the compiler happy
}
