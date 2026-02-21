package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) Filmography(personID int) (*CombinedCreditsResponse, error) {
	path := fmt.Sprintf("/person/%d/combined_credits", personID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting filmography: %w", err)
	}
	var resp CombinedCreditsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TVDetails(seriesID int) (*TVDetails, error) {
	path := fmt.Sprintf("/tv/%d", seriesID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting TV details: %w", err)
	}
	var resp TVDetails
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SeasonDetails(seriesID, seasonNumber int) (*SeasonDetails, error) {
	path := fmt.Sprintf("/tv/%d/season/%d", seriesID, seasonNumber)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting season details: %w", err)
	}
	var resp SeasonDetails
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetRatedMovies() (*SearchMoviesResponse, error) {
	path := fmt.Sprintf("/account/%d/rated/movies", c.accountID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting rated movies: %w", err)
	}
	var resp SearchMoviesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetRatedTV() (*SearchTVResponse, error) {
	path := fmt.Sprintf("/account/%d/rated/tv", c.accountID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting rated TV: %w", err)
	}
	var resp SearchTVResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetWatchlistMovies() (*SearchMoviesResponse, error) {
	path := fmt.Sprintf("/account/%d/watchlist/movies", c.accountID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting movie watchlist: %w", err)
	}
	var resp SearchMoviesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetWatchlistTV() (*SearchTVResponse, error) {
	path := fmt.Sprintf("/account/%d/watchlist/tv", c.accountID)
	data, err := c.get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting TV watchlist: %w", err)
	}
	var resp SearchTVResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
