package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const baseURLv4 = "https://api.themoviedb.org/4"

type AccountRating struct {
	CreatedAt string  `json:"created_at"`
	Value     float64 `json:"value"`
}

type RatedMovie struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	ReleaseDate    string        `json:"release_date"`
	Overview       string        `json:"overview"`
	VoteAverage    float64       `json:"vote_average"`
	AccountRating  AccountRating `json:"account_rating"`
}

type RatedTV struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	FirstAirDate   string        `json:"first_air_date"`
	Overview       string        `json:"overview"`
	VoteAverage    float64       `json:"vote_average"`
	AccountRating  AccountRating `json:"account_rating"`
}

type RatedMoviesResponse struct {
	Page         int          `json:"page"`
	Results      []RatedMovie `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

type RatedTVResponse struct {
	Page         int       `json:"page"`
	Results      []RatedTV `json:"results"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
}

// getV4 makes a GET request to the V4 API.
func (c *Client) getV4(path string, params url.Values) (json.RawMessage, error) {
	u := baseURLv4 + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := newGetRequest(u, c.token)
	if err != nil {
		return nil, err
	}
	return doRequest(c.http, req)
}

// GetAllRatedMovies fetches all pages of rated movies from V4 API (includes rating timestamps).
func (c *Client) GetAllRatedMovies() ([]RatedMovie, error) {
	var all []RatedMovie
	page := 1
	for {
		params := url.Values{
			"page":    {strconv.Itoa(page)},
			"sort_by": {"created_at.desc"},
		}
		path := fmt.Sprintf("/account/%s/movie/rated", c.accountObjectID)
		data, err := c.getV4(path, params)
		if err != nil {
			return nil, fmt.Errorf("getting rated movies: %w", err)
		}
		var resp RatedMoviesResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Results...)
		if page >= resp.TotalPages {
			break
		}
		page++
	}
	return all, nil
}

// GetAllRatedTV fetches all pages of rated TV from V4 API.
func (c *Client) GetAllRatedTV() ([]RatedTV, error) {
	var all []RatedTV
	page := 1
	for {
		params := url.Values{
			"page":    {strconv.Itoa(page)},
			"sort_by": {"created_at.desc"},
		}
		path := fmt.Sprintf("/account/%s/tv/rated", c.accountObjectID)
		data, err := c.getV4(path, params)
		if err != nil {
			return nil, fmt.Errorf("getting rated TV: %w", err)
		}
		var resp RatedTVResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Results...)
		if page >= resp.TotalPages {
			break
		}
		page++
	}
	return all, nil
}

// GetRatedMoviesPage fetches a single page of rated movies (sorted by newest first).
func (c *Client) GetRatedMoviesPage(page, pageSize int) (*RatedMoviesResponse, error) {
	params := url.Values{
		"page":    {strconv.Itoa(page)},
		"sort_by": {"created_at.desc"},
	}
	path := fmt.Sprintf("/account/%s/movie/rated", c.accountObjectID)
	data, err := c.getV4(path, params)
	if err != nil {
		return nil, fmt.Errorf("getting rated movies: %w", err)
	}
	var resp RatedMoviesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
