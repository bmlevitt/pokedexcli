package main

import (
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// config holds the application state and configuration settings.
// It maintains the API client, pagination URLs for location navigation,
// and the user's collection of caught Pokémon (the Pokédex).
type config struct {
	pokeapiClient   pokeapi.Client                     // Client for making Pokemon API requests
	nextLocationURL *string                            // URL for the next page of locations
	prevLocationURL *string                            // URL for the previous page of locations
	pokedex         map[string]pokeapi.PokemonDataResp // Map of caught Pokemon indexed by name
}

// main initializes the application and starts the command-line interface.
// It creates a new API client with a 1-hour cache duration to reduce API calls
// and initializes an empty Pokédex to store caught Pokémon.
func main() {
	// Initialize the configuration with a new Pokemon API client
	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		pokedex:       make(map[string]pokeapi.PokemonDataResp),
	}

	// Start the REPL (Read-Eval-Print Loop) with our config
	startREPL(&cfg)
}
