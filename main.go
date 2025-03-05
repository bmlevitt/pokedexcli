package main

import (
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// config holds the application state
type config struct {
	pokeapiClient   pokeapi.Client // Client for making Pokemon API requests
	nextLocationURL *string        // URL for the next page of locations
	prevLocationURL *string        // URL for the previous page of locations
	pokedex         map[string]pokeapi.PokemonDataResp
}

func main() {
	// Initialize the configuration with a new Pokemon API client
	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		pokedex:       make(map[string]pokeapi.PokemonDataResp),
	}

	// Start the REPL (Read-Eval-Print Loop) with our config
	startREPL(&cfg)
}
