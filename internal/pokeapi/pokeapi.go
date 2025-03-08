// Package pokeapi provides a client for interacting with the PokeAPI service.
// It handles making HTTP requests, caching responses, and provides typed
// structures for working with Pokémon data.
package pokeapi

import (
	"net/http"
	"time"

	"github.com/bmlevitt/pokedexcli/internal/pokecache"
)

// baseURL is the root endpoint for the PokeAPI v2 service
const baseURL = "https://pokeapi.co/api/v2"

// Client represents a PokeAPI client that handles API requests with caching.
// It uses an internal cache to reduce the number of HTTP requests made to the API,
// improving performance and reducing load on the API service.
type Client struct {
	cache      pokecache.Cache // Cache for storing API responses
	httpClient http.Client     // HTTP client for making API requests
}

// NewClient creates a new PokeAPI client with the specified cache duration.
// The cache helps avoid redundant API calls by storing responses for the specified duration.
//
// Parameters:
//   - cacheInterval: How long cached items should remain valid before expiring
//
// Returns:
//   - A configured Client ready to make API requests with caching
func NewClient(cacheInterval time.Duration) Client {
	return Client{
		cache: *pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: time.Minute, // Set a 1-minute timeout for all requests
		},
	}
}
