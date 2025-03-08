package main

import (
	"fmt"
	"strconv"
)

// commandAutoSave controls the auto-save feature of the application.
// It allows users to enable or disable automatic saving of their Pokédex
// after changes (catching, releasing, or evolving Pokémon).
//
// Parameters:
//   - cfg: The application configuration
//   - params: Command parameters, where params[0] is either "on", "off", or omitted
//
// Returns:
//   - An error if the parameter is invalid
func commandAutoSave(cfg *config, params []string) error {
	// If no parameter is provided, display the current status
	if len(params) == 0 {
		status := "enabled"
		if !cfg.autoSaveEnabled {
			status = "disabled"
		}
		fmt.Printf("Auto-save is currently %s\n", status)
		return nil
	}

	// Otherwise, update the setting based on the provided parameter
	switch params[0] {
	case "on", "true", "1", "enable", "enabled":
		cfg.autoSaveEnabled = true
		fmt.Println("Auto-save enabled. Your Pokédex will be saved automatically after changes.")
	case "off", "false", "0", "disable", "disabled":
		cfg.autoSaveEnabled = false
		fmt.Println("Auto-save disabled. Use 'save' command to manually save your Pokédex.")
	default:
		return fmt.Errorf("invalid parameter: %s (use 'on' or 'off')", params[0])
	}

	// Save the configuration itself, including the new autosave setting
	return savePokedexData(cfg)
}

// autoSaveIfEnabled automatically saves the Pokédex data if auto-save is enabled.
// This is a helper function to be called after operations that modify the Pokédex.
//
// Parameters:
//   - cfg: The application configuration containing the Pokédex data
//
// Returns:
//   - An error if the save operation fails
func autoSaveIfEnabled(cfg *config) error {
	if cfg.autoSaveEnabled {
		return savePokedexData(cfg)
	}
	return nil
}

// commandSaveInterval sets or displays the interval for auto-saving.
// This allows users to configure how frequently their Pokédex is saved
// when changes are made, in terms of number of operations.
//
// Parameters:
//   - cfg: The application configuration
//   - params: Command parameters, where params[0] is the interval or omitted
//
// Returns:
//   - An error if the parameter is invalid
func commandSaveInterval(cfg *config, params []string) error {
	// If no parameter is provided, display the current interval
	if len(params) == 0 {
		if cfg.autoSaveInterval == 1 {
			fmt.Println("Auto-save occurs after every change to your Pokédex.")
		} else {
			fmt.Printf("Auto-save occurs after every %d changes to your Pokédex.\n", cfg.autoSaveInterval)
		}
		return nil
	}

	// Parse the provided interval
	interval, err := strconv.Atoi(params[0])
	if err != nil || interval < 1 {
		return fmt.Errorf("invalid interval: %s (must be a positive number)", params[0])
	}

	// Update the interval
	cfg.autoSaveInterval = interval

	// Provide feedback
	if interval == 1 {
		fmt.Println("Auto-save will occur after every change to your Pokédex.")
	} else {
		fmt.Printf("Auto-save will occur after every %d changes to your Pokédex.\n", interval)
	}

	// Save the configuration itself, including the new interval setting
	return savePokedexData(cfg)
}
