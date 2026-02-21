package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	req, err := newGetRequest("https://example.com/test", "my-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Header.Get("Authorization") != "Bearer my-token" {
		t.Errorf("auth header = %q, want %q", req.Header.Get("Authorization"), "Bearer my-token")
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("accept header = %q, want %q", req.Header.Get("Accept"), "application/json")
	}
	if req.Method != "GET" {
		t.Errorf("method = %q, want GET", req.Method)
	}
}

func TestDoRequestSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"page":1,"results":[]}`))
	}))
	defer ts.Close()

	req, _ := newGetRequest(ts.URL+"/test", "tok")
	data, err := doRequest(ts.Client(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"page":1,"results":[]}` {
		t.Errorf("body = %q", string(data))
	}
}

func TestDoRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"status_code":34,"status_message":"not found"}`))
	}))
	defer ts.Close()

	req, _ := newGetRequest(ts.URL+"/missing", "tok")
	_, err := doRequest(ts.Client(), req)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestSearchMoviesResponse(t *testing.T) {
	raw := `{
		"page": 1,
		"total_pages": 1,
		"total_results": 1,
		"results": [
			{"id": 603, "title": "The Matrix", "release_date": "1999-03-24", "vote_average": 8.2}
		]
	}`
	var resp SearchMoviesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Page != 1 {
		t.Errorf("Page = %d, want 1", resp.Page)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("got %d results, want 1", len(resp.Results))
	}
	m := resp.Results[0]
	if m.ID != 603 {
		t.Errorf("ID = %d, want 603", m.ID)
	}
	if m.Title != "The Matrix" {
		t.Errorf("Title = %q", m.Title)
	}
	if m.VoteAverage != 8.2 {
		t.Errorf("VoteAverage = %f, want 8.2", m.VoteAverage)
	}
}

func TestRatedMoviesResponse(t *testing.T) {
	raw := `{
		"page": 1,
		"total_pages": 2,
		"total_results": 30,
		"results": [
			{
				"id": 100,
				"title": "Test Movie",
				"release_date": "2025-01-01",
				"vote_average": 7.5,
				"account_rating": {
					"created_at": "2026-01-02T14:28:27.350Z",
					"value": 8.0
				}
			}
		]
	}`
	var resp RatedMoviesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.TotalPages != 2 {
		t.Errorf("TotalPages = %d, want 2", resp.TotalPages)
	}
	r := resp.Results[0]
	if r.AccountRating.Value != 8.0 {
		t.Errorf("rating value = %f, want 8.0", r.AccountRating.Value)
	}
	if r.AccountRating.CreatedAt != "2026-01-02T14:28:27.350Z" {
		t.Errorf("created_at = %q", r.AccountRating.CreatedAt)
	}
}

func TestRatedTVResponse(t *testing.T) {
	raw := `{
		"page": 1,
		"total_pages": 1,
		"results": [
			{
				"id": 200,
				"name": "Test Show",
				"first_air_date": "2020-05-01",
				"vote_average": 9.0,
				"account_rating": {
					"created_at": "2025-11-15T10:00:00.000Z",
					"value": 9.0
				}
			}
		]
	}`
	var resp RatedTVResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("got %d results", len(resp.Results))
	}
	if resp.Results[0].Name != "Test Show" {
		t.Errorf("Name = %q", resp.Results[0].Name)
	}
}

func TestSearchTVResponse(t *testing.T) {
	raw := `{"page":1,"results":[{"id":1396,"name":"Breaking Bad","first_air_date":"2008-01-20","vote_average":8.9}]}`
	var resp SearchTVResponse
	json.Unmarshal([]byte(raw), &resp)
	if resp.Results[0].Name != "Breaking Bad" {
		t.Errorf("Name = %q", resp.Results[0].Name)
	}
}

func TestPersonResponse(t *testing.T) {
	raw := `{"page":1,"results":[{"id":287,"name":"Brad Pitt","known_for_department":"Acting"}]}`
	var resp SearchPersonResponse
	json.Unmarshal([]byte(raw), &resp)
	if resp.Results[0].KnownForDepartment != "Acting" {
		t.Errorf("dept = %q", resp.Results[0].KnownForDepartment)
	}
}

func TestCombinedCreditsResponse(t *testing.T) {
	raw := `{"cast":[{"id":603,"title":"The Matrix","media_type":"movie","release_date":"1999-03-24","character":"Neo"},{"id":1396,"name":"Breaking Bad","media_type":"tv","first_air_date":"2008-01-20","character":"Jesse"}]}`
	var resp CombinedCreditsResponse
	json.Unmarshal([]byte(raw), &resp)
	if len(resp.Cast) != 2 {
		t.Fatalf("got %d cast", len(resp.Cast))
	}
	if resp.Cast[0].MediaType != "movie" {
		t.Errorf("cast[0].MediaType = %q", resp.Cast[0].MediaType)
	}
	if resp.Cast[1].MediaType != "tv" {
		t.Errorf("cast[1].MediaType = %q", resp.Cast[1].MediaType)
	}
}
