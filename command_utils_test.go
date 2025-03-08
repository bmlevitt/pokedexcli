package main

import (
	"testing"
)

// TestValidatePokemonParam tests the validation of Pokémon parameters
func TestValidatePokemonParam(t *testing.T) {
	cases := []struct {
		name          string
		params        []string
		expected      string
		expectedError bool
	}{
		{
			name:          "With valid pokemon name",
			params:        []string{"pikachu"},
			expected:      "pikachu",
			expectedError: false,
		},
		{
			name:          "With empty params",
			params:        []string{},
			expected:      "",
			expectedError: true,
		},
		{
			name:          "With empty pokemon name",
			params:        []string{""},
			expected:      "",
			expectedError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePokemonParam(tc.params)
			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected result %v but got %v", tc.expected, result)
			}
		})
	}
}

// TestFormatPokemonInput tests the FormatPokemonInput function
func TestFormatPokemonInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected PokemonNameInfo
	}{
		{
			name:  "Simple name",
			input: "pikachu",
			expected: PokemonNameInfo{
				Input:     "pikachu",
				Formatted: "Pikachu",
				APIFormat: "pikachu",
			},
		},
		{
			name:  "Name with spaces",
			input: "mr mime",
			expected: PokemonNameInfo{
				Input:     "mr mime",
				Formatted: "Mr-mime",
				APIFormat: "mr-mime",
			},
		},
		{
			name:  "Name with hyphen",
			input: "ho-oh",
			expected: PokemonNameInfo{
				Input:     "ho-oh",
				Formatted: "Ho-oh",
				APIFormat: "ho-oh",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatPokemonInput(tc.input)

			if result.Input != tc.expected.Input {
				t.Errorf("Expected Input %v but got %v", tc.expected.Input, result.Input)
			}

			if result.Formatted != tc.expected.Formatted {
				t.Errorf("Expected Formatted %v but got %v", tc.expected.Formatted, result.Formatted)
			}

			if result.APIFormat != tc.expected.APIFormat {
				t.Errorf("Expected APIFormat %v but got %v", tc.expected.APIFormat, result.APIFormat)
			}
		})
	}
}
