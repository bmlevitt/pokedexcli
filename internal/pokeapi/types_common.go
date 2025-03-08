package pokeapi

// NamedAPIResource represents a resource with a name and URL
// This type is used widely throughout the PokeAPI to represent resources
type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PaginatedResponse is a common structure for paginated responses
type PaginatedResponse struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}
