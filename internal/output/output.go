package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yareeh/themoviedb-cli/internal/api"
)

func Movies(movies []api.MovieResult, asJSON bool) {
	if asJSON {
		printJSON(movies)
		return
	}
	for i, m := range movies {
		year := yearFrom(m.ReleaseDate)
		rating := ""
		if m.Rating > 0 {
			rating = fmt.Sprintf(" [rated %.1f]", m.Rating)
		}
		fmt.Printf("%d. [%d] %s (%s) ★%.1f%s\n", i+1, m.ID, m.Title, year, m.VoteAverage, rating)
	}
}

func TVShows(shows []api.TVResult, asJSON bool) {
	if asJSON {
		printJSON(shows)
		return
	}
	for i, s := range shows {
		year := yearFrom(s.FirstAirDate)
		rating := ""
		if s.Rating > 0 {
			rating = fmt.Sprintf(" [rated %.1f]", s.Rating)
		}
		fmt.Printf("%d. [%d] %s (%s) ★%.1f%s\n", i+1, s.ID, s.Name, year, s.VoteAverage, rating)
	}
}

func People(people []api.PersonResult, asJSON bool) {
	if asJSON {
		printJSON(people)
		return
	}
	for i, p := range people {
		fmt.Printf("%d. [%d] %s (%s)\n", i+1, p.ID, p.Name, p.KnownForDepartment)
	}
}

func Filmography(credits []api.CastCredit, asJSON bool) {
	if asJSON {
		printJSON(credits)
		return
	}
	for i, c := range credits {
		title := c.Title
		date := c.ReleaseDate
		kind := "movie"
		if c.MediaType == "tv" {
			title = c.Name
			date = c.FirstAirDate
			kind = "tv"
		}
		year := yearFrom(date)
		char := ""
		if c.Character != "" {
			char = fmt.Sprintf(" as %s", c.Character)
		}
		fmt.Printf("%d. [%s:%d] %s (%s)%s\n", i+1, kind, c.ID, title, year, char)
	}
}

func Seasons(seasons []api.TVSeason, showName string, asJSON bool) {
	if asJSON {
		printJSON(seasons)
		return
	}
	fmt.Printf("%s — Seasons:\n", showName)
	for _, s := range seasons {
		fmt.Printf("  S%02d: %s (%d episodes, %s)\n", s.SeasonNumber, s.Name, s.EpisodeCount, yearFrom(s.AirDate))
	}
}

func Episodes(episodes []api.TVEpisode, seasonName string, asJSON bool) {
	if asJSON {
		printJSON(episodes)
		return
	}
	fmt.Printf("%s:\n", seasonName)
	for _, e := range episodes {
		fmt.Printf("  S%02dE%02d: %s ★%.1f (%s)\n", e.SeasonNumber, e.EpisodeNumber, e.Name, e.VoteAverage, e.AirDate)
	}
}

func RatedMovies(movies []api.RatedMovie, asJSON bool) {
	if asJSON {
		printJSON(movies)
		return
	}
	for i, m := range movies {
		year := yearFrom(m.ReleaseDate)
		ratedDate := dateOnly(m.AccountRating.CreatedAt)
		fmt.Printf("%d. [%d] %s (%s) ★%.1f [rated %.0f on %s]\n",
			i+1, m.ID, m.Title, year, m.VoteAverage, m.AccountRating.Value, ratedDate)
	}
	if len(movies) > 0 {
		fmt.Printf("\n%d movies\n", len(movies))
	}
}

func RatedTVShows(shows []api.RatedTV, asJSON bool) {
	if asJSON {
		printJSON(shows)
		return
	}
	for i, s := range shows {
		year := yearFrom(s.FirstAirDate)
		ratedDate := dateOnly(s.AccountRating.CreatedAt)
		fmt.Printf("%d. [%d] %s (%s) ★%.1f [rated %.0f on %s]\n",
			i+1, s.ID, s.Name, year, s.VoteAverage, s.AccountRating.Value, ratedDate)
	}
	if len(shows) > 0 {
		fmt.Printf("\n%d shows\n", len(shows))
	}
}

func dateOnly(ts string) string {
	if len(ts) >= 10 {
		return ts[:10]
	}
	return ts
}

func yearFrom(date string) string {
	if len(date) >= 4 {
		return date[:4]
	}
	return "????"
}

func printJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func Status(action string, err error) {
	if err != nil {
		fmt.Printf("Error: %s: %v\n", action, err)
	} else {
		fmt.Printf("OK: %s\n", action)
	}
}

func Wrap(text string, width int) string {
	words := strings.Fields(text)
	var lines []string
	line := ""
	for _, w := range words {
		if len(line)+len(w)+1 > width {
			lines = append(lines, line)
			line = w
		} else if line == "" {
			line = w
		} else {
			line += " " + w
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
