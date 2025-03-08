package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests the creation of a new PokeAPI client
func TestNewClient(t *testing.T) {
	client := NewClient(time.Hour)

	// Check that the client was created successfully
	if client.httpClient.Timeout != time.Minute {
		t.Errorf("Expected default timeout of 1 minute, got %v", client.httpClient.Timeout)
	}
}

// TestGetPokemonData tests the GetPokemonData method
func TestGetPokemonData(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is for a Pokémon
		if r.URL.Path == "/api/v2/pokemon/pikachu" {
			// Return a mock Pokémon response
			pokemon := PokemonDataResp{
				Name:   "pikachu",
				Height: 4,
				Weight: 60,
				Types: []struct {
					Slot int              `json:"slot"`
					Type NamedAPIResource `json:"type"`
				}{
					{
						Slot: 1,
						Type: NamedAPIResource{
							Name: "electric",
							URL:  "https://pokeapi.co/api/v2/type/13/",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(pokemon)
			return
		}

		// Return 404 for unknown Pokémon
		if r.URL.Path == "/api/v2/pokemon/unknown" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Default response
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create a client with a short cache duration and custom base URL
	client := NewClient(time.Millisecond * 10)

	// Store the original HTTP client to restore later
	originalClient := client.httpClient

	// Create a custom HTTP client that redirects to our test server
	client.httpClient = http.Client{
		Transport: &testTransport{
			testServer: server,
		},
	}

	// Restore the original client when done
	defer func() {
		client.httpClient = originalClient
	}()

	// Test getting a valid Pokémon
	t.Run("Valid Pokemon", func(t *testing.T) {
		pokemon, err := client.GetPokemonData("pikachu")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if pokemon.Name != "pikachu" {
			t.Errorf("Expected name 'pikachu', got '%s'", pokemon.Name)
		}

		if pokemon.Height != 4 {
			t.Errorf("Expected height 4, got %d", pokemon.Height)
		}

		if pokemon.Weight != 60 {
			t.Errorf("Expected weight 60, got %d", pokemon.Weight)
		}

		if len(pokemon.Types) != 1 {
			t.Errorf("Expected 1 type, got %d", len(pokemon.Types))
		}

		if pokemon.Types[0].Type.Name != "electric" {
			t.Errorf("Expected type 'electric', got '%s'", pokemon.Types[0].Type.Name)
		}
	})

	// Test getting an invalid Pokémon
	t.Run("Invalid Pokemon", func(t *testing.T) {
		_, err := client.GetPokemonData("unknown")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

// testTransport is a custom http.RoundTripper that redirects requests to a test server
type testTransport struct {
	testServer *httptest.Server
}

// RoundTrip implements the http.RoundTripper interface
func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Replace the request URL with the test server URL
	url := t.testServer.URL + req.URL.Path
	if req.URL.RawQuery != "" {
		url += "?" + req.URL.RawQuery
	}

	// Create a new request with the test server URL
	newReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers
	newReq.Header = req.Header

	// Send the request to the test server
	return http.DefaultTransport.RoundTrip(newReq)
}
