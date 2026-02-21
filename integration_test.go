package main

import (
	"os"
	"testing"
	"time"

	"github.com/yareeh/themoviedb-cli/internal/api"
)

// Test data: well-known TMDB entries
const (
	testMovieID    = 603   // The Matrix
	testMovieTitle = "The Matrix"
	testSeriesID   = 1396  // Breaking Bad
	testSeriesName = "Breaking Bad"
	testPersonID   = 287   // Brad Pitt
	testPersonName = "Brad Pitt"
	testSeasonNum  = 5
	testEpisodeSxx = "S05E16" // Felina
	testEpisodeSn  = 5
	testEpisodeEn  = 16
	testRating     = 7.0
)

func integrationClient(t *testing.T) *api.Client {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	token := os.Getenv("TMDB_ACCESS_TOKEN")
	if token == "" {
		t.Skip("TMDB_ACCESS_TOKEN not set, skipping integration test")
	}
	sessionID := os.Getenv("TMDB_SESSION_ID")
	accountID := 0
	if v := os.Getenv("TMDB_ACCOUNT_ID"); v != "" {
		for _, c := range v {
			accountID = accountID*10 + int(c-'0')
		}
	}
	accountObjectID := os.Getenv("TMDB_ACCOUNT_OBJECT_ID")
	if accountObjectID == "" {
		// Try to extract from JWT
		accountObjectID = extractJWTSub(token)
	}
	return api.New(token, sessionID, accountID, accountObjectID)
}

func TestIntegrationSearchMovie(t *testing.T) {
	client := integrationClient(t)
	resp, err := client.SearchMovies("The Matrix")
	if err != nil {
		t.Fatalf("SearchMovies: %v", err)
	}
	if len(resp.Results) == 0 {
		t.Fatal("SearchMovies returned no results")
	}
	found := false
	for _, m := range resp.Results {
		if m.ID == testMovieID {
			found = true
			if m.Title != testMovieTitle {
				t.Errorf("expected title %q, got %q", testMovieTitle, m.Title)
			}
		}
	}
	if !found {
		t.Errorf("expected to find movie ID %d in results", testMovieID)
	}
}

func TestIntegrationSearchTV(t *testing.T) {
	client := integrationClient(t)
	resp, err := client.SearchTV("Breaking Bad")
	if err != nil {
		t.Fatalf("SearchTV: %v", err)
	}
	if len(resp.Results) == 0 {
		t.Fatal("SearchTV returned no results")
	}
	found := false
	for _, s := range resp.Results {
		if s.ID == testSeriesID {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find TV ID %d in results", testSeriesID)
	}
}

func TestIntegrationSearchPerson(t *testing.T) {
	client := integrationClient(t)
	resp, err := client.SearchPerson("Brad Pitt")
	if err != nil {
		t.Fatalf("SearchPerson: %v", err)
	}
	if len(resp.Results) == 0 {
		t.Fatal("SearchPerson returned no results")
	}
	found := false
	for _, p := range resp.Results {
		if p.ID == testPersonID {
			found = true
			if p.Name != testPersonName {
				t.Errorf("expected name %q, got %q", testPersonName, p.Name)
			}
		}
	}
	if !found {
		t.Errorf("expected to find person ID %d in results", testPersonID)
	}
}

func TestIntegrationFilmography(t *testing.T) {
	client := integrationClient(t)
	resp, err := client.Filmography(testPersonID)
	if err != nil {
		t.Fatalf("Filmography: %v", err)
	}
	if len(resp.Cast) == 0 {
		t.Fatal("Filmography returned no credits")
	}
	// Brad Pitt should have a decent filmography
	if len(resp.Cast) < 10 {
		t.Errorf("expected at least 10 credits, got %d", len(resp.Cast))
	}
}

func TestIntegrationTVSeasons(t *testing.T) {
	client := integrationClient(t)
	details, err := client.TVDetails(testSeriesID)
	if err != nil {
		t.Fatalf("TVDetails: %v", err)
	}
	if details.Name != testSeriesName {
		t.Errorf("expected name %q, got %q", testSeriesName, details.Name)
	}
	// Breaking Bad has 5 seasons + specials
	if len(details.Seasons) < 5 {
		t.Errorf("expected at least 5 seasons, got %d", len(details.Seasons))
	}
}

func TestIntegrationTVEpisodes(t *testing.T) {
	client := integrationClient(t)
	details, err := client.SeasonDetails(testSeriesID, testSeasonNum)
	if err != nil {
		t.Fatalf("SeasonDetails: %v", err)
	}
	if len(details.Episodes) == 0 {
		t.Fatal("SeasonDetails returned no episodes")
	}
	// Breaking Bad S5 has 16 episodes
	if len(details.Episodes) != 16 {
		t.Errorf("expected 16 episodes, got %d", len(details.Episodes))
	}
}

func TestIntegrationRateAndUnrateMovie(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_SESSION_ID") == "" {
		t.Skip("TMDB_SESSION_ID not set, skipping write test")
	}

	// Rate the movie
	err := client.RateMovie(testMovieID, testRating)
	if err != nil {
		t.Fatalf("RateMovie: %v", err)
	}

	// TMDB needs time to propagate ratings
	time.Sleep(2 * time.Second)

	// Clean up: remove rating
	err = client.DeleteMovieRating(testMovieID)
	if err != nil {
		t.Fatalf("DeleteMovieRating: %v", err)
	}
}

