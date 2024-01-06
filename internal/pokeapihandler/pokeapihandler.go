package pokeapihandler

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"time"
	"errors"
	"internal/pokecache"
)

const pokeBaseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache pokecache.Cache
}

func NewHTTPClient(timeout time.Duration, cache pokecache.Cache) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: cache,
	}
}

func (c *Client) ListLocations(limit *int, offset *int, url *string) (PaginatedResponse, error) {

	var finalURL string

	if url != nil {
		finalURL = *url
	} else {
		limitParam := ""
		offsetParam := ""

		if limit != nil {
			limitParam += fmt.Sprintf("limit=%d", *limit)
		}

		if offset != nil {
			offsetParam += fmt.Sprintf("offset=%d", *offset)
		}

		params := fmt.Sprintf("?%s&%s", limitParam, offsetParam)

		finalURL = pokeBaseURL + "/location-area" + params
	}

	// Read-through cache strategy
	body, ok := c.cache.Get(finalURL)

	if !ok {
		//	return res, nil
		res, err := http.Get(finalURL)

		if err != nil {
			return PaginatedResponse{}, err
		}

		// Read body
		body, err = io.ReadAll(res.Body)
		// Is the developer's responsability to close res.Body
		res.Body.Close()

		if res.StatusCode > 299 {
			return PaginatedResponse{}, errors.New(fmt.Sprintf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body))
		}

		if err != nil {
			return PaginatedResponse{}, err
		}

		c.cache.Add(finalURL, []byte(body))
	}
	// Unmarshal the JSON body
	response := PaginatedResponse{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return PaginatedResponse{}, err
	}

	return response, nil
}

func (c *Client) GetLocation(name string) (LocationArea, error){
	if name == "" {
		return LocationArea{}, errors.New("Provide a valid location area name.")
	}

	finalURL := pokeBaseURL + "/location-area/" + name

	// Read-through cache strategy
	body, ok := c.cache.Get(finalURL)

	if !ok {
		// Data not in cache, then we make the request
		res, err := http.Get(finalURL)

		if err != nil {
			return LocationArea{}, err
		}

		// Read body
		body, err = io.ReadAll(res.Body)
		// Is the developer's responsability to close res.Body
		res.Body.Close()

		if res.StatusCode > 299 {
			return LocationArea{}, errors.New(fmt.Sprintf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body))
		}

		if err != nil {
			return LocationArea{}, err
		}

		c.cache.Add(finalURL, []byte(body))
	}
	// Unmarshal the JSON body
	response := LocationArea{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return LocationArea{}, err
	}

	return response, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	if name == "" {
		return Pokemon{}, errors.New("Provide a valid Pokemon name (e.g. pikachu, mewtwo, gengar, etc)")
	}

	finalURL := pokeBaseURL + "/pokemon/" + name

	// Read-through cache strategy
	body, ok := c.cache.Get(finalURL)

	if !ok {
		// Data not in cache, then we make the request
		res, err := http.Get(finalURL)

		if err != nil {
			return Pokemon{}, err
		}

		// Read body
		body, err = io.ReadAll(res.Body)

		// Is the developer's responsability to close res.Body
		res.Body.Close()

		if res.StatusCode > 299 {
			return Pokemon{}, errors.New(fmt.Sprintf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body))
		}

		if err != nil {
			return Pokemon{}, err
		}

		c.cache.Add(finalURL, []byte(body))
	}

	// Unmarshal the JSON body
	response := Pokemon{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return Pokemon{}, err
	}

	return response, nil
}
