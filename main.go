package main

import (
	"fmt"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// config holds the application state and configuration settings.
// It maintains the API client, pagination URLs for location navigation,
// the user's collection of caught Pokémon (the Pokédex), and auto-save settings.
type config struct {
	pokeapiClient    pokeapi.Client                     // Client for making Pokemon API requests
	nextLocationURL  *string                            // URL for the next page of locations
	prevLocationURL  *string                            // URL for the previous page of locations
	pokedex          map[string]pokeapi.PokemonDataResp // Map of caught Pokemon indexed by name
	autoSaveEnabled  bool                               // Whether to automatically save after changes
	autoSaveInterval int                                // How many changes before auto-saving (if enabled)
	changesSinceSync int                                // Counter for changes since last save
	recentLocations  []pokeapi.NamedAPIResource         // Most recent list of locations displayed
}

// main initializes the application and starts the command-line interface.
// It creates a new API client with a 1-hour cache duration to reduce API calls,
// initializes an empty Pokédex to store caught Pokémon, and loads any saved data.
func main() {
	// Initialize the configuration with a new Pokemon API client and default settings
	cfg := config{
		pokeapiClient:    pokeapi.NewClient(time.Hour),
		pokedex:          make(map[string]pokeapi.PokemonDataResp),
		autoSaveEnabled:  true, // Auto-save is enabled by default
		autoSaveInterval: 1,    // Save after every change by default
		changesSinceSync: 0,    // No changes yet
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
