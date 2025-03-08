package main

import (
	"encoding/json"
	"fmt"
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