func TestIntegrationRateAndUnrateTV(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_SESSION_ID") == "" {
		t.Skip("TMDB_SESSION_ID not set, skipping write test")
	}

	err := client.RateTV(testSeriesID, testRating)
	if err != nil {
		t.Fatalf("RateTV: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = client.DeleteTVRating(testSeriesID)
	if err != nil {
		t.Fatalf("DeleteTVRating: %v", err)
	}
}

func TestIntegrationRateAndUnrateEpisode(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_SESSION_ID") == "" {
		t.Skip("TMDB_SESSION_ID not set, skipping write test")
	}

	err := client.RateEpisode(testSeriesID, testEpisodeSn, testEpisodeEn, testRating)
	if err != nil {
		t.Fatalf("RateEpisode: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = client.DeleteEpisodeRating(testSeriesID, testEpisodeSn, testEpisodeEn)
	if err != nil {
		t.Fatalf("DeleteEpisodeRating: %v", err)
	}
}

func TestIntegrationWatchlistAddAndRemove(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_SESSION_ID") == "" {
		t.Skip("TMDB_SESSION_ID not set, skipping write test")
	}

	// Add to watchlist
	err := client.AddToWatchlist("movie", testMovieID)
	if err != nil {
		t.Fatalf("AddToWatchlist: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Verify it's in the watchlist
	resp, err := client.GetWatchlistMovies()
	if err != nil {
		t.Fatalf("GetWatchlistMovies: %v", err)
	}
	found := false
	for _, m := range resp.Results {
		if m.ID == testMovieID {
			found = true
		}
	}
	if !found {
		t.Error("movie not found in watchlist after adding")
	}

	// Remove from watchlist
	err = client.RemoveFromWatchlist("movie", testMovieID)
	if err != nil {
		t.Fatalf("RemoveFromWatchlist: %v", err)
	}
}

func TestIntegrationRatedMoviesList(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_ACCOUNT_OBJECT_ID") == "" {
		accountObjectID := extractJWTSub(os.Getenv("TMDB_ACCESS_TOKEN"))
		if accountObjectID == "" {
			t.Skip("Cannot determine account object ID, skipping rated list test")
		}
	}

	movies, err := client.GetAllRatedMovies()
	if err != nil {
		t.Fatalf("GetAllRatedMovies: %v", err)
	}
	// Should return some results (the test account has ratings)
	t.Logf("Found %d rated movies", len(movies))

	// Verify timestamps are present
	if len(movies) > 0 {
		if movies[0].AccountRating.CreatedAt == "" {
			t.Error("expected account_rating.created_at to be populated")
		}
		if movies[0].AccountRating.Value == 0 {
			t.Error("expected account_rating.value to be non-zero")
		}
	}
}

func TestIntegrationRatedTVList(t *testing.T) {
	client := integrationClient(t)
	if os.Getenv("TMDB_ACCOUNT_OBJECT_ID") == "" {
		accountObjectID := extractJWTSub(os.Getenv("TMDB_ACCESS_TOKEN"))
		if accountObjectID == "" {
			t.Skip("Cannot determine account object ID, skipping rated list test")
		}
	}

	shows, err := client.GetAllRatedTV()
	if err != nil {
		t.Fatalf("GetAllRatedTV: %v", err)
	}
	t.Logf("Found %d rated TV shows", len(shows))
}
