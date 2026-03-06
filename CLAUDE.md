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

## Project Structure

```
cmd/server/main.go          # Entry point, embeds React build
internal/
  config/config.go          # Env config (port, data path)
  player/model.go           # Player and Video structs
  player/store.go           # Load JSON data, search, lookup, random video pick
  game/guess.go             # AES-GCM tokens, EvaluateGuess -> color-coded feedback
  game/guess_test.go        # Tests for guess evaluation
  api/router.go             # Route registration
  api/middleware.go          # CORS, logging
  api/handlers_game.go      # GET /api/game/random, POST /api/game/guess
  api/handlers_players.go   # GET /api/players/search?q=
  api/response.go           # JSON helpers
  media/server.go           # http.FileServer for /media/
data/
  players.json              # Player details (one entry per player)
  videos.json               # Video metadata linked to players (many per player)
  raw_videos/               # Input cricket video clips
  silhouette_videos/        # Processed silhouette output videos
web/                        # React + Vite frontend
```

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

- AES-GCM encrypted tokens (not base64 + HMAC)
- Token encodes the target player ID
- Server is stateless — no sessions, no database
- Each token is unique (random nonce) even for the same player
- Token is created once per round, sent back with every guess
- Secret key must be exactly 16, 24, or 32 bytes (AES requirement)

## Development Setup

### Go backend
```bash
mise exec -- go run ./cmd/server
```

### Run tests
```bash
mise exec -- go test ./...
```

### Frontend dev server (with proxy to Go backend)
```bash
cd web && npm install && npm run dev
```

### Docker (production)
```bash
docker-compose -f docker-compose.prod.yml up --build
```

### Docker (dev)
```bash
docker-compose up --build
```

## API Endpoints

- `GET /api/health` — health check
- `GET /api/game/random` — picks random video, returns `{token, video_url}`
- `POST /api/game/guess` — body: `{token, player_id}`, returns feedback + correct flag
- `GET /api/players/search?q=` — autocomplete, returns `[{id, name}]` (max 10)

## Key Design Decisions

- Random per visit (not daily puzzle)
- AES-GCM encrypted token for round state (server stays stateless)
- Unlimited guesses, repeats allowed
- Videos served from blob storage (video IDs in API, not filenames)
- No external router — Go 1.22+ stdlib only
- Video metadata separate from player data (one player → many videos)
- Wicket-keeping is a boolean ability, not a role position
