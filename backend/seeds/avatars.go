// Package seeds contains repository-owned, fictional development fixtures.
package seeds

import (
	"embed"
	"sort"
)

//go:embed photos/*.jpg
var photos embed.FS

// Some players intentionally have no fixture to exercise the SVG placeholder.
var avatarFiles = map[string]string{
	"seed-player-201": "coral", "seed-player-202": "sage", "seed-player-203": "blue",
	"seed-player-206": "blue", "seed-player-207": "coral", "seed-player-209": "blue",
	"seed-player-210": "coral", "seed-player-211": "sage", "seed-player-213": "coral",
	"seed-player-214": "sage", "seed-player-215": "blue", "seed-player-217": "sage",
	"seed-player-218": "blue", "seed-player-219": "coral",
}

type AvatarFixture struct {
	FirebaseUID string
	JPEG        []byte
}

func AvatarFixtures() ([]AvatarFixture, error) {
	uids := make([]string, 0, len(avatarFiles))
	for uid := range avatarFiles {
		uids = append(uids, uid)
	}
	sort.Strings(uids)
	result := make([]AvatarFixture, 0, len(uids))
	for _, uid := range uids {
		data, err := photos.ReadFile("photos/" + avatarFiles[uid] + ".jpg")
		if err != nil {
			return nil, err
		}
		result = append(result, AvatarFixture{FirebaseUID: uid, JPEG: data})
	}
	return result, nil
}
