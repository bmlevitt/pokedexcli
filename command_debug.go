package main

import "fmt"

// commandToggleDebug toggles the debug mode setting in the application.
// When debug mode is enabled, detailed error information is logged to stderr,
// which can be helpful for troubleshooting issues.
//
// Parameters:
//   - cfg: The application configuration containing the debug mode setting
//   - params: Not used for this command
//
// Returns:
//   - Always returns nil as this command cannot fail
func commandToggleDebug(cfg *config, params []string) error {
	// Toggle the debug mode setting
	cfg.debugMode = !cfg.debugMode

	// Display the new debug mode status
	if cfg.debugMode {
		fmt.Println("Debug mode is now enabled. Detailed error information will be logged.")
	} else {
		fmt.Println("Debug mode is now disabled. Only user-friendly error messages will be shown.")
	}
	fmt.Println("-----")

	return nil
}
