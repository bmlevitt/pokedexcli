package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func commandShowOff(cfg *config, params []string) error {
	// Check for pokemon name parameter
	if len(params) == 0 {
		return errors.New("no pokemon name provided")
	}
	pokemonName := params[0]

	// Check if the pokemon exists in the pokedex
	pokemon, exists := cfg.pokedex[pokemonName]
	if !exists {
		return fmt.Errorf("%s is not in your pokedex", pokemonName)
	}

	// Check if the pokemon has any moves
	if len(pokemon.Moves) == 0 {
		return fmt.Errorf("%s doesn't know any moves", pokemonName)
	}

	// Select a random move
	randomIndex := rand.Intn(len(pokemon.Moves))
	moveName := pokemon.Moves[randomIndex].Move.Name

	// Format the move name for better display (replace hyphens with spaces and capitalize words)
	formattedMove := formatMoveName(moveName)

	// Show off the pokemon using the move
	fmt.Printf("%s used %s!\n", capitalizeFirstLetter(pokemonName), formattedMove)

	return nil
}

// formatMoveName formats a move name by replacing hyphens with spaces and capitalizing each word
func formatMoveName(name string) string {
	// Replace hyphens with spaces
	spaced := strings.ReplaceAll(name, "-", " ")

	// Capitalize each word
	words := strings.Fields(spaced)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// capitalizeFirstLetter capitalizes the first letter of a string
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
