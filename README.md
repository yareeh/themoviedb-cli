# themoviedb-cli

[![CI](https://github.com/yareeh/themoviedb-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/yareeh/themoviedb-cli/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/yareeh/themoviedb-cli/graph/badge.svg)](https://codecov.io/gh/yareeh/themoviedb-cli)
[![Go](https://img.shields.io/github/go-mod/go-version/yareeh/themoviedb-cli)](https://go.dev)
[![Release](https://img.shields.io/github/v/release/yareeh/themoviedb-cli)](https://github.com/yareeh/themoviedb-cli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yareeh/themoviedb-cli)](https://goreportcard.com/report/github.com/yareeh/themoviedb-cli)
[![Built for AI Agents](https://img.shields.io/badge/built%20for-AI%20agents-blueviolet)](https://github.com/yareeh/themoviedb-cli)

A concise CLI for [The Movie Database (TMDB)](https://www.themoviedb.org), designed for agentic use with AI assistants.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install yareeh/tap/themoviedb-cli
```

### Go install

```bash
go install github.com/yareeh/themoviedb-cli@latest
```

### Binary download

Download pre-built binaries from the [Releases](https://github.com/yareeh/themoviedb-cli/releases) page.

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
