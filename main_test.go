package main

import (
	"testing"

	"github.com/yareeh/themoviedb-cli/internal/api"
)

func TestParseEpisodeCode(t *testing.T) {
	tests := []struct {
		code    string
		season  int
		episode int
		wantErr bool
	}{
		{"S01E02", 1, 2, false},
		{"S05E16", 5, 16, false},
		{"s01e02", 1, 2, false}, // lowercase
		{"S10E01", 10, 1, false},
		{"S00E00", 0, 0, false},
		{"S1E1", 1, 1, false},   // single digit
		{"E01S02", 0, 0, true},  // wrong order
		{"S01", 0, 0, true},     // missing episode
		{"SE01", 0, 0, true},    // missing season number
		{"01E02", 0, 0, true},   // missing S prefix
		{"", 0, 0, true},
		{"hello", 0, 0, true},
		{"SXXE01", 0, 0, true},  // non-numeric season
		{"S01EXX", 0, 0, true},  // non-numeric episode
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			season, episode, err := parseEpisodeCode(tt.code)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseEpisodeCode(%q) expected error, got season=%d episode=%d", tt.code, season, episode)
				}
				return
			}
			if err != nil {
				t.Errorf("parseEpisodeCode(%q) unexpected error: %v", tt.code, err)
				return
			}
			if season != tt.season || episode != tt.episode {
				t.Errorf("parseEpisodeCode(%q) = (%d, %d), want (%d, %d)", tt.code, season, episode, tt.season, tt.episode)
			}
		})
	}
}

func TestHasFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flag     string
		want     bool
		wantArgs []string
	}{
		{"flag present", []string{"search", "--json", "query"}, "--json", true, []string{"search", "query"}},
		{"flag absent", []string{"search", "query"}, "--json", false, []string{"search", "query"}},
		{"empty args", []string{}, "--json", false, []string{}},
		{"flag only", []string{"--json"}, "--json", true, []string{}},
		{"multiple flags", []string{"--json", "--json"}, "--json", true, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]string, len(tt.args))
			copy(args, tt.args)
			got := hasFlag(&args, tt.flag)
			if got != tt.want {
				t.Errorf("hasFlag(%v, %q) = %v, want %v", tt.args, tt.flag, got, tt.want)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("after hasFlag, args = %v, want %v", args, tt.wantArgs)
				return
			}
			for i, a := range args {
				if a != tt.wantArgs[i] {
					t.Errorf("after hasFlag, args[%d] = %q, want %q", i, a, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestExtractJWTSub(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			"valid TMDB token",
			"eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJ0ZXN0Iiwic3ViIjoiYWJjMTIzNDU2Nzg5MCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdfQ.sig",
			"abc1234567890",
		},
		{"empty string", "", ""},
		{"not a JWT", "notajwt", ""},
		{"two parts only", "part1.part2", ""},
		{"invalid base64", "header.!!!invalid!!!.sig", ""},
		{"valid base64 but not JSON", "header.dGVzdA.sig", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJWTSub(tt.token)
			if got != tt.want {
				t.Errorf("extractJWTSub() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFilterRatedMovies(t *testing.T) {
	movies := []api.RatedMovie{
		{ID: 1, Title: "Movie A", AccountRating: api.AccountRating{CreatedAt: "2026-02-15T10:00:00.000Z", Value: 8}},
		{ID: 2, Title: "Movie B", AccountRating: api.AccountRating{CreatedAt: "2026-01-05T10:00:00.000Z", Value: 7}},
		{ID: 3, Title: "Movie C", AccountRating: api.AccountRating{CreatedAt: "2025-12-20T10:00:00.000Z", Value: 9}},
		{ID: 4, Title: "Movie D", AccountRating: api.AccountRating{CreatedAt: "2025-06-01T10:00:00.000Z", Value: 6}},
		{ID: 5, Title: "Movie E", AccountRating: api.AccountRating{CreatedAt: "2024-03-15T10:00:00.000Z", Value: 5}},
	}

	t.Run("all", func(t *testing.T) {
		result := filterRatedMovies(movies, "all", "")
		if len(result) != 5 {
			t.Errorf("all: got %d, want 5", len(result))
		}
	})

	t.Run("empty mode", func(t *testing.T) {
		result := filterRatedMovies(movies, "", "")
		if len(result) != 5 {
			t.Errorf("empty mode: got %d, want 5", len(result))
		}
	})

	t.Run("last 2", func(t *testing.T) {
		result := filterRatedMovies(movies, "last", "2")
		if len(result) != 2 {
			t.Errorf("last 2: got %d, want 2", len(result))
		}
		if result[0].ID != 1 || result[1].ID != 2 {
			t.Errorf("last 2: got IDs %d,%d, want 1,2", result[0].ID, result[1].ID)
		}
	})

	t.Run("last more than total", func(t *testing.T) {
		result := filterRatedMovies(movies, "last", "100")
		if len(result) != 5 {
			t.Errorf("last 100: got %d, want 5", len(result))
		}
	})

	t.Run("from date", func(t *testing.T) {
		result := filterRatedMovies(movies, "from", "2026-01-01")
		if len(result) != 2 {
			t.Errorf("from 2026-01-01: got %d, want 2", len(result))
		}
	})

	t.Run("from date no results", func(t *testing.T) {
		result := filterRatedMovies(movies, "from", "2027-01-01")
		if len(result) != 0 {
			t.Errorf("from 2027: got %d, want 0", len(result))
		}
	})
}

func TestFilterRatedTV(t *testing.T) {
	shows := []api.RatedTV{
		{ID: 1, Name: "Show A", AccountRating: api.AccountRating{CreatedAt: "2026-02-10T10:00:00.000Z", Value: 9}},
		{ID: 2, Name: "Show B", AccountRating: api.AccountRating{CreatedAt: "2025-11-01T10:00:00.000Z", Value: 7}},
		{ID: 3, Name: "Show C", AccountRating: api.AccountRating{CreatedAt: "2025-03-15T10:00:00.000Z", Value: 8}},
	}

	t.Run("all", func(t *testing.T) {
		result := filterRatedTV(shows, "all", "")
		if len(result) != 3 {
			t.Errorf("all: got %d, want 3", len(result))
		}
	})

	t.Run("last 1", func(t *testing.T) {
		result := filterRatedTV(shows, "last", "1")
		if len(result) != 1 {
			t.Errorf("last 1: got %d, want 1", len(result))
		}
	})

	t.Run("from date", func(t *testing.T) {
		result := filterRatedTV(shows, "from", "2025-10-01")
		if len(result) != 2 {
			t.Errorf("from 2025-10-01: got %d, want 2", len(result))
		}
	})
}
