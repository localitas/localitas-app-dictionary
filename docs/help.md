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

### Version

```bash
./dictionary-server --version
```

### Build from source

```bash
# Development (native)
cd apps/dictionary && go build -o bin/dictionary-server ./cmd/dictionary-server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/dictionary-server-linux-amd64 ./cmd/dictionary-server
```

### Docker

Build a Docker image directly from the binary:

```bash
# Default base image (debian:12-slim)
./dictionary-server docker-build

# Custom base image
./dictionary-server docker-build --base ubuntu:24.04

# Custom Dockerfile
./dictionary-server docker-build --dockerfile ./my.Dockerfile

# Tag and push to registry
./dictionary-server docker-build --tag ghcr.io/localitas/dictionary:latest --push
```

The `docker-build` command requires a Linux amd64 binary in the same directory. Run `make deploy-build` from the project root first.

### Download

Pre-built binaries are available on the [GitHub releases page](https://github.com/localitas/localitas/releases).

Each release includes three builds per app:
- `dictionary-server-darwin-arm64` (macOS Apple Silicon)
- `dictionary-server-linux-amd64` (Linux x86_64)
- `dictionary-server-linux-arm64` (Linux ARM64)

Download with the GitHub CLI:

    gh release download --repo localitas/localitas --pattern 'dictionary-server-*'

### Release

All app binaries are published to GitHub releases as part of `make deploy-upload-image`.
