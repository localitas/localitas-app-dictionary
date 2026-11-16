---
title: Dictionary
description: Word lookup and definitions
---

# Dictionary

Look up word definitions from multiple dictionary sources simultaneously.

## Word Lookup

Search for a word and receive definitions from all configured sources in parallel. The app queries multiple dictionaries and returns combined results.

**GET /api/lookup?word={word}** - Look up a word
**GET /api/lookup/{word}** - Look up a word (path parameter variant)

## Sources

The dictionary app supports two sources:

- **dictionary** - Free Dictionary API (dictionaryapi.dev), provides formal definitions, phonetics, parts of speech, and example sentences
- **urban** - Urban Dictionary, provides crowd-sourced definitions with upvotes and downvotes

Both sources are enabled by default. Configure which sources to use at app initialization by passing a comma-separated list.

## Response Format

Each lookup returns a list of dictionary entries. Each entry contains:
- The source name (dictionary or urban)
- Word and phonetic pronunciation
- Part of speech
- Definitions with examples
- Audio pronunciation URL (when available from the dictionary source)

## Web Interface

The app provides a browser-based UI for quick lookups. Type a word and see definitions from all sources rendered inline.

## Build & Deploy

Run these in the app's repository, `~/localitas-app-dictionary`.

### Version

```bash
bin/dictionary-server --version
```

### Build from source

```bash
# Native binary
go build -o bin/dictionary-server ./cmd/dictionary-server

# Native binary plus the Linux amd64 binary the Docker image is built from (dictionary-server-linux-amd64)
make build-linux
```

`make lint` runs gofmt and go vet; `make test` runs them and then the tests. `make build` runs only the lint step.

### Run natively

```bash
make start   # writes the launchd service com.localitas.app.dictionary for bin/dictionary-server and loads it on port 9205
make stop
make logs    # follows ~/.localitas/logs/dictionary/stdout.log
```

The app announces itself to core when it starts, and core serves it at `/apps/ext/dictionary/`.

### Docker

```bash
make build-docker   # builds the Linux binary and the dictionary:latest image
make start-docker   # runs container "dictionary" on port 9205, labelled localitas.app=true
make stop-docker
make logs-docker
```

Core registers containers labelled `localitas.app=true` once Docker reports them healthy.

### Publish

```bash
# Run the tests, build the image, push ghcr.io/localitas/localitas-app-dictionary:latest and :<VERSION>
make docker-push

# Build dist/dictionary-server-darwin-arm64 and dist/dictionary-server-darwin-amd64
# and create GitHub release vX.Y.Z on this repository with them
make release VERSION=vX.Y.Z
```

`VERSION` defaults to `git describe --tags --always --dirty`.

### Download

Pre-built macOS binaries are on the [releases page](https://github.com/localitas/localitas-app-dictionary/releases):
- `dictionary-server-darwin-arm64` (Apple Silicon)
- `dictionary-server-darwin-amd64` (Intel)

Every push to `main` publishes a release: `.github/workflows/release.yml` runs the tests and releases both binaries as `v0.1.<commit count>.<UTC timestamp>`. `make release` publishes one by hand.

Download with the GitHub CLI:

    gh release download --repo localitas/localitas-app-dictionary --pattern 'dictionary-server-*'
