package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// cliCommand represents a command that can be executed in the CLI.
// Each command has a name, description, and callback function to execute.
type cliCommand struct {
	name        string                        // Name of the command
	description string                        // Description shown in help
	callback    func(*config, []string) error // Function to execute when command is called
}

// getCommands returns a map of all available CLI commands.
// This function acts as a registry for all commands supported by the application.
// New commands should be added here to make them available in the REPL.
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "List available commands",
			callback:    commandHelp,
		},
		"explore": {
			name:        "explore",
			description: "List the pokemon found at the specified map location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "List the stats of the specified pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all pokemon currently in your pokedex",
			callback:    commandPokedex,
		},
		"release": {
			name:        "release",
			description: "Release a caught pokemon from your pokedex",
			callback:    commandRelease,
		},
		"showoff": {
			name:        "showoff",
			description: "Show off a caught pokemon using one of its moves",
			callback:    commandShowOff,
		},
		"describe": {
			name:        "describe",
			description: "Display information about a caught pokemon",
			callback:    commandDescribe,
		},
		"evolve": {
			name:        "evolve",
			description: "Evolve a pokemon that is in your pokedex",
			callback:    commandEvolve,
		},
		"save": {
			name:        "save",
			description: "Save your current Pokédex to a file",
			callback:    commandSave,
		},
		"reset": {
			name:        "reset",
			description: "Clear your Pokédex and start fresh",
			callback:    commandReset,
		},
		"autosave": {
			name:        "autosave",
			description: "Enable or disable automatic saving (on/off)",
			callback:    commandAutoSave,
		},
		"saveinterval": {
			name:        "saveinterval",
			description: "Set how often to auto-save (number of changes)",
			callback:    commandSaveInterval,
		},
		"map": {
			name:        "map",
			description: "Navigates to the first page of locations",
			callback:    commandMap,
		},
		"next": {
			name:        "next",
			description: "Navigate to the next page of locations",
			callback:    commandNext,
		},
		"prev": {
			name:        "prev",
			description: "Navigate to the previous page of locations",
			callback:    commandPrev,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

// cleanInput normalizes and splits user input into words.
// It handles whitespace and converts all text to lowercase for case-insensitive command matching.
//
// Parameters:
//   - text: The raw input string from the user
//
// Returns:
//   - A slice of lowercase words parsed from the input
func cleanInput(text string) []string {
	slice := make([]string, 0)
	words := strings.Fields(text) // Split on whitespace

	// Convert all words to lowercase for case-insensitive commands
	for _, word := range words {
		slice = append(slice, strings.ToLower(word))
	}

	return slice
}

// startREPL starts the Read-Eval-Print Loop for the Pokedex CLI.
// This is the main interaction loop that:
//  1. Reads user input (Read)
//  2. Evaluates the command (Eval)
//  3. Prints the result (Print)
//  4. Loops back for the next command
//
// The REPL continues until the user explicitly exits using the 'exit' command.
//
// Parameters:
//   - cfg: The application configuration shared across all commands
func startREPL(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	var parameters []string

	for {
		// Display prompt and get user input
		parameters = []string{}
		fmt.Print("Pokedex > ")
		reader.Scan()

		// Clean and parse the input
		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		// Get command name
		commandName := words[0]
		// fmt.Printf("command entered: %s | ", commandName) // hide later

		// Get second word as parameter if it exists
		if len(words) > 1 {
			parameters = append(parameters, words[1])
		}

		// fmt.Printf("parameter(s) applied: %s\n", parameters) // hide later

		// Look up and execute the command if it exists
		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, parameters)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("unknown command")
			continue
		}
	}
}
