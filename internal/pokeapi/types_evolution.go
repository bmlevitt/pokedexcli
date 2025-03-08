package pokeapi

// EvolutionChainResp represents the response from the evolution-chain endpoint
type EvolutionChainResp struct {
	ID              int         `json:"id"`
	BabyTriggerItem interface{} `json:"baby_trigger_item"`
	Chain           ChainLink   `json:"chain"`
}

// ChainLink represents a link in the evolution chain
type ChainLink struct {
	IsBaby           bool              `json:"is_baby"`
	Species          NamedAPIResource  `json:"species"`
	EvolutionDetails []EvolutionDetail `json:"evolution_details"`
	EvolvesTo        []ChainLink       `json:"evolves_to"`
}

// EvolutionDetail contains the details of how a Pokémon evolves
type EvolutionDetail struct {
	Item                  interface{}      `json:"item"`
	Trigger               NamedAPIResource `json:"trigger"`
	Gender                interface{}      `json:"gender"`
	HeldItem              interface{}      `json:"held_item"`
	KnownMove             interface{}      `json:"known_move"`
	KnownMoveType         interface{}      `json:"known_move_type"`
	Location              interface{}      `json:"location"`
	MinLevel              int              `json:"min_level"`
	MinHappiness          interface{}      `json:"min_happiness"`
	MinBeauty             interface{}      `json:"min_beauty"`
	MinAffection          interface{}      `json:"min_affection"`
	NeedsOverworldRain    bool             `json:"needs_overworld_rain"`
	PartySpecies          interface{}      `json:"party_species"`
	PartyType             interface{}      `json:"party_type"`
	RelativePhysicalStats interface{}      `json:"relative_physical_stats"`
	TimeOfDay             string           `json:"time_of_day"`
	TradeSpecies          interface{}      `json:"trade_species"`
	TurnUpsideDown        bool             `json:"turn_upside_down"`
}
