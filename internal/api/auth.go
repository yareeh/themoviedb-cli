package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateRequestToken() (string, error) {
	data, err := c.get("/authentication/token/new", nil)
	if err != nil {
		return "", fmt.Errorf("creating request token: %w", err)
	}
	var resp RequestTokenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("failed to create request token")
	}
	return resp.RequestToken, nil
}

func (c *Client) CreateSession(requestToken string) (string, error) {
	payload := map[string]string{"request_token": requestToken}
	data, err := c.post("/authentication/session/new", payload)
	if err != nil {
		return "", fmt.Errorf("creating session: %w", err)
	}
	var resp SessionResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("failed to create session")
	}
	return resp.SessionID, nil
}

func (c *Client) GetAccount() (*AccountResponse, error) {
	data, err := c.get("/account", nil)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}
	var resp AccountResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
