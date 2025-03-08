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
