package main

import (
	"fmt"
	"math/rand"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// commandCatch attempts to catch a specified Pokémon and add it to the user's Pokédex.
// The function simulates the catching mechanic from the Pokémon games by using
// the Pokémon's capture rate from the API to determine catch probability.
//
// The catch probability is calculated by comparing a random number (0-255)
// against the Pokémon's capture rate. If the random number is less than the
// capture rate, the catch is successful.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex and API client
//   - params: Command parameters where params[0] is the Pokémon name to catch
//
// Returns:
//   - An error if:
//   - No Pokémon name is provided (InvalidParameterError)
//   - The specified Pokémon doesn't exist (InvalidPokemonNameError)
//   - There's an API connection issue (NetworkError)
//   - The API response cannot be processed (InternalError)
//   - Returns nil on successful execution, even if the catch attempt fails
func commandCatch(cfg *config, params []string) error {
	// Validate the Pokemon parameter
	pokemonName, err := ValidatePokemonParam(params)
	if err != nil {
		// Use standardized error handling
		if HandleCommandError(cfg, "catch", err) {
			return err
		}
		return nil
	}

	// Process the Pokémon name input
	nameInfo := FormatPokemonInput(pokemonName)

	// Fetch pokemon capture rate
	resp, err := cfg.pokeapiClient.GetPokemonCaptureRate(nameInfo.APIFormat)
	if err != nil {
		// Check if this is an invalid Pokémon name (doesn't exist) error
		if errorhandling.IsNotFoundError(err) {
			// Convert to our standard invalid Pokémon name error
			invalidNameErr := errorhandling.InvalidPokemonNameError(nameInfo.Formatted)
			return invalidNameErr
		}

		// Use standardized error handling for other errors
		if HandleCommandError(cfg, "catch", err) {
			return err
		}
		return nil
	}

	captureRate := resp.CaptureRate

	// Scale the capture rate for rare Pokémon
	// Pokémon with capture rates below 50 are considered rare
	// The scaling ensures rare Pokémon (capture rate < 50) have at least a 10-20% catch rate
	// This means they should be caught within 5-10 attempts on average
	effectiveCaptureRate := captureRate
	isRare := captureRate < 50
	if isRare {
		// For rare Pokémon, boost the capture rate to be between 25-50
		// This gives approximately a 10-20% chance per throw
		effectiveCaptureRate = captureRate + (50-captureRate)/2
	}

	randNum := rand.Intn(256)
	caught := randNum < effectiveCaptureRate

	// Determine what ball to use in the message
	if isRare && caught {
		fmt.Println("You found a Masterball lying nearby...!")
		fmt.Printf("Throwing a Masterball at %s...\n", nameInfo.Formatted)
	} else {
		fmt.Printf("Throwing a Pokéball at %s...\n", nameInfo.Formatted)
	}

	if caught {
		pokeData, err := cfg.pokeapiClient.GetPokemonData(nameInfo.APIFormat)
		if err != nil {
			// Use standardized error handling
			if HandleCommandError(cfg, "catch", err) {
				return err
			}
			return nil
		}

		// Lock the config before modifying the pokedex
		cfg.mutex.Lock()
		cfg.pokedex[nameInfo.APIFormat] = pokeData
		cfg.mutex.Unlock()

		fmt.Printf("%s was caught!\n", nameInfo.Formatted)

		// Auto-save after catching a Pokémon
		if err := UpdatePokedexAndSave(cfg); err != nil {
			// Use standardized error handling but don't return the error
			// since we still want to show the success message
			HandleCommandError(cfg, "catch", err)
		}
	} else {
		fmt.Printf("%s escaped!\n", nameInfo.Formatted)
	}
	fmt.Println("-----")
	return nil
}
