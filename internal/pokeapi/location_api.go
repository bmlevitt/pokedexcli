package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListLocationAreas retrieves a list of location areas from the PokeAPI
// If pageURL is provided, it fetches that specific page, otherwise fetches the first page
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResp, error) {
	endpoint := "/location-area?offset=0&limit=20"
	fullURL := baseURL + endpoint
	if pageURL != nil {
		fullURL = *pageURL
	}

	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		locationAreasResp := LocationAreasResp{}
		err := json.Unmarshal(data, &locationAreasResp)
		if err != nil {
			return LocationAreasResp{}, err
		}
		return locationAreasResp, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreasResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return LocationAreasResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResp{}, err
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	locationAreasResp := LocationAreasResp{}
	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return LocationAreasResp{}, err
	}

	return locationAreasResp, nil
}

// ExploreLocation retrieves information about Pokemon that can be found in a specific location area
func (c *Client) ExploreLocation(location string) (LocationExploreResp, error) {
	endpoint := "/location-area/"
	fullURL := baseURL + endpoint + location

	fmt.Printf("exploring %s...\n", location)

	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		locationExploreResp := LocationExploreResp{}
		err := json.Unmarshal(data, &locationExploreResp)
		if err != nil {
			return LocationExploreResp{}, err
		}
		return locationExploreResp, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationExploreResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationExploreResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return LocationExploreResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationExploreResp{}, err
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	locationExploreResp := LocationExploreResp{}
	err = json.Unmarshal(body, &locationExploreResp)
	if err != nil {
		return LocationExploreResp{}, err
	}

	return locationExploreResp, nil
}
