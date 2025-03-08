package main

import (
	"fmt"
	"os"
)

// commandExit terminates the application with a goodbye message.
// This command provides a clean way for users to exit the Pokédex CLI.
// It prints a farewell message and then immediately terminates the program
// with a successful status code (0).
//
// Parameters:
//   - cfg: The application configuration (not used in this command)
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - This function does not actually return as it calls os.Exit()
//     which immediately terminates the program
func commandExit(cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
