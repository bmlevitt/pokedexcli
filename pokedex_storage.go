package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
	"github.com/gofrs/flock"
)

// defaultSaveFile is the default location for storing Pokédex data.
// The file is stored in the user's home directory.
const defaultSaveFile = ".pokedexcli_save.json"

// lockTimeout is the maximum time to wait for acquiring a file lock
const lockTimeout = 5 * time.Second

// lockRetryInterval is the interval between lock attempts
const lockRetryInterval = 100 * time.Millisecond

// SaveData represents the structure of data saved to disk.
// It includes the Pokédex data and other persistent state.
type SaveData struct {
	Pokedex   map[string]pokeapi.PokemonDataResp `json:"pokedex"`   // User's caught Pokémon
	LastSaved time.Time                          `json:"lastSaved"` // Timestamp of the last save
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

// getLockFilePath returns the path to the lock file based on the save file path
func getLockFilePath(saveFilePath string) string {
	return saveFilePath + ".lock"
}

// savePokedexData saves the current Pokédex and navigation state to disk.
// It serializes the data as JSON and writes it to the save file.
// The function uses file locking to prevent concurrent access issues and
// performs atomic writes by writing to a temporary file first.
//
// Parameters:
//   - cfg: The application configuration containing data to save
//
// Returns:
//   - An error if the save operation fails
func savePokedexData(cfg *config) error {
	// Get save file path
	saveFilePath, err := getSaveFilePath()
	if err != nil {
		return fmt.Errorf("error determining save file path: %w", err)
	}

	// Create a file lock
	fileLock := flock.New(getLockFilePath(saveFilePath))

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), lockTimeout)
	defer cancel()

	// Acquire an exclusive lock with a timeout
	locked, err := fileLock.TryLockContext(ctx, lockRetryInterval)
	if err != nil {
		return fmt.Errorf("error acquiring file lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("could not acquire lock on save file: timeout after %v", lockTimeout)
	}

	// Release the lock when we're done
	defer fileLock.Unlock()

	// Acquire read lock on the config to get a consistent snapshot
	cfg.mutex.RLock()
	saveData := SaveData{
		Pokedex:   cfg.pokedex,
		LastSaved: time.Now(),
	}
	cfg.mutex.RUnlock()

	// Serialize data to JSON
	data, err := json.Marshal(saveData)
	if err != nil {
		return fmt.Errorf("error serializing Pokédex data: %w", err)
	}

	// Create a temporary file in the same directory
	tempFilePath := saveFilePath + ".tmp"

	// Write data to the temporary file
	err = os.WriteFile(tempFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing temporary save file: %w", err)
	}

	// Atomically replace the old file with the new one
	err = os.Rename(tempFilePath, saveFilePath)
	if err != nil {
		// Try to clean up the temporary file if rename fails
		os.Remove(tempFilePath)
		return fmt.Errorf("error replacing save file: %w", err)
	}

	return nil
}

// loadPokedexData loads saved Pokédex and navigation state from disk.
// It deserializes the JSON data from the save file and updates the configuration.
// The function uses file locking to prevent concurrent access issues.
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

	// Create a file lock
	fileLock := flock.New(getLockFilePath(saveFilePath))

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), lockTimeout)
	defer cancel()

	// Acquire a shared lock with a timeout
	locked, err := fileLock.TryRLockContext(ctx, lockRetryInterval)
	if err != nil {
		return fmt.Errorf("error acquiring file lock for reading: %w", err)
	}
	if !locked {
		return fmt.Errorf("could not acquire read lock on save file: timeout after %v", lockTimeout)
	}

	// Release the lock when we're done
	defer fileLock.Unlock()

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

	// Update configuration with loaded data - acquire a write lock
	cfg.mutex.Lock()
	cfg.pokedex = saveData.Pokedex
	// Don't load map navigation URLs - user must run 'map' command first
	cfg.nextLocationURL = nil
	cfg.prevLocationURL = nil
	cfg.mapViewedThisSession = false
	cfg.mutex.Unlock()

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
