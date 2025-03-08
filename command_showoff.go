package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// commandShowOff displays a Pokémon from the user's Pokédex performing a random move.
// This command simulates a Pokémon using one of its moves for display purposes.
// It randomly selects a move from the Pokémon's known moves and presents it in
// a formatted message.
//
// The command can only be used with Pokémon that are currently in the user's Pokédex
// and that know at least one move.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to show off
//
// Returns:
//   - An error if no Pokémon name is provided, if the Pokémon is not in the Pokédex,
//     or if the Pokémon doesn't know any moves
func commandShowOff(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}

	// Process the Pokémon name and check if it exists
	apiName, exists, pokemonData := CheckPokemonExists(cfg, params[0])
	nameInfo := FormatPokemonInput(apiName)

	if !exists {
		return HandlePokemonNotFound(nameInfo.APIFormat)
	}

	// Type assertion for the pokemon data
	pokemon, ok := pokemonData.(pokeapi.PokemonDataResp)
	if !ok {
		return fmt.Errorf("unexpected data type for %s", nameInfo.Formatted)
	}

	// Check if the pokemon has any moves
	if len(pokemon.Moves) == 0 {
		return fmt.Errorf("%s doesn't know any moves", nameInfo.Formatted)
	}

	// Select a random move
	randomIndex := rand.Intn(len(pokemon.Moves))
	moveName := pokemon.Moves[randomIndex].Move.Name

	// Format the move name for better display
	formattedMove := FormatMoveName(moveName)

	// Show off the pokemon using the move
	fmt.Printf("%s used %s!\n", nameInfo.Formatted, formattedMove)
	fmt.Println("-----")

	return nil
}
