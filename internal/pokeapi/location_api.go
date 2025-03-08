package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// ListLocationAreas retrieves a paginated list of location areas from the PokeAPI.
// Location areas are specific places within the Pokémon world where Pokémon can be encountered.
// This function supports pagination, allowing navigation through all available locations.
// Results are cached to improve performance and reduce API calls.
//
// Parameters:
//   - pageURL: Optional URL for a specific page of results. If nil, retrieves the first page.
//
// Returns:
//   - A LocationAreasResp containing the list of location areas and pagination URLs
//   - An error if the API request fails
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResp, error) {
	endpoint := "/location-area?offset=0&limit=20"
	fullURL := baseURL + endpoint
	if pageURL != nil {
		fullURL = *pageURL
	}

	data, ok := c.cache.Get(fullURL)
	if ok {
		locationAreasResp := LocationAreasResp{}
		err := json.Unmarshal(data, &locationAreasResp)
		if err != nil {
			return LocationAreasResp{}, fmt.Errorf("error unmarshaling cached location data: %w", err)
		}
		return locationAreasResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreasResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		reqEndpoint := endpoint
		if pageURL != nil {
			reqEndpoint = *pageURL
		}
		return LocationAreasResp{}, errorhandling.NewAPIError(resp.StatusCode, reqEndpoint, fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	locationAreasResp := LocationAreasResp{}
	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return LocationAreasResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return locationAreasResp, nil
}

// ExploreLocation retrieves a list of Pokémon that can be encountered at a specific location area.
// This function is used by the "explore" command to show which Pokémon are available for catching
// at a given location.
// Results are cached to improve performance and reduce API calls.
//
// Parameters:
//   - location: The name or ID of the location area to explore (in lowercase with hyphens)
//
// Returns:
//   - A LocationExploreResp containing the list of Pokémon encounters at the location
//   - An error if the API request fails or the location doesn't exist
func (c *Client) ExploreLocation(location string) (LocationExploreResp, error) {
	endpoint := "/location-area/"
	fullURL := baseURL + endpoint + location

	data, ok := c.cache.Get(fullURL)
	if ok {
		locationExploreResp := LocationExploreResp{}
		err := json.Unmarshal(data, &locationExploreResp)
		if err != nil {
			return LocationExploreResp{}, fmt.Errorf("error unmarshaling cached location data: %w", err)
		}
		return locationExploreResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationExploreResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationExploreResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusNotFound {
			return LocationExploreResp{}, errorhandling.LocationNotFoundError(location, fmt.Errorf("HTTP 404"))
		}
		return LocationExploreResp{}, errorhandling.NewAPIError(resp.StatusCode, endpoint+location, fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationExploreResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	locationExploreResp := LocationExploreResp{}
	err = json.Unmarshal(body, &locationExploreResp)
	if err != nil {
		return LocationExploreResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return locationExploreResp, nil
}
