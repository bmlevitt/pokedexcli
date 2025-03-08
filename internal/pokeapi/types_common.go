// This file defines common data structures used throughout the pokeapi package.
// These type definitions represent standard response formats and resources that
// appear in multiple API endpoints in the PokeAPI.
package pokeapi

// NamedAPIResource represents a resource with a name and URL in the PokeAPI.
// This type is used extensively throughout the API to reference other objects
// like Pokémon, moves, types, locations, etc. It serves as a pointer to another
// resource that can be fetched with an additional API call.
type NamedAPIResource struct {
	Name string `json:"name"` // The name of the referenced resource (in lowercase with hyphens)
	URL  string `json:"url"`  // The URL to fetch the complete data for the referenced resource
}

// PaginatedResponse is a common structure for paginated responses in the PokeAPI.
// Many endpoints that return lists of resources use this format to allow for
// pagination through large result sets.
//
// This structure is used by multiple API endpoints including:
// - /location-area (for mapping locations)
// - /pokemon (for listing all available Pokémon)
// - /type (for listing all available types)
// - Many other list-based endpoints
//
// The client code typically uses the Next and Previous fields to implement
// navigation between pages of results.
type PaginatedResponse struct {
	Count    int                `json:"count"`    // The total number of resources available
	Next     *string            `json:"next"`     // URL to the next page of results, or null if this is the last page
	Previous *string            `json:"previous"` // URL to the previous page of results, or null if this is the first page
	Results  []NamedAPIResource `json:"results"`  // The list of resources on this page
}
