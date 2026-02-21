package output

import (
	"testing"
)

func TestYearFrom(t *testing.T) {
	tests := []struct {
		date string
		want string
	}{
		{"2024-11-01", "2024"},
		{"1999-03-24", "1999"},
		{"2025", "2025"},
		{"", "????"},
		{"abc", "????"},
		{"99", "????"},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			got := yearFrom(tt.date)
			if got != tt.want {
				t.Errorf("yearFrom(%q) = %q, want %q", tt.date, got, tt.want)
			}
		})
	}
}

func TestDateOnly(t *testing.T) {
	tests := []struct {
		ts   string
		want string
	}{
		{"2026-01-02T14:28:27.350Z", "2026-01-02"},
		{"2025-12-20T00:00:00.000Z", "2025-12-20"},
		{"2024-03-15", "2024-03-15"},
		{"short", "short"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ts, func(t *testing.T) {
			got := dateOnly(tt.ts)
			if got != tt.want {
				t.Errorf("dateOnly(%q) = %q, want %q", tt.ts, got, tt.want)
			}
		})
	}
}

func TestTmdbURL(t *testing.T) {
	tests := []struct {
		mediaType string
		id        int
		want      string
	}{
		{"movie", 603, "https://www.themoviedb.org/movie/603"},
		{"tv", 1396, "https://www.themoviedb.org/tv/1396"},
		{"person", 287, "https://www.themoviedb.org/person/287"},
	}

	for _, tt := range tests {
		t.Run(tt.mediaType, func(t *testing.T) {
			got := tmdbURL(tt.mediaType, tt.id)
			if got != tt.want {
				t.Errorf("tmdbURL(%q, %d) = %q, want %q", tt.mediaType, tt.id, got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{"short text", "hello world", 80, "hello world"},
		{"wrap at width", "hello world foo bar", 11, "hello world\nfoo bar"},
		{"empty", "", 80, ""},
		{"single long word", "superlongword", 80, "superlongword"},
		{"exact width", "hello world", 11, "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.text, tt.width)
			if got != tt.want {
				t.Errorf("Wrap(%q, %d) = %q, want %q", tt.text, tt.width, got, tt.want)
			}
		})
	}
}
