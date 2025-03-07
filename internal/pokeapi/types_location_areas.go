package pokeapi

type LocationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationExploreResp struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}
type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonDataResp struct {
	// Basic information
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`

	// Stats information
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`

	// Type information
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`

	// Moves information
	Moves []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`

	// Species reference
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`

	CaptureRate int `json:"capture_rate"`
}

type PokemonCaptureRateResp struct {
	CaptureRate int `json:"capture_rate"`
}
