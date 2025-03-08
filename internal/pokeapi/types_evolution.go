// This file defines the data structures for working with Pokémon evolution chain data
// from the PokeAPI. These structures model how Pokémon can evolve from one form to another
// and what conditions are required for evolution to occur.
package pokeapi

// EvolutionChainResp represents the response from the evolution-chain endpoint in the PokeAPI.
// It contains the complete evolution lineage of a Pokémon species, from its basic form
// through all possible evolutions. This is used by the evolve command to determine
// the next evolutions for a Pokémon.
type EvolutionChainResp struct {
	ID              int         `json:"id"`                // The identifier for this evolution chain
	BabyTriggerItem interface{} `json:"baby_trigger_item"` // The item that triggers a baby Pokémon
	Chain           ChainLink   `json:"chain"`             // The base chain link that starts this evolution chain
}

// ChainLink represents a link in the evolution chain - a single Pokémon in the evolution line.
// Each link contains the Pokémon species data, evolution details, and links to what it can evolve into.
// The structure is recursive, with each Pokémon potentially having multiple evolution paths.
type ChainLink struct {
	IsBaby           bool              `json:"is_baby"`           // Whether this is a baby Pokémon
	Species          NamedAPIResource  `json:"species"`           // The Pokémon species at this stage in the evolution chain
	EvolutionDetails []EvolutionDetail `json:"evolution_details"` // The details for the evolution to this Pokémon
	EvolvesTo        []ChainLink       `json:"evolves_to"`        // The Pokémon species that evolve from this one
}

// EvolutionDetail contains the conditions required for a Pokémon to evolve.
// There are many different ways Pokémon can evolve in the games (level up, trading,
// using specific items, etc.), and this structure captures all those possibilities.
type EvolutionDetail struct {
	Item                  interface{}      `json:"item"`                    // The item required to trigger evolution
	Trigger               NamedAPIResource `json:"trigger"`                 // The evolution trigger (e.g., level-up, trade)
	Gender                interface{}      `json:"gender"`                  // The gender the Pokémon must be
	HeldItem              interface{}      `json:"held_item"`               // The item the Pokémon must be holding
	KnownMove             interface{}      `json:"known_move"`              // The move that must be known
	KnownMoveType         interface{}      `json:"known_move_type"`         // The type of move that must be known
	Location              interface{}      `json:"location"`                // The location where evolution must occur
	MinLevel              int              `json:"min_level"`               // The minimum level required
	MinHappiness          interface{}      `json:"min_happiness"`           // The minimum happiness required
	MinBeauty             interface{}      `json:"min_beauty"`              // The minimum beauty required
	MinAffection          interface{}      `json:"min_affection"`           // The minimum affection required
	NeedsOverworldRain    bool             `json:"needs_overworld_rain"`    // Whether it must be raining
	PartySpecies          interface{}      `json:"party_species"`           // The species that must be in the party
	PartyType             interface{}      `json:"party_type"`              // The type that must be in the party
	RelativePhysicalStats interface{}      `json:"relative_physical_stats"` // The relative physical stats (attack vs defense)
	TimeOfDay             string           `json:"time_of_day"`             // The time of day (day or night)
	TradeSpecies          interface{}      `json:"trade_species"`           // The species that must be traded
	TurnUpsideDown        bool             `json:"turn_upside_down"`        // Whether the 3DS must be turned upside-down
}
