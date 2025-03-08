// This file defines the data structures for working with location data from the PokeAPI.
// These structures model the geographic areas in the Pokémon world and the Pokémon
// that can be encountered in each location.
package pokeapi

// LocationAreasResp represents the response from the location-area endpoint in the PokeAPI.
// It contains a paginated list of location areas and navigation URLs for accessing
// additional pages of results. This is used by the map, next, and prev commands.
type LocationAreasResp struct {
	Count    int                `json:"count"`    // The total number of location areas available in the API
	Next     *string            `json:"next"`     // URL to the next page of results, or null if this is the last page
	Previous *string            `json:"previous"` // URL to the previous page of results, or null if this is the first page
	Results  []NamedAPIResource `json:"results"`  // The list of location areas on this page
}

// LocationExploreResp represents the response when exploring a specific location area.
// It contains information about the Pokémon that can be encountered at that location.
// This is used by the explore command to show available Pokémon at a location.
type LocationExploreResp struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"` // The list of Pokémon that can be encountered
}

// PokemonEncounter represents a Pokémon that can be encountered in a location area.
// It contains a reference to the Pokémon species that can be found at that location.
type PokemonEncounter struct {
	Pokemon NamedAPIResource `json:"pokemon"` // Reference to the Pokémon that can be encountered
}
