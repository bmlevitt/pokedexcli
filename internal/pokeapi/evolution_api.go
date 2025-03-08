package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// GetEvolutionChainBySpecies gets the evolution chain for a Pokemon species
func (c *Client) GetEvolutionChainBySpecies(pokemonName string) (EvolutionChainResp, error) {
	// First, get the species data to find the evolution chain URL
	speciesData, err := c.GetPokemonSpecies(pokemonName)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("error fetching species data: %v", err)
	}

	// Check if evolution chain URL exists in the species data
	evolutionURL, err := getEvolutionChainURL(speciesData)
	if err != nil {
		return EvolutionChainResp{}, err
	}

	// Extract the evolution chain ID from the URL
	parts := strings.Split(evolutionURL, "/")
	if len(parts) < 2 {
		return EvolutionChainResp{}, fmt.Errorf("invalid evolution chain URL format")
	}

	idStr := parts[len(parts)-2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return EvolutionChainResp{}, fmt.Errorf("invalid evolution chain ID: %v", err)
	}

	// Get the evolution chain data
	return c.GetEvolutionChain(id)
}

// Helper function to extract the evolution chain URL from the species data
func getEvolutionChainURL(speciesData PokemonSpeciesResp) (string, error) {
	if speciesData.EvolutionChain.URL == "" {
		return "", fmt.Errorf("evolution chain URL not found")
	}
	return speciesData.EvolutionChain.URL, nil
}

// GetEvolutionChain gets the evolution chain for a given ID
func (c *Client) GetEvolutionChain(id int) (EvolutionChainResp, error) {
	endpoint := "/evolution-chain/"
	fullURL := baseURL + endpoint + strconv.Itoa(id)

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		evolutionChainResp := EvolutionChainResp{}
		err := json.Unmarshal(data, &evolutionChainResp)
		if err != nil {
			return EvolutionChainResp{}, err
		}
		return evolutionChainResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return EvolutionChainResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return EvolutionChainResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return EvolutionChainResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EvolutionChainResp{}, err
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	evolutionChainResp := EvolutionChainResp{}
	err = json.Unmarshal(body, &evolutionChainResp)
	if err != nil {
		return EvolutionChainResp{}, err
	}

	return evolutionChainResp, nil
}
