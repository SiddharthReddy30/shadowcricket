package player

import (
	"encoding/json"
	"math/rand/v2"
	"os"
	"strings"
)

type Store struct {
	players []Player
	videos  []Video
	byID    map[int]Player
}

func LoadStore(playersFile string, videosFile string) (*Store, error) {
	playerData, err := os.ReadFile(playersFile)
	if err != nil {
		return nil, err
	}

	var players []Player
	if err := json.Unmarshal(playerData, &players); err != nil {
		return nil, err
	}

	videoData, err := os.ReadFile(videosFile)
	if err != nil {
		return nil, err
	}

	var videos []Video
	if err := json.Unmarshal(videoData, &videos); err != nil {
		return nil, err
	}

	byID := make(map[int]Player, len(players))
	for _, p := range players {
		byID[p.ID] = p
	}

	return &Store{players: players, videos: videos, byID: byID}, nil
}

func (s *Store) RandomVideo() (Video, Player) {
	v := s.videos[rand.IntN(len(s.videos))]
	p := s.byID[v.PlayerID]
	return v, p
}

func (s *Store) Lookup(id int) (Player, bool) {
	p, ok := s.byID[id]
	return p, ok
}

func (s *Store) Search(query string) []Player {
	query = strings.ToLower(query)
	var results []Player
	for _, p := range s.players {
		if strings.Contains(strings.ToLower(p.Name), query) {
			results = append(results, p)
			if len(results) >= 10 {
				break
			}
		}
	}
	return results
}
