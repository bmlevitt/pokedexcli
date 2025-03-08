package main

import (
	"fmt"
)

// commandHelp displays a list of all available commands with their descriptions.
// This is the main help system for the application, providing users with information
// about what commands are available and what they do.
//
// The function iterates through all registered commands from getCommands() and
// prints each command name alongside its description in a formatted list.
//
// Parameters:
//   - cfg: The application configuration (not used in this command)
//   - params: Command parameters (not used in this command)
//
// Returns:
//   - Always returns nil as this command cannot fail under normal circumstances
//
// Side Effects:
//   - Prints the welcome message and list of available commands to stdout
func commandHelp(cfg *config, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("-----")
	for _, cmd := range getCommands() {
		fmt.Printf("%s | %s \n", cmd.name, cmd.description)
	}
	fmt.Println("-----")
	return nil
}
