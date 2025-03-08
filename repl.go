package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
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
			description: "List the pokemon found at the specified map location number (1-20)",
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
			description: "Navigate to the first page of locations",
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
		"debug": {
			name:        "debug",
			description: "Toggle debug mode to show detailed error information",
			callback:    commandToggleDebug,
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

// startREPL begins the read-eval-print loop for the CLI application.
// It continuously reads user input, processes commands, and displays the results
// until the user chooses to exit the application with the 'exit' command.
//
// Each command is looked up in the command map, validated, and executed with
// the provided parameters. Errors are handled and displayed to the user.
//
// Parameters:
//   - cfg: The application configuration to be shared with all commands
func startREPL(cfg *config) {
	reader := bufio.NewReader(os.Stdin)
	commands := getCommands()

	// Display initial welcome and instructions
	fmt.Println("Welcome to the Pokédex!")
	fmt.Println("Type 'help' for a list of commands.")

	// Set up debug logging if enabled
	if cfg.debugMode {
		log.SetOutput(os.Stderr)
	} else {
		// Discard logs when not in debug mode
		log.SetOutput(nil)
	}

	// Loop until exit
	for {
		fmt.Print("Pokédex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			// Check if it's an EOF error, which happens when piping commands
			if err.Error() == "EOF" {
				// Exit gracefully on EOF
				fmt.Println("Exiting Pokédex. Goodbye!")
				return
			}

			// For other errors, log and continue
			fmt.Println("Error reading input:", err)
			continue
		}

		// Clean input and split into command and parameters
		cleaned := cleanInput(input)
		if len(cleaned) == 0 {
			continue
		}
		commandName := cleaned[0]
		parameters := []string{}
		if len(cleaned) > 1 {
			parameters = cleaned[1:]
		}

		// Find the command in our available commands
		command, exists := commands[commandName]
		if !exists {
			fmt.Printf("Unknown command: %s\n", commandName)
			fmt.Println("Type 'help' for a list of commands.")
			fmt.Println("-----")
			continue
		}

		// Execute the command
		err = command.callback(cfg, parameters)
		if err != nil {
			// Log the full error for debugging
			if cfg.debugMode {
				log.Printf("ERROR: [%s] %v", commandName, err)
			}

			// Format error message for display to user
			fmt.Printf("Error: %s\n", errorhandling.FormatUserMessage(err))
			fmt.Println("-----")
		}
	}
}
