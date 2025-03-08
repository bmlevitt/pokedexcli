package main

import (
	"errors"
	"fmt"
)

// commandRelease removes a Pokémon from the user's Pokédex.
// This command simulates releasing a caught Pokémon back into the wild,
// removing it from the user's collection.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to release
//
// Returns:
//   - An error if no Pokémon name is provided or if the Pokémon is not in the Pokédex
func commandRelease(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}

	// Process the Pokémon name and check if it exists
	apiName, exists, _ := CheckPokemonExists(cfg, params[0])
	nameInfo := FormatPokemonInput(apiName)

	if exists {
		// Lock the config before modifying the pokedex
		cfg.mutex.Lock()
		// Remove the pokemon from the pokedex
		delete(cfg.pokedex, apiName)
		cfg.mutex.Unlock()

		fmt.Printf("%s was released. Bye, %s!\n", nameInfo.Formatted, nameInfo.Formatted)
		fmt.Println("-----")

		// Auto-save after releasing a Pokémon
		if err := UpdatePokedexAndSave(cfg); err != nil {
			return err
		}

		return nil
	}

	return HandlePokemonNotFound(nameInfo.APIFormat)
}
