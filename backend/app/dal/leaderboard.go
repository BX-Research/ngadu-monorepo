package dal

import (
	"context"
	"encoding/json"
	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
)

type Leaderboard struct{ db Querier }

func NewLeaderboard(db Querier) *Leaderboard { return &Leaderboard{db: db} }

// Leaderboard uses one statement so the count and top ten share a snapshot.
func (l *Leaderboard) Leaderboard(ctx context.Context, communityID string) (types.Leaderboard, error) {
	result := types.Leaderboard{Players: []types.RankedPlayer{}}
	var data []byte
	err := l.db.QueryRow(ctx, `
 SELECT c.id::text, c.name,
 (SELECT count(*) FROM player_ratings WHERE community_id=c.id),
 COALESCE((SELECT jsonb_agg(to_jsonb(p) ORDER BY p.rating DESC, p.user_id ASC) FROM (
   SELECT r.user_id::text, u.display_name, r.rating, r.wins, r.losses, u.profile_photo_key AS photo_key
   FROM player_ratings r JOIN users u ON u.id=r.user_id
   WHERE r.community_id=c.id ORDER BY r.rating DESC, r.user_id ASC LIMIT 10
 ) p), '[]'::jsonb)
 FROM communities c WHERE c.id=$1`, communityID).Scan(&result.Community.ID, &result.Community.Name, &result.TotalPlayers, &data)
	if err != nil {
		return result, err
	}
	var players []struct {
		types.RankedPlayer
		Key *string `json:"photo_key"`
	}
	if err := json.Unmarshal(data, &players); err != nil {
		return result, err
	}
	for i, p := range players {
		p.Rank = i + 1
		if p.Key != nil {
			p.PhotoKey = *p.Key
		}
		result.Players = append(result.Players, p.RankedPlayer)
	}
	return result, nil
}
