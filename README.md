# themoviedb-cli

A concise CLI for [The Movie Database (TMDB)](https://www.themoviedb.org), designed for agentic use with AI assistants.

## Installation

```bash
go build -o themoviedb-cli .
# Or install globally:
go install github.com/yareeh/themoviedb-cli@latest
```

## Setup

1. Get a TMDB API Read Access Token from [TMDB Settings](https://www.themoviedb.org/settings/api)
2. Run `themoviedb-cli login` and paste your token
3. Approve access in the browser when prompted

Config is stored at `~/.config/themoviedb-cli/config.json`.

## Usage

### Search

```bash
# Search movies (default)
themoviedb-cli search "The Matrix"

# Search TV shows
themoviedb-cli search "tv:Breaking Bad"

# Search people
themoviedb-cli search "person:Brad Pitt"

# JSON output
themoviedb-cli search "The Matrix" --json
```

### Filmography

```bash
themoviedb-cli filmography 287    # Brad Pitt's filmography
```

### Rate

```bash
# Rate a movie (scale: 1-10)
themoviedb-cli rate movie 603 9

# Rate a TV series
themoviedb-cli rate tv 1396 10

# Rate a specific episode
themoviedb-cli rate episode 1396 S05E16 10
```

### Watchlist

```bash
themoviedb-cli watchlist add movie 603
themoviedb-cli watchlist add tv 1396
themoviedb-cli watchlist list
themoviedb-cli watchlist list tv
themoviedb-cli watchlist remove movie 603
```

### TV Seasons & Episodes

```bash
# List seasons
themoviedb-cli seasons 1396

# List episodes in a season
themoviedb-cli episodes 1396 5
```

### Your Ratings

```bash
themoviedb-cli rated movie
themoviedb-cli rated tv
```

## Output Formats

All commands support `--json` for machine-readable JSON output. Default is human-readable text with numbered results.

Text format: `1. [id] Title (year) ★rating`

## License

MIT
