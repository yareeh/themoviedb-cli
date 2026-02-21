package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) RateMovie(movieID int, rating float64) error {
	payload := map[string]float64{"value": rating}
	data, err := c.post(fmt.Sprintf("/movie/%d/rating", movieID), payload)
	if err != nil {
		return fmt.Errorf("rating movie: %w", err)
	}
	var resp StatusResponse
	json.Unmarshal(data, &resp)
	return nil
}

func (c *Client) RateTV(seriesID int, rating float64) error {
	payload := map[string]float64{"value": rating}
	data, err := c.post(fmt.Sprintf("/tv/%d/rating", seriesID), payload)
	if err != nil {
		return fmt.Errorf("rating TV: %w", err)
	}
	var resp StatusResponse
	json.Unmarshal(data, &resp)
	return nil
}

func (c *Client) RateEpisode(seriesID, season, episode int, rating float64) error {
	payload := map[string]float64{"value": rating}
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/rating", seriesID, season, episode)
	data, err := c.post(path, payload)
	if err != nil {
		return fmt.Errorf("rating episode: %w", err)
	}
	var resp StatusResponse
	json.Unmarshal(data, &resp)
	return nil
}

func (c *Client) AddToWatchlist(mediaType string, mediaID int) error {
	payload := map[string]any{
		"media_type": mediaType,
		"media_id":   mediaID,
		"watchlist":  true,
	}
	_, err := c.post(fmt.Sprintf("/account/%d/watchlist", c.accountID), payload)
	if err != nil {
		return fmt.Errorf("adding to watchlist: %w", err)
	}
	return nil
}

func (c *Client) RemoveFromWatchlist(mediaType string, mediaID int) error {
	payload := map[string]any{
		"media_type": mediaType,
		"media_id":   mediaID,
		"watchlist":  false,
	}
	_, err := c.post(fmt.Sprintf("/account/%d/watchlist", c.accountID), payload)
	if err != nil {
		return fmt.Errorf("removing from watchlist: %w", err)
	}
	return nil
}

func (c *Client) AddFavorite(mediaType string, mediaID int) error {
	payload := map[string]any{
		"media_type": mediaType,
		"media_id":   mediaID,
		"favorite":   true,
	}
	_, err := c.post(fmt.Sprintf("/account/%d/favorite", c.accountID), payload)
	if err != nil {
		return fmt.Errorf("adding favorite: %w", err)
	}
	return nil
}
