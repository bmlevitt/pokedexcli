// This file defines the data structures for working with Pokémon-related data from the PokeAPI.
// It includes type definitions for Pokémon details, capture rates, and species information.
package pokeapi

// PokemonDataResp represents the response from the pokemon endpoint in the PokeAPI.
// It contains comprehensive data about a Pokémon, including its attributes, stats,
// types, moves, and species information. This is the primary data structure used
// when displaying Pokémon information in the application.
type PokemonDataResp struct {
	// Basic information
	Name   string `json:"name"`   // The name of the Pokémon (lowercase with hyphens)
	Height int    `json:"height"` // The height of the Pokémon in decimeters
	Weight int    `json:"weight"` // The weight of the Pokémon in hectograms

	// Stats information
	Stats []struct {
		BaseStat int              `json:"base_stat"` // The base value for the stat
		Effort   int              `json:"effort"`    // The effort points (EVs) gained
		Stat     NamedAPIResource `json:"stat"`      // The stat that these values belong to
	} `json:"stats"`

	// Type information
	Types []struct {
		Slot int              `json:"slot"` // The slot that this type occupies in this Pokémon (1 or 2)
		Type NamedAPIResource `json:"type"` // The type the Pokémon has
	} `json:"types"`

	// Moves information
	Moves []struct {
		Move NamedAPIResource `json:"move"` // The move that can be learned
	} `json:"moves"`

	// Species reference
	Species NamedAPIResource `json:"species"` // The species this Pokémon belongs to

	CaptureRate int `json:"capture_rate"` // The capture rate (not in the standard API response, added manually)
}

// PokemonCaptureRateResp represents a specialized response containing just the capture rate.
// This is used by the catch command to determine the probability of successfully
// catching a Pokémon by comparing the capture rate against a random number.
//
// In the original Pokémon games, capture rates range from 0-255, with higher values
// meaning the Pokémon is easier to catch. This struct is populated based on data
// from the Pokémon species endpoint.
type PokemonCaptureRateResp struct {
	CaptureRate int `json:"capture_rate"` // The base capture rate between 0-255 (higher = easier to catch)
}

// PokemonSpeciesResp represents the response from the pokemon-species endpoint.
// It contains data about a Pokémon species, including Pokédex entries (flavor text),
// genus information, and evolution chain references. This is used by the describe
// and evolve commands.
type PokemonSpeciesResp struct {
	// Basic information
	ID   int    `json:"id"`   // The identifier for this Pokémon species
	Name string `json:"name"` // The name of this Pokémon species (lowercase with hyphens)

	// Flavor text entries from different games
	FlavorTextEntries []struct {
		FlavorText string           `json:"flavor_text"` // The localized flavor text for this species in different games
		Language   NamedAPIResource `json:"language"`    // The language this flavor text is in
		Version    NamedAPIResource `json:"version"`     // The game version this flavor text is from
	} `json:"flavor_text_entries"`

	// Form descriptions
	FormDescriptions []struct {
		Description string           `json:"description"` // The localized description of this form
		Language    NamedAPIResource `json:"language"`    // The language this description is in
	} `json:"form_descriptions"`

	// Genus information (e.g., "Mouse Pokémon" for Pikachu)
	Genera []struct {
		Genus    string           `json:"genus"`    // The localized genus of this Pokémon species
		Language NamedAPIResource `json:"language"` // The language this genus is in
	} `json:"genera"`

	// Evolution chain reference
	EvolutionChain struct {
		URL string `json:"url"` // The URL to the evolution chain data for this species
	} `json:"evolution_chain"`

	// Reference to the Pokémon species that evolves into this one
	EvolvesFromSpecies *NamedAPIResource `json:"evolves_from_species"` // The species that evolves into this one, if any
}
