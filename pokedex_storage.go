package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// defaultSaveFile is the default location for storing Pokédex data.
// The file is stored in the user's home directory.
const defaultSaveFile = ".pokedexcli_save.json"

// SaveData represents the structure of data saved to disk.
// It includes the Pokédex data and other persistent state.
type SaveData struct {
	Pokedex         map[string]pokeapi.PokemonDataResp `json:"pokedex"`         // User's caught Pokémon
	NextLocationURL *string                            `json:"nextLocationURL"` // URL for next location page
	PrevLocationURL *string                            `json:"prevLocationURL"` // URL for previous location page
	LastSaved       time.Time                          `json:"lastSaved"`       // Timestamp of the last save
}

// getSaveFilePath returns the full path to the save file.
// It tries to use the user's home directory, falling back to the current directory.
func getSaveFilePath() (string, error) {
	// Try to get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if home can't be determined
		return defaultSaveFile, nil
	}
	return filepath.Join(homeDir, defaultSaveFile), nil
}

// savePokedexData saves the current Pokédex and navigation state to disk.
// It serializes the data as JSON and writes it to the save file.
//
// Parameters:
//   - cfg: The application configuration containing data to save
//
// Returns:
//   - An error if the save operation fails
func savePokedexData(cfg *config) error {
	saveData := SaveData{
		Pokedex:         cfg.pokedex,
		NextLocationURL: cfg.nextLocationURL,
		PrevLocationURL: cfg.prevLocationURL,
		LastSaved:       time.Now(),
	}

	// Serialize data to JSON
	data, err := json.Marshal(saveData)
	if err != nil {
		return fmt.Errorf("error serializing Pokédex data: %w", err)
	}

	// Get save file path
	saveFilePath, err := getSaveFilePath()
	if err != nil {
		return fmt.Errorf("error determining save file path: %w", err)
	}

	// Write data to file
	err = os.WriteFile(saveFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing save file: %w", err)
	}

	return nil
}

// loadPokedexData loads saved Pokédex and navigation state from disk.
// It deserializes the JSON data from the save file and updates the configuration.
//
// Parameters:
//   - cfg: The application configuration to update with loaded data
//
// Returns:
//   - An error if the load operation fails
func loadPokedexData(cfg *config) error {
	// Get save file path
	saveFilePath, err := getSaveFilePath()
	if err != nil {
		return fmt.Errorf("error determining save file path: %w", err)
	}

	// Check if the file exists
	if _, err := os.Stat(saveFilePath); os.IsNotExist(err) {
		// No save file exists, nothing to load
		return nil
	}

	// Read data from file
	data, err := os.ReadFile(saveFilePath)
	if err != nil {
		return fmt.Errorf("error reading save file: %w", err)
	}

	// Deserialize JSON data
	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return fmt.Errorf("error deserializing Pokédex data: %w", err)
	}

	// Update configuration with loaded data
	cfg.pokedex = saveData.Pokedex
	cfg.nextLocationURL = saveData.NextLocationURL
	cfg.prevLocationURL = saveData.PrevLocationURL

	return nil
}

// commandSave explicitly saves the current Pokédex state to disk.
// This is useful when users want to ensure their progress is saved.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex data
//   - params: Command parameters (unused)
//
// Returns:
//   - An error if the save operation fails
func commandSave(cfg *config, params []string) error {
	err := savePokedexData(cfg)
	if err != nil {
		return err
	}
	fmt.Println("Pokédex saved successfully!")
	return nil
}

// commandReset clears the user's Pokédex, returning it to an empty state.
// This allows users to start fresh without existing caught Pokémon.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex to clear
//   - params: Command parameters (unused)
//
// Returns:
//   - An error if the user cancels the operation
func commandReset(cfg *config, params []string) error {
	// Confirm with the user before clearing data
	fmt.Print("Are you sure you want to clear your Pokédex? This cannot be undone. (y/N): ")
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return errors.New("operation cancelled")
	}

	// Clear the Pokédex
	cfg.pokedex = make(map[string]pokeapi.PokemonDataResp)
	fmt.Println("Pokédex cleared! All Pokémon have been released.")

	// Save the empty state
	err := savePokedexData(cfg)
	if err != nil {
		return fmt.Errorf("error saving empty Pokédex: %w", err)
	}

	return nil
}
