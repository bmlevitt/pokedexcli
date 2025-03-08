package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// config holds the application's global configuration and state.
// It includes API clients, navigation state, and the user's Pokédex data.
type config struct {
	pokeapiClient        pokeapi.Client                     // Client for making Pokemon API requests
	nextLocationURL      *string                            // URL for the next page of map locations
	prevLocationURL      *string                            // URL for the previous page of map locations
	pokedex              map[string]pokeapi.PokemonDataResp // Map of caught Pokemon indexed by name
	autoSaveEnabled      bool                               // Whether to automatically save after changes
	autoSaveInterval     int                                // How many changes before auto-saving (if enabled)
	changesSinceSync     int                                // Counter for changes since last save
	recentLocations      []pokeapi.NamedAPIResource         // Most recent list of map locations displayed
	mapViewedThisSession bool                               // Whether the map command has been used in this session
	debugMode            bool                               // Whether to show detailed error messages
	mutex                sync.RWMutex                       // Mutex to protect access to shared data
	// Only one mutex -- risk is low in this simple app
}

// main is the entry point for the Pokédex CLI application.
// It creates a new API client with a 1-hour cache duration to reduce API calls,
// initializes an empty Pokédex to store caught Pokémon, and loads any saved data.
// After initialization, it starts the interactive REPL (Read-Eval-Print Loop)
// that accepts user commands and processes them.
//
// The function handles startup errors gracefully, particularly for loading saved data,
// by displaying friendly error messages to the user instead of crashing.
//
// Side Effects:
//   - Creates and initializes the global application state
//   - Attempts to read from the saved data file
//   - Prints startup messages to stdout
//   - Starts the interactive command loop that runs until program exit
func main() {
	// Initialize the configuration with a new Pokemon API client and default settings
	cfg := config{
		pokeapiClient:        pokeapi.NewClient(time.Hour),
		pokedex:              make(map[string]pokeapi.PokemonDataResp),
		autoSaveEnabled:      true,  // Auto-save is enabled by default
		autoSaveInterval:     1,     // Save after every change by default
		changesSinceSync:     0,     // No changes yet
		mapViewedThisSession: false, // Map hasn't been viewed in this session yet
		debugMode:            false, // Debug mode is disabled by default
	}

	// Try to load saved data
	err := loadPokedexData(&cfg)
	if err != nil {
		fmt.Printf("Warning: Could not load saved Pokédex data: %v\n", err)
	} else if len(cfg.pokedex) > 0 {
		fmt.Printf("Loaded Pokédex with %d Pokémon\n", len(cfg.pokedex))
	}
	fmt.Println("-----")

	// Start the REPL (Read-Eval-Print Loop) with our config
	startREPL(&cfg)
}
