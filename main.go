package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yareeh/themoviedb-cli/internal/api"
	"github.com/yareeh/themoviedb-cli/internal/config"
	"github.com/yareeh/themoviedb-cli/internal/output"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]
	jsonFlag := hasFlag(&args, "--json")

	switch cmd {
	case "login":
		doLogin()
	case "search":
		doSearch(args, jsonFlag)
	case "filmography":
		doFilmography(args, jsonFlag)
	case "rate":
		doRate(args)
	case "watchlist":
		doWatchlist(args, jsonFlag)
	case "seasons":
		doSeasons(args, jsonFlag)
	case "episodes":
		doEpisodes(args, jsonFlag)
	case "rated":
		doRated(args, jsonFlag)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`themoviedb-cli — TMDB client for agents

Usage: themoviedb-cli <command> [options]

Commands:
  login                          Authenticate with TMDB
  search <query>                 Search movies, TV, people (prefix: movie:, tv:, person:)
  filmography <person_id>        List filmography of a person
  rate <movie|tv|episode> <id> <rating>  Rate (1-10, use S01E02 format for episodes)
  watchlist <add|remove|list> [movie|tv] [id]  Manage watchlist
  seasons <series_id>            List seasons of a TV series
  episodes <series_id> <season>  List episodes of a season
  rated [movie|tv] [all|ytd|last N|from YYYY-MM-DD]  List rated

Options:
  --json    Output as JSON instead of text

Examples:
  themoviedb-cli search "The Matrix"
  themoviedb-cli search "tv:Breaking Bad"
  themoviedb-cli search "person:Brad Pitt"
  themoviedb-cli filmography 287
  themoviedb-cli rate movie 603 9
  themoviedb-cli rate episode 1396 S05E16 10
  themoviedb-cli watchlist add movie 603
  themoviedb-cli watchlist list
  themoviedb-cli seasons 1396
  themoviedb-cli episodes 1396 5
  themoviedb-cli rated movie
  themoviedb-cli rated movie ytd
  themoviedb-cli rated movie last 10
  themoviedb-cli rated tv from 2025-06-01
`)
}

func mustClient() *api.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	if cfg.AccessToken == "" {
		fmt.Fprintln(os.Stderr, "Not logged in. Run: themoviedb-cli login")
		os.Exit(1)
	}
	return api.New(cfg.AccessToken, cfg.SessionID, cfg.AccountID)
}

func doLogin() {
	fmt.Print("Enter your TMDB API Read Access Token: ")
	var token string
	fmt.Scanln(&token)
	token = strings.TrimSpace(token)
	if token == "" {
		fmt.Fprintln(os.Stderr, "Token cannot be empty")
		os.Exit(1)
	}

	client := api.New(token, "", 0)

	// Create request token
	reqToken, err := client.CreateRequestToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	approveURL := fmt.Sprintf("https://www.themoviedb.org/authenticate/%s", reqToken)
	fmt.Printf("\nOpen this URL to approve access:\n  %s\n\nPress Enter after approving...", approveURL)
	fmt.Scanln()

	// Create session
	sessionID, err := client.CreateSession(reqToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating session: %v\n", err)
		os.Exit(1)
	}

	// Get account info
	sessionClient := api.New(token, sessionID, 0)
	account, err := sessionClient.GetAccount()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting account: %v\n", err)
		os.Exit(1)
	}

	cfg := &config.Config{
		AccessToken: token,
		SessionID:   sessionID,
		AccountID:   account.ID,
	}
	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Logged in as %s (account %d)\n", account.Username, account.ID)
}

func doSearch(args []string, jsonFlag bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli search <query>")
		os.Exit(1)
	}
	query := strings.Join(args, " ")
	client := mustClient()

	// Check for type prefix
	switch {
	case strings.HasPrefix(query, "tv:"):
		q := strings.TrimPrefix(query, "tv:")
		resp, err := client.SearchTV(strings.TrimSpace(q))
		exitOnErr(err)
		output.TVShows(resp.Results, jsonFlag)

	case strings.HasPrefix(query, "person:"):
		q := strings.TrimPrefix(query, "person:")
		resp, err := client.SearchPerson(strings.TrimSpace(q))
		exitOnErr(err)
		output.People(resp.Results, jsonFlag)

	case strings.HasPrefix(query, "movie:"):
		q := strings.TrimPrefix(query, "movie:")
		resp, err := client.SearchMovies(strings.TrimSpace(q))
		exitOnErr(err)
		output.Movies(resp.Results, jsonFlag)

	default:
		// Default: search movies
		resp, err := client.SearchMovies(query)
		exitOnErr(err)
		output.Movies(resp.Results, jsonFlag)
	}
}

