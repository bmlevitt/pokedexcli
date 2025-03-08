package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// GetEvolutionChainBySpecies retrieves the complete evolution chain for a specific Pokémon.
// This function first gets species data for the Pokémon, then uses the evolution chain URL
// from that data to fetch the full evolution chain. It's used by the "evolve" command
// to determine possible evolutions for a Pokémon.
//
// Parameters:
//   - pokemonName: The name of the Pokémon to get evolution chain for (in lowercase with hyphens)
//
// Returns:
//   - An EvolutionChainResp containing the complete evolution chain data
//   - An error if the API request fails, the Pokémon doesn't exist, or it has no evolution data
func (c *Client) GetEvolutionChainBySpecies(pokemonName string) (EvolutionChainResp, error) {
	// First, get the species data to find the evolution chain URL
	speciesData, err := c.GetPokemonSpecies(pokemonName)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("error fetching species data: %w", err)
	}

	// Check if evolution chain URL exists in the species data
	evolutionURL, err := getEvolutionChainURL(speciesData)
	if err != nil {
		return EvolutionChainResp{}, errorhandling.EvolutionNotFoundError(pokemonName, err)
	}

	// Extract the evolution chain ID from the URL
	parts := strings.Split(evolutionURL, "/")
	if len(parts) < 2 {
		return EvolutionChainResp{}, fmt.Errorf("invalid evolution chain URL format: %s", evolutionURL)
	}

	// The ID should be the last part of the URL (minus any trailing slash)
	idStr := parts[len(parts)-1]
	if idStr == "" && len(parts) >= 3 {
		idStr = parts[len(parts)-2] // Handle trailing slash
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("invalid evolution chain ID: %s", idStr)
	}

	// Now get the evolution chain data
	return c.GetEvolutionChain(id)
}

// getEvolutionChainURL extracts the evolution chain URL from a PokemonSpeciesResp.
// This is a helper function for GetEvolutionChainBySpecies.
//
// Parameters:
//   - speciesData: The species data containing the evolution chain URL
//
// Returns:
//   - The URL string for the evolution chain
//   - An error if the URL is missing or empty
func getEvolutionChainURL(speciesData PokemonSpeciesResp) (string, error) {
	if speciesData.EvolutionChain.URL == "" {
		return "", fmt.Errorf("evolution chain URL not found")
	}
	return speciesData.EvolutionChain.URL, nil
}

// GetEvolutionChain retrieves evolution chain data directly by its ID.
// Evolution chains describe how Pokémon can evolve from one form to another.
// Results are cached to improve performance and reduce API calls.
//
// Parameters:
//   - id: The unique identifier of the evolution chain to retrieve
//
// Returns:
//   - An EvolutionChainResp containing the complete evolution chain data
//   - An error if the API request fails or the evolution chain doesn't exist
func (c *Client) GetEvolutionChain(id int) (EvolutionChainResp, error) {
	endpoint := "/evolution-chain/"
	fullURL := baseURL + endpoint + strconv.Itoa(id)

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		evolutionChainResp := EvolutionChainResp{}
		err := json.Unmarshal(data, &evolutionChainResp)
		if err != nil {
			return EvolutionChainResp{}, fmt.Errorf("error unmarshaling cached evolution data: %w", err)
		}
		return evolutionChainResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return EvolutionChainResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return EvolutionChainResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusNotFound {
			return EvolutionChainResp{}, errorhandling.FormatResourceNotFoundError(
				errorhandling.ResourceEvolutionChain,
				fmt.Sprintf("ID: %d", id),
				fmt.Errorf("HTTP 404"))
		}
		return EvolutionChainResp{}, errorhandling.NewAPIError(resp.StatusCode, fmt.Sprintf("%s%d", endpoint, id), fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	evolutionChainResp := EvolutionChainResp{}
	err = json.Unmarshal(body, &evolutionChainResp)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return evolutionChainResp, nil
}
