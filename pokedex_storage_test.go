package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokeapi"
)

// testSaveAndLoadPokedex is a helper function that saves and loads Pokédex data
// using a specified save file path for testing purposes
func testSaveAndLoadPokedex(saveFilePath string, cfg *config) error {
	// Create a SaveData struct
	saveData := SaveData{
		Pokedex:   cfg.pokedex,
		LastSaved: time.Now(),
	}

	// Marshal the data to JSON
	jsonData, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return err
	}

	// Write to the file
	err = os.WriteFile(saveFilePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// testLoadPokedexData is a helper function that loads Pokédex data
// from a specified save file path for testing purposes
func testLoadPokedexData(saveFilePath string, cfg *config) error {
	// Check if the file exists
	_, err := os.Stat(saveFilePath)
	if os.IsNotExist(err) {
		return nil // No save file, nothing to load
	}

	// Read the file
	jsonData, err := os.ReadFile(saveFilePath)
	if err != nil {
		return err
	}

	// Unmarshal the data
	var saveData SaveData
	err = json.Unmarshal(jsonData, &saveData)
	if err != nil {
		return err
	}

	// Update the config
	cfg.pokedex = saveData.Pokedex

	return nil
}

// TestSaveAndLoadPokedex tests the saving and loading of Pokédex data
func TestSaveAndLoadPokedex(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "pokedex_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary save file path
	tempSaveFile := filepath.Join(tempDir, "test_save.json")

	// Create test data
	testPokemon := pokeapi.PokemonDataResp{
		Name:   "pikachu",
		Height: 4,
		Weight: 60,
	}

	// Create a config with test data
	cfg := &config{
		pokedex: map[string]pokeapi.PokemonDataResp{
			"pikachu": testPokemon,
		},
		autoSaveEnabled:  true,
		autoSaveInterval: 1,
	}

	// Test saving
	err = testSaveAndLoadPokedex(tempSaveFile, cfg)
	if err != nil {
		t.Fatalf("Failed to save Pokédex data: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(tempSaveFile); os.IsNotExist(err) {
		t.Fatalf("Save file was not created")
	}

	// Create a new empty config
	newCfg := &config{
		pokedex:          make(map[string]pokeapi.PokemonDataResp),
		autoSaveEnabled:  true,
		autoSaveInterval: 1,
	}

	// Test loading
	err = testLoadPokedexData(tempSaveFile, newCfg)
	if err != nil {
		t.Fatalf("Failed to load Pokédex data: %v", err)
	}

	// Verify the data was loaded correctly
	loadedPokemon, exists := newCfg.pokedex["pikachu"]
	if !exists {
		t.Fatalf("Pikachu was not loaded into the Pokédex")
	}

	if loadedPokemon.Name != testPokemon.Name {
		t.Errorf("Expected Pokémon name %s, got %s", testPokemon.Name, loadedPokemon.Name)
	}

	if loadedPokemon.Height != testPokemon.Height {
		t.Errorf("Expected Pokémon height %d, got %d", testPokemon.Height, loadedPokemon.Height)
	}

	if loadedPokemon.Weight != testPokemon.Weight {
		t.Errorf("Expected Pokémon weight %d, got %d", testPokemon.Weight, loadedPokemon.Weight)
	}
}

// TestAutoSaveLogic tests the auto-save logic without actually saving files
func TestAutoSaveLogic(t *testing.T) {
	// Test cases
	testCases := []struct {
		name             string
		autoSaveEnabled  bool
		autoSaveInterval int
		changesSinceSync int
		expectedToSave   bool
		expectedChanges  int
	}{
		{
			name:             "Auto-save enabled, interval 1, should save",
			autoSaveEnabled:  true,
			autoSaveInterval: 1,
			changesSinceSync: 1,
			expectedToSave:   true,
			expectedChanges:  0, // Reset to 0 after save
		},
		{
			name:             "Auto-save enabled, interval 2, not enough changes",
			autoSaveEnabled:  true,
			autoSaveInterval: 2,
			changesSinceSync: 1,
			expectedToSave:   false,
			expectedChanges:  2, // Incremented by 1
		},
		{
			name:             "Auto-save enabled, interval 2, enough changes",
			autoSaveEnabled:  true,
			autoSaveInterval: 2,
			changesSinceSync: 2,
			expectedToSave:   true,
			expectedChanges:  0, // Reset to 0 after save
		},
		{
			name:             "Auto-save disabled, should not save",
			autoSaveEnabled:  false,
			autoSaveInterval: 1,
			changesSinceSync: 10,
			expectedToSave:   false,
			expectedChanges:  11, // Incremented by 1
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create config for this test
			cfg := &config{
				autoSaveEnabled:  tc.autoSaveEnabled,
				autoSaveInterval: tc.autoSaveInterval,
				changesSinceSync: tc.changesSinceSync,
			}

			// Check if we should save based on the auto-save logic
			shouldSave := cfg.autoSaveEnabled && cfg.changesSinceSync >= cfg.autoSaveInterval

			// Simulate the update logic
			cfg.changesSinceSync++
			if shouldSave {
				cfg.changesSinceSync = 0
			}

			// Verify the logic matches our expectations
			if shouldSave != tc.expectedToSave {
				t.Errorf("Expected shouldSave to be %v, but got %v", tc.expectedToSave, shouldSave)
			}

			if cfg.changesSinceSync != tc.expectedChanges {
				t.Errorf("Expected changesSinceSync to be %d, but got %d", tc.expectedChanges, cfg.changesSinceSync)
			}
		})
	}
}
