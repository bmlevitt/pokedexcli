// This file provides utility functions for text formatting and display
// throughout the PokédexCLI application. It includes functions for formatting
// Pokémon names, types, moves, and other text elements for user-friendly display.
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PrettyPrint formats and displays any data structure as indented JSON.
// This utility function is useful for debugging and development purposes,
// allowing for clear visualization of complex nested data structures.
//
// The function uses JSON encoding with indentation to create a human-readable
// representation of the data, which is then printed to standard output.
//
// Parameters:
//   - data: Any Go data structure that can be marshaled to JSON
//
// Note:
// If marshaling fails, an error message is printed and the function returns without
// displaying the data.
func PrettyPrint(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

// FormatMoveName converts API move names (like "thunder-punch") to a user-friendly format (like "Thunder Punch").
//
// Parameters:
//   - name: The raw move name with hyphens
//
// Returns:
//   - A formatted move name with spaces and proper capitalization
func FormatMoveName(name string) string {
	// Replace hyphens with spaces
	name = strings.ReplaceAll(name, "-", " ")

	// Split the name into words
	words := strings.Fields(name)
	for i, word := range words {
		// Capitalize each word
		words[i] = CapitalizeFirstLetter(word)
	}

	// Join the words back together
	return strings.Join(words, " ")
}

// CapitalizeFirstLetter capitalizes the first letter of a string.
//
// Parameters:
//   - s: The input string
//
// Returns:
//   - The string with its first letter capitalized, or the original string if empty
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return cases.Title(language.English).String(s)
}

// FormatLocationName converts API location names (like "cerulean-city") to a user-friendly format (like "Cerulean City").
//
// Parameters:
//   - name: The raw location name with hyphens
//
// Returns:
//   - A formatted location name with spaces and proper capitalization
func FormatLocationName(name string) string {
	// Replace hyphens with spaces
	name = strings.ReplaceAll(name, "-", " ")

	// Split the name into words
	words := strings.Fields(name)
	for i, word := range words {
		// Capitalize each word
		words[i] = CapitalizeFirstLetter(word)
	}

	// Join the words back together
	return strings.Join(words, " ")
}

// FormatTypeName converts API type names (like "fire") to a capitalized format (like "Fire").
//
// Parameters:
//   - name: The raw type name
//
// Returns:
//   - A formatted type name with proper capitalization
func FormatTypeName(name string) string {
	return CapitalizeFirstLetter(name)
}

// FormatStatName converts API stat names (like "special-attack") to a user-friendly format (like "Special Attack").
//
// Parameters:
//   - name: The raw stat name with hyphens
//
// Returns:
//   - A formatted stat name with spaces and proper capitalization
func FormatStatName(name string) string {
	// Replace hyphens with spaces
	name = strings.ReplaceAll(name, "-", " ")

	// Split the name into words
	words := strings.Fields(name)
	for i, word := range words {
		// Capitalize each word
		words[i] = CapitalizeFirstLetter(word)
	}

	// Join the words back together
	return strings.Join(words, " ")
}

// FormatPokemonName converts API Pokémon names (like "pikachu") to a properly capitalized format (like "Pikachu").
//
// Parameters:
//   - name: The raw Pokémon name
//
// Returns:
//   - A formatted Pokémon name with proper capitalization
func FormatPokemonName(name string) string {
	if len(name) == 0 {
		return name
	}

	// Split the name by hyphens
	parts := strings.Split(name, "-")

	// Capitalize first letter of first part
	parts[0] = cases.Title(language.English).String(parts[0])

	// Keep the rest of the parts lowercase
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToLower(parts[i])
	}

	// Join the parts back with hyphens
	return strings.Join(parts, "-")
}

// ConvertToAPIFormat converts a user-friendly formatted string (like "Mt Coronet 5F")
// back to the API format (like "mt-coronet-5f").
//
// Parameters:
//   - formattedName: The formatted name with spaces and proper capitalization
//
// Returns:
//   - The name in API format with lowercase and hyphens instead of spaces
func ConvertToAPIFormat(formattedName string) string {
	// Handle special case Pokemon names
	specialCasePokemon := map[string]string{
		"mr mime":   "mr-mime",
		"mr. mime":  "mr-mime",
		"mime jr":   "mime-jr",
		"mime jr.":  "mime-jr",
		"type null": "type-null",
		"type: null": "type-null",
		"tapu koko": "tapu-koko",
		"tapu lele": "tapu-lele",
		"tapu bulu": "tapu-bulu",
		"tapu fini": "tapu-fini",
		"ho oh":    "ho-oh",
		"porygon z": "porygon-z",
		"jangmo o":  "jangmo-o",
		"hakamo o":  "hakamo-o",
		"kommo o":   "kommo-o",
	}

	// Check if this is a special case Pokemon name (case insensitive)
	lowerName := strings.ToLower(formattedName)
	if apiName, isSpecial := specialCasePokemon[lowerName]; isSpecial {
		return apiName
	}

	// Handle special cases for "é" character in "Pokémon"
	formattedName = strings.ReplaceAll(formattedName, "é", "e")

	// Remove special characters that aren't letters, numbers, or spaces
	// This handles things like parentheses, quotation marks, etc.
	var result strings.Builder
	for _, r := range formattedName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' {
			result.WriteRune(r)
		}
	}
	formattedName = result.String()

	// Convert to lowercase
	formattedName = strings.ToLower(formattedName)

	// Replace multiple spaces with a single space
	formattedName = strings.Join(strings.Fields(formattedName), " ")

	// Replace spaces with hyphens for the final API format
	return strings.ReplaceAll(formattedName, " ", "-")
}
