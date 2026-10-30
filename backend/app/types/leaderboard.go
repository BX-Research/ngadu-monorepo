package types

type RankedPlayer struct {
	Rank        int     `json:"rank"`
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Rating      int     `json:"rating"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	PhotoURL    *string `json:"photo_url"`
}

type Leaderboard struct {
	TotalPlayers int            `json:"total_players"`
	Players      []RankedPlayer `json:"players"`
}
