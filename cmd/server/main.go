package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/siddharthreddy/shadowcricket/internal/api"
	"github.com/siddharthreddy/shadowcricket/internal/config"
	"github.com/siddharthreddy/shadowcricket/internal/player"
)

func main() {
	cfg, err := config.Load(".env.dev")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	store, err := player.LoadStore(
		filepath.Join(cfg.DataDir, "players.json"),
		filepath.Join(cfg.DataDir, "videos.json"),
	)
	if err != nil {
		log.Fatalf("data: %v", err)
	}

	handler := api.NewRouter(store, cfg.TokenSecret)

	log.Printf("server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}
