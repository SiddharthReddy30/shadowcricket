package player

type Player struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	JerseyNumber int    `json:"jersey_number"`
	Role         string `json:"role"`
	IPLTeam        string `json:"ipl_team"`
	IsWicketKeeper bool   `json:"is_wicket_keeper"`
}

type Video struct {
	ID              int    `json:"id"`
	PlayerID        int    `json:"player_id"`
	RawVideo        string `json:"raw_video"`
	SilhouetteVideo string `json:"silhouette_video"`
}
