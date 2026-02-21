package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) SearchMovies(query string) (*SearchMoviesResponse, error) {
	params := url.Values{"query": {query}}
	data, err := c.get("/search/movie", params)
	if err != nil {
		return nil, fmt.Errorf("searching movies: %w", err)
	}
	var resp SearchMoviesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SearchTV(query string) (*SearchTVResponse, error) {
	params := url.Values{"query": {query}}
	data, err := c.get("/search/tv", params)
	if err != nil {
		return nil, fmt.Errorf("searching TV: %w", err)
	}
	var resp SearchTVResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SearchPerson(query string) (*SearchPersonResponse, error) {
	params := url.Values{"query": {query}}
	data, err := c.get("/search/person", params)
	if err != nil {
		return nil, fmt.Errorf("searching people: %w", err)
	}
	var resp SearchPersonResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
