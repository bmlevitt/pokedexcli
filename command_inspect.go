package main

import (
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
//   - An error if no Pokémon name is provided or if the Pokémon is not in the Pokédex
func commandInspect(cfg *config, params []string) error {
	// Use the utility function to validate the Pokemon parameter and check if it exists
	_, nameInfo, pokemonData, _, err := GetPokemonIfExists(cfg, params)
	if err != nil {
		return err
	}

	// The Pokemon exists, so convert to the typed data structure
	data, err := GetTypedPokemonData(pokemonData, nameInfo.Formatted)
	if err != nil {
		return err
	}

	// Display Pokemon information
	fmt.Printf("Name: %s\n", nameInfo.Formatted)
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

	return nil
}