func doFilmography(args []string, jsonFlag bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli filmography <person_id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	exitOnErr(err)
	client := mustClient()
	resp, err := client.Filmography(id)
	exitOnErr(err)
	output.Filmography(resp.Cast, jsonFlag)
}

func doRate(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli rate <movie|tv|episode> <id> <rating>")
		fmt.Fprintln(os.Stderr, "  For episodes: themoviedb-cli rate episode <series_id> S01E02 <rating>")
		os.Exit(1)
	}
	client := mustClient()
	mediaType := args[0]

	switch mediaType {
	case "movie":
		id, err := strconv.Atoi(args[1])
		exitOnErr(err)
		rating, err := strconv.ParseFloat(args[2], 64)
		exitOnErr(err)
		err = client.RateMovie(id, rating)
		output.Status(fmt.Sprintf("rated movie %d as %.1f", id, rating), err)

	case "tv":
		id, err := strconv.Atoi(args[1])
		exitOnErr(err)
		rating, err := strconv.ParseFloat(args[2], 64)
		exitOnErr(err)
		err = client.RateTV(id, rating)
		output.Status(fmt.Sprintf("rated TV %d as %.1f", id, rating), err)

	case "episode":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli rate episode <series_id> S01E02 <rating>")
			os.Exit(1)
		}
		seriesID, err := strconv.Atoi(args[1])
		exitOnErr(err)
		season, episode, err := parseEpisodeCode(args[2])
		exitOnErr(err)
		rating, err := strconv.ParseFloat(args[3], 64)
		exitOnErr(err)
		err = client.RateEpisode(seriesID, season, episode, rating)
		output.Status(fmt.Sprintf("rated S%02dE%02d of %d as %.1f", season, episode, seriesID, rating), err)

	default:
		fmt.Fprintf(os.Stderr, "Unknown media type: %s (use movie, tv, or episode)\n", mediaType)
		os.Exit(1)
	}
}

func doWatchlist(args []string, jsonFlag bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli watchlist <add|remove|list> [movie|tv] [id]")
		os.Exit(1)
	}
	client := mustClient()
	action := args[0]

	switch action {
	case "list":
		mediaType := "movie"
		if len(args) > 1 {
			mediaType = args[1]
		}
		if mediaType == "tv" {
			resp, err := client.GetWatchlistTV()
			exitOnErr(err)
			output.TVShows(resp.Results, jsonFlag)
		} else {
			resp, err := client.GetWatchlistMovies()
			exitOnErr(err)
			output.Movies(resp.Results, jsonFlag)
		}

	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli watchlist add <movie|tv> <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[2])
		exitOnErr(err)
		err = client.AddToWatchlist(args[1], id)
		output.Status(fmt.Sprintf("added %s %d to watchlist", args[1], id), err)

	case "remove":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli watchlist remove <movie|tv> <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[2])
		exitOnErr(err)
		err = client.RemoveFromWatchlist(args[1], id)
		output.Status(fmt.Sprintf("removed %s %d from watchlist", args[1], id), err)

	default:
		fmt.Fprintf(os.Stderr, "Unknown watchlist action: %s\n", action)
		os.Exit(1)
	}
}

func doSeasons(args []string, jsonFlag bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli seasons <series_id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	exitOnErr(err)
	client := mustClient()
	details, err := client.TVDetails(id)
	exitOnErr(err)
	output.Seasons(details.Seasons, details.Name, jsonFlag)
}

func doEpisodes(args []string, jsonFlag bool) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: themoviedb-cli episodes <series_id> <season_number>")
		os.Exit(1)
	}
	seriesID, err := strconv.Atoi(args[0])
	exitOnErr(err)
	seasonNum, err := strconv.Atoi(args[1])
	exitOnErr(err)
	client := mustClient()
	details, err := client.SeasonDetails(seriesID, seasonNum)
	exitOnErr(err)
	output.Episodes(details.Episodes, details.Name, jsonFlag)
}

