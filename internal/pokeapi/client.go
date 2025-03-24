package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mickk/pokedexcli/internal/pokecache"
)

var pokeCache *pokecache.Cache

func init() {
	pokeCache = pokecache.NewCache(5 * time.Minute)
}

func get(url string) ([]byte, error) {
	if cachedData, ok := pokeCache.Get(url); ok {
		return cachedData, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return []byte{}, fmt.Errorf("resource not found for url %v", url)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	pokeCache.Add(url, responseBody)

	return responseBody, nil
}

func GetLocationAreas(url string) (Resource, error) {
	responseBody, err := get(url)
	if err != nil {
		return Resource{}, err
	}

	var result Resource
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return Resource{}, err
	}

	return result, nil
}

func GetLocationArea(url string) (LocationArea, error) {
	responseBody, err := get(url)
	if err != nil {
		return LocationArea{}, err
	}

	var result LocationArea
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return LocationArea{}, err
	}

	return result, nil
}

func GetPokemonInfo(url string) (Pokemon, error) {
	responseBody, err := get(url)
	if err != nil {
		return Pokemon{}, err
	}

	var result Pokemon
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return Pokemon{}, err
	}

	return result, nil
}
