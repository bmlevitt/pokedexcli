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
			name:  "Name with spaces (Mr Mime)",
			input: "mr mime",
			expected: PokemonNameInfo{
				Input:     "mr mime",
				Formatted: "Mr-mime",
				APIFormat: "mr-mime",
			},
		},
		{
			name:  "Name with period (Mr. Mime)",
			input: "mr. mime",
			expected: PokemonNameInfo{
				Input:     "mr. mime",
				Formatted: "Mr-mime",
				APIFormat: "mr-mime",
			},
		},
		{
			name:  "Name with period and junior (Mime Jr.)",
			input: "mime jr.",
			expected: PokemonNameInfo{
				Input:     "mime jr.",
				Formatted: "Mime-jr",
				APIFormat: "mime-jr",
			},
		},
		{
			name:  "Name with colon (Type: Null)",
			input: "type: null",
			expected: PokemonNameInfo{
				Input:     "type: null",
				Formatted: "Type-null",
				APIFormat: "type-null",
			},
		},
		{
			name:  "Name with hyphen (Ho-oh)",
			input: "ho-oh",
			expected: PokemonNameInfo{
				Input:     "ho-oh",
				Formatted: "Ho-oh",
				APIFormat: "ho-oh",
			},
		},
		{
			name:  "Name with spaces that has a hyphen variant (Ho oh)",
			input: "ho oh",
			expected: PokemonNameInfo{
				Input:     "ho oh",
				Formatted: "Ho-oh",
				APIFormat: "ho-oh",
			},
		},
		{
			name:  "Multi-part name (Tapu Koko)",
			input: "tapu koko",
			expected: PokemonNameInfo{
				Input:     "tapu koko",
				Formatted: "Tapu-koko",
				APIFormat: "tapu-koko",
			},
		},
		{
			name:  "Name with Z (Porygon Z)",
			input: "porygon z",
			expected: PokemonNameInfo{
				Input:     "porygon z",
				Formatted: "Porygon-z",
				APIFormat: "porygon-z",
			},
		},
		{
			name:  "Name with trailing o (Jangmo-o)",
			input: "jangmo o",
			expected: PokemonNameInfo{
				Input:     "jangmo o",
				Formatted: "Jangmo-o",
				APIFormat: "jangmo-o",
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
