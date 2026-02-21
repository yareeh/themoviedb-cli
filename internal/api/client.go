package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://api.themoviedb.org/3"

type Client struct {
	token           string
	sessionID       string
	accountID       int
	accountObjectID string
	http            *http.Client
}

func New(token, sessionID string, accountID int, accountObjectID string) *Client {
	return &Client{
		token:           token,
		sessionID:       sessionID,
		accountID:       accountID,
		accountObjectID: accountObjectID,
		http:            &http.Client{},
	}
}

func newGetRequest(u, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func doRequest(client *http.Client, req *http.Request) (json.RawMessage, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) get(path string, params url.Values) (json.RawMessage, error) {
	u := baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := newGetRequest(u, c.token)
	if err != nil {
		return nil, err
	}
	return doRequest(c.http, req)
}

func (c *Client) post(path string, payload any) (json.RawMessage, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	u := baseURL + path
	if c.sessionID != "" {
		u += "?session_id=" + c.sessionID
	}
	req, err := http.NewRequest("POST", u, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) delete(path string) (json.RawMessage, error) {
	u := baseURL + path
	if c.sessionID != "" {
		u += "?session_id=" + c.sessionID
	}
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
