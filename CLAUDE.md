# themoviedb-cli

Go CLI for themoviedb.org, designed for agentic use.

## Build

```
go build -o themoviedb-cli .
```

## Project Structure

```
main.go                    # CLI entry point, all command routing
internal/
  api/
    client.go              # HTTP client (Bearer token auth)
    types.go               # Response types
    auth.go                # Auth flow (request token -> session)
    search.go              # Search movies/TV/people
    actions.go             # Rate, watchlist, favorites
    details.go             # Filmography, TV details, seasons, rated lists
  config/
    config.go              # JSON config at ~/.config/themoviedb-cli/config.json
  output/
    output.go              # Text and JSON formatters
```

## Auth Flow

1. User provides API Read Access Token (from TMDB settings)
2. CLI creates a request token, user approves in browser
3. CLI creates session_id and stores it with account_id

## API

- Base URL: https://api.themoviedb.org/3
- Auth: Bearer token header + session_id query param for writes
- Rating scale: 0.5–10.0 (TMDB scale, not 1–5)