func doRated(args []string, jsonFlag bool) {
	// Parse: rated [movie|tv] [all|ytd|last N|from YYYY-MM-DD]
	mediaType := "movie"
	filterMode := "all"
	filterValue := ""

	i := 0
	if i < len(args) && (args[i] == "movie" || args[i] == "tv") {
		mediaType = args[i]
		i++
	}
	if i < len(args) {
		filterMode = args[i]
		i++
	}
	if i < len(args) {
		filterValue = args[i]
	}

	client := mustClient()

	if mediaType == "tv" {
		shows, err := client.GetAllRatedTV()
		exitOnErr(err)
		shows = filterRatedTV(shows, filterMode, filterValue)
		output.RatedTVShows(shows, jsonFlag)
	} else {
		movies, err := client.GetAllRatedMovies()
		exitOnErr(err)
		movies = filterRatedMovies(movies, filterMode, filterValue)
		output.RatedMovies(movies, jsonFlag)
	}
}

func filterRatedMovies(movies []api.RatedMovie, mode, value string) []api.RatedMovie {
	switch mode {
	case "all", "":
		return movies
	case "ytd":
		cutoff := fmt.Sprintf("%d-01-01", currentYear())
		var filtered []api.RatedMovie
		for _, m := range movies {
			if m.AccountRating.CreatedAt >= cutoff {
				filtered = append(filtered, m)
			}
		}
		return filtered
	case "last":
		n, err := strconv.Atoi(value)
		if err != nil || n <= 0 {
			fmt.Fprintln(os.Stderr, "Usage: rated [movie|tv] last <N>")
			os.Exit(1)
		}
		if n > len(movies) {
			n = len(movies)
		}
		return movies[:n]
	case "from":
		if value == "" {
			fmt.Fprintln(os.Stderr, "Usage: rated [movie|tv] from YYYY-MM-DD")
			os.Exit(1)
		}
		var filtered []api.RatedMovie
		for _, m := range movies {
			if m.AccountRating.CreatedAt >= value {
				filtered = append(filtered, m)
			}
		}
		return filtered
	default:
		fmt.Fprintf(os.Stderr, "Unknown filter: %s (use all, ytd, last N, or from YYYY-MM-DD)\n", mode)
		os.Exit(1)
		return nil
	}
}

func filterRatedTV(shows []api.RatedTV, mode, value string) []api.RatedTV {
	switch mode {
	case "all", "":
		return shows
	case "ytd":
		cutoff := fmt.Sprintf("%d-01-01", currentYear())
		var filtered []api.RatedTV
		for _, s := range shows {
			if s.AccountRating.CreatedAt >= cutoff {
				filtered = append(filtered, s)
			}
		}
		return filtered
	case "last":
		n, err := strconv.Atoi(value)
		if err != nil || n <= 0 {
			fmt.Fprintln(os.Stderr, "Usage: rated [movie|tv] last <N>")
			os.Exit(1)
		}
		if n > len(shows) {
			n = len(shows)
		}
		return shows[:n]
	case "from":
		if value == "" {
			fmt.Fprintln(os.Stderr, "Usage: rated [movie|tv] from YYYY-MM-DD")
			os.Exit(1)
		}
		var filtered []api.RatedTV
		for _, s := range shows {
			if s.AccountRating.CreatedAt >= value {
				filtered = append(filtered, s)
			}
		}
		return filtered
	default:
		fmt.Fprintf(os.Stderr, "Unknown filter: %s (use all, ytd, last N, or from YYYY-MM-DD)\n", mode)
		os.Exit(1)
		return nil
	}
}

func currentYear() int {
	return time.Now().Year()
}

// parseEpisodeCode parses "S01E02" into season=1, episode=2
func parseEpisodeCode(code string) (int, int, error) {
	code = strings.ToUpper(code)
	if !strings.HasPrefix(code, "S") || !strings.Contains(code, "E") {
		return 0, 0, fmt.Errorf("invalid episode code %q (use S01E02 format)", code)
	}
	parts := strings.SplitN(code[1:], "E", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid episode code %q", code)
	}
	season, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid season in %q: %w", code, err)
	}
	episode, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid episode in %q: %w", code, err)
	}
	return season, episode, nil
}

func hasFlag(args *[]string, flag string) bool {
	filtered := make([]string, 0, len(*args))
	found := false
	for _, a := range *args {
		if a == flag {
			found = true
		} else {
			filtered = append(filtered, a)
		}
	}
	*args = filtered
	return found
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
