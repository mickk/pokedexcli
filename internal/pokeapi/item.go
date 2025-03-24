package pokeapi

type ResourceResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Resource struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []ResourceResults `json:"results"`
}
