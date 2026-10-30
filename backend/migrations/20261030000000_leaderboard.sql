-- +goose Up
ALTER TABLE users ADD COLUMN profile_photo_key text;

CREATE TABLE player_ratings (
    community_id uuid NOT NULL,
    user_id uuid NOT NULL,
    rating integer NOT NULL DEFAULT 1200,
    wins integer NOT NULL DEFAULT 0 CHECK (wins >= 0),
    losses integer NOT NULL DEFAULT 0 CHECK (losses >= 0),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (community_id, user_id),
    FOREIGN KEY (user_id, community_id) REFERENCES memberships(user_id, community_id) ON DELETE CASCADE
);
CREATE INDEX player_ratings_ladder_idx ON player_ratings (community_id, rating DESC, user_id ASC);

-- +goose Down
DROP TABLE player_ratings;
ALTER TABLE users DROP COLUMN profile_photo_key;
