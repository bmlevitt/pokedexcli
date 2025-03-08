package main

import (
	"errors"
	"fmt"
)

// commandInspect displays detailed information about a Pokémon in the user's Pokédex.
// This command shows various attributes of a caught Pokémon, including:
//   - Base stats (HP, Attack, Defense, etc.)
//   - Physical attributes (Height and Weight)
//   - Types (Fire, Water, etc.)
//
// The information is only available for Pokémon that have been caught and are
// currently in the user's Pokédex.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex
//   - params: Command parameters where params[0] is the Pokémon name to inspect
//
// Returns:
//   - An error if no Pokémon name is provided
func commandInspect(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	inputName := params[0]

	// Convert the input name to API format if it's in a formatted style
	apiPokemonName := ConvertToAPIFormat(inputName)
	formattedName := FormatPokemonName(apiPokemonName)

	data, exists := cfg.pokedex[apiPokemonName]
	if exists {
		fmt.Printf("Name: %s\n", formattedName)
		fmt.Printf("Height: %v\n", data.Height)
		fmt.Printf("Weight: %v\n", data.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range data.Stats {
			formattedStat := FormatStatName(stat.Stat.Name)
			fmt.Printf(" - %s: %v\n", formattedStat, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, typ := range data.Types {
			formattedType := FormatTypeName(typ.Type.Name)
			fmt.Printf(" - %s\n", formattedType)
		}
		fmt.Println("-----")

	} else {
		// Check if it's a capitalization issue by trying all keys in lowercase
		found := false
		for key := range cfg.pokedex {
			if ConvertToAPIFormat(key) == apiPokemonName {
				data = cfg.pokedex[key]
				found = true
				break
			}
		}

		if found {
			fmt.Printf("Name: %s\n", formattedName)
			fmt.Printf("Height: %v\n", data.Height)
			fmt.Printf("Weight: %v\n", data.Weight)
			fmt.Printf("Stats:\n")
			for _, stat := range data.Stats {
				formattedStat := FormatStatName(stat.Stat.Name)
				fmt.Printf(" - %s: %v\n", formattedStat, stat.BaseStat)
			}
			fmt.Printf("Types:\n")
			for _, typ := range data.Types {
				formattedType := FormatTypeName(typ.Type.Name)
				fmt.Printf(" - %s\n", formattedType)
			}
			fmt.Println("-----")
		} else {
			fmt.Printf("%s has not been caught yet\n", formattedName)
		}
	}
	return nil
}
