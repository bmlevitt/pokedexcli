package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ExploreLocation(location string) (LocationExploreResp, error) {
	endpoint := "/location-area/"
	fullURL := baseURL + endpoint + location

	fmt.Printf("exploring %s...\n", location)

	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		locationExploreResp := LocationAreasResp{}
		err := json.Unmarshal(data, &locationExploreResp)
		if err != nil {
			return LocationExploreResp{}, err
		}
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationExploreResp{}, err
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationExploreResp{}, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode > 399 {
		return LocationExploreResp{}, fmt.Errorf("bad status received: %v", resp.StatusCode)
	}

	// Read the response body
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return LocationExploreResp{}, nil
	}

	// Parse the JSON response
	locationExploreResp := LocationExploreResp{}
	err = json.Unmarshal(data, &locationExploreResp)
	if err != nil {
		return LocationExploreResp{}, err
	}

	c.cache.Add(fullURL, data)
	return locationExploreResp, nil
}
