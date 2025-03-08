package pokeapi

// LocationAreasResp represents the response from the location-area endpoint
type LocationAreasResp struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

// LocationExploreResp represents the response when exploring a specific location area
type LocationExploreResp struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter represents a pokemon that can be encountered in a location area
type PokemonEncounter struct {
	Pokemon NamedAPIResource `json:"pokemon"`
}
