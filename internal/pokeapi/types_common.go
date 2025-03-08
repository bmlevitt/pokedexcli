package pokeapi

// NamedAPIResource represents a resource with a name and URL in the PokeAPI.
// This type is used extensively throughout the API to reference other objects
// like Pokémon, moves, types, locations, etc. It serves as a pointer to another
// resource that can be fetched with an additional API call.
type NamedAPIResource struct {
	Name string `json:"name"` // The name of the referenced resource (in lowercase with hyphens)
	URL  string `json:"url"`  // The URL to fetch the complete data for the referenced resource
}
