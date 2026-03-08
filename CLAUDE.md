# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ShadowCricket is a cricket silhouette guessing game. Users see a silhouette video of a cricketer and guess their identity, receiving color-coded feedback (green/yellow/white) on fields like name, country, jersey number, role, and IPL team. Unlimited guesses, random video each round.

## Stack

- **Backend**: Go (net/http stdlib, Go 1.23)
- **Frontend**: React + Vite + TypeScript
- **Data**: JSON files loaded into memory (no DB)
- **Media**: Silhouette videos served via blob storage (video IDs in API, not filenames)
- **Production**: Single Go binary with React build embedded via `embed.FS`
- **Deployment**: Docker multi-stage build
- **Tooling**: mise for Go/Node version management
- **Dependencies**: godotenv (env file loading), gotestsum (test output formatting)

## Project Structure

```
cmd/server/main.go              # Entry point, embeds React build (not yet created)
internal/
  config/config.go              # Env config using godotenv (.env.dev / .env.prod)
  player/model.go               # Player and Video structs
  player/store.go               # Load JSON data, search, lookup, random video pick
  game/
    data_models.go              # Types, constants, color enums, lookup maps
    token.go                    # AES-GCM token creation and decryption
    compare.go                  # Field comparison functions (name, country, jersey, role, IPL)
    guess.go                    # EvaluateGuess orchestrator
    test_helpers_test.go        # Shared test player data and assertion helpers
    compare_test.go             # Compare function tests (13 tests)
    token_test.go               # Token tests (6 tests)
    guess_test.go               # EvaluateGuess tests (2 tests)
  api/                          # Not yet created
    router.go                   # Route registration
    middleware.go               # CORS, logging
    handlers_game.go            # GET /api/game/random, POST /api/game/guess
    handlers_players.go         # GET /api/players/search?q=
    response.go                 # JSON helpers
data/
  players.json                  # Player details (5 sample cricketers)
  videos.json                   # Video metadata linked to players (7 entries)
  raw_videos/                   # Input cricket video clips
  silhouette_videos/            # Processed silhouette output videos
web/                            # React + Vite frontend (not yet created)
```

## Config

- Uses godotenv to load `.env.dev` or `.env.prod`
- `PORT` (default: 8080), `DATA_DIR` (default: data), `TOKEN_SECRET` (required, exactly 32 bytes)
- AES-256 enforced — server refuses to start if TOKEN_SECRET is not 32 bytes
- `.env` files are in `.gitignore` (secrets never committed)

## Data Design

Two JSON files (like two database tables):
- **`players.json`** — player details: id, name, country, jersey_number, role, ipl_team, is_wicket_keeper
- **`videos.json`** — video metadata: id, player_id, raw_video, silhouette_video

One player can have multiple videos. The game picks a random video, not a random player.

### Player Roles
- Opening Batsman (rank 1)
- Middle-Order Batsman (rank 2)
- Finisher (rank 3)
- All-Rounder (rank 4)
- Bowler (rank 5)

Wicket-keeping is a separate boolean field, not a role.

## Feedback Colors

- **Green** = exact match (or both are wicket-keepers for role field)
- **Yellow** = close/nearby match
- **White** = no match

### Field-specific rules
- **Name**: green (exact) / white
- **Country**: green (same country) / yellow (same continent) / white
- **Jersey Number**: green (exact) / yellow (within 5) / white
- **Role**: if both wicket-keepers → green "Wicket-Keeper"; otherwise compare by rank: green (same) / yellow (1 apart) / white (2+ apart)
- **IPL Team**: green (same team) / yellow (both in IPL, different teams) / white (one or both not in IPL)

## Token Design

- AES-256-GCM encrypted tokens (32-byte key enforced)
- Token encodes the target player ID
- Server is stateless — no sessions, no database
- Each token is unique (random nonce) even for the same player
- Token is created once per round, sent back with every guess

## Development Setup

### Run tests
```bash
mise run verbosetest                              # all tests, colored output
mise run verbosetest ./internal/game/             # specific package
mise run verbosetest -- ./internal/game/ -run TestCompare  # specific tests
go test ./...                                     # plain output (if mise activated)
```

### Go backend
```bash
go run ./cmd/server
```

### Frontend dev server (with proxy to Go backend)
```bash
cd web && npm install && npm run dev
```

### Docker
```bash
docker-compose up --build              # dev
docker-compose -f docker-compose.prod.yml up --build  # prod
```

## API Endpoints

- `GET /api/health` — health check
- `GET /api/game/random` — picks random video, returns `{token, video_id}`
- `POST /api/game/guess` — body: `{token, player_id}`, returns feedback + correct flag
- `GET /api/players/search?q=` — autocomplete, returns `[{id, name}]` (max 10)

## Key Design Decisions

- Users can only guess players in the database — frontend disables the guess button until a player is selected from autocomplete results, backend validates player_id exists (returns 400 if not found)
- Random per visit (not daily puzzle)
- AES-256-GCM encrypted token for round state (server stays stateless)
- Unlimited guesses, repeats allowed
- Videos served from blob storage (video IDs in API, not filenames)
- No external router — Go 1.22+ stdlib only
- Video metadata separate from player data (one player → many videos)
- Wicket-keeping is a boolean ability, not a role position
- Game package split by responsibility: data_models, token, compare, guess
- Test files split to match: compare_test, token_test, guess_test, test_helpers_test

## Current Status

**Completed:**
- Core game logic (player models, data store, tokens, guess evaluation)
- Config package with godotenv
- 21 unit tests (all passing)

**Next:**
- API layer (response.go, middleware.go, handlers, router.go, main.go)
- Then: React frontend, Docker setup
