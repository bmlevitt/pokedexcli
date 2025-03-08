package pokeapi

// PokemonDataResp represents the response from the pokemon endpoint
type PokemonDataResp struct {
	// Basic information
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`

	// Stats information
	Stats []struct {
		BaseStat int              `json:"base_stat"`
		Effort   int              `json:"effort"`
		Stat     NamedAPIResource `json:"stat"`
	} `json:"stats"`

	// Type information
	Types []struct {
		Slot int              `json:"slot"`
		Type NamedAPIResource `json:"type"`
	} `json:"types"`

	// Moves information
	Moves []struct {
		Move NamedAPIResource `json:"move"`
	} `json:"moves"`

	// Species reference
	Species NamedAPIResource `json:"species"`

	CaptureRate int `json:"capture_rate"`
}

// PokemonCaptureRateResp represents the response for a pokemon's capture rate
type PokemonCaptureRateResp struct {
	CaptureRate int `json:"capture_rate"`
}

// PokemonSpeciesResp represents the response from the pokemon-species endpoint
type PokemonSpeciesResp struct {
	// Basic information
	ID   int    `json:"id"`
	Name string `json:"name"`

	// Flavor text entries from different games
	FlavorTextEntries []struct {
		FlavorText string           `json:"flavor_text"`
		Language   NamedAPIResource `json:"language"`
		Version    NamedAPIResource `json:"version"`
	} `json:"flavor_text_entries"`

	// Form descriptions
	FormDescriptions []struct {
		Description string           `json:"description"`
		Language    NamedAPIResource `json:"language"`
	} `json:"form_descriptions"`

	// Genus information (e.g., "Mouse Pokémon" for Pikachu)
	Genera []struct {
		Genus    string           `json:"genus"`
		Language NamedAPIResource `json:"language"`
	} `json:"genera"`

	// Evolution chain reference
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`

	// Reference to the Pokémon species that evolves into this one
	EvolvesFromSpecies *NamedAPIResource `json:"evolves_from_species"`
}
