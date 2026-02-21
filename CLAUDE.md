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

- Base URL: https://api.themoviedb.org/3 (v4 for rated items with timestamps)
- Auth: Bearer token header + session_id query param for writes
- Rating scale: 0.5–10.0 (TMDB scale, not 1–5)
- V4 rated endpoints use account_object_id (hex, from JWT sub claim), not numeric account_id

## Releasing

Tag and release using `gh`:

```bash
# Tag the release
git tag v0.1.0
git push origin v0.1.0

# Create GitHub release
gh release create v0.1.0 --generate-notes

# Or with a title and notes
gh release create v0.1.0 --title "v0.1.0" --notes "First release"
```

To attach binaries, cross-compile first:

```bash
GOOS=linux GOARCH=amd64 go build -o themoviedb-cli-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o themoviedb-cli-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o themoviedb-cli-windows-amd64.exe .
gh release create v0.1.0 --generate-notes \
  themoviedb-cli-linux-amd64 \
  themoviedb-cli-darwin-arm64 \
  themoviedb-cli-windows-amd64.exe
```

## Codecov

Coverage is uploaded automatically by CI. To activate:
1. Sign in to [codecov.io](https://codecov.io) with GitHub
2. Add the repo
3. Add `CODECOV_TOKEN` to repo secrets (Settings > Secrets)
