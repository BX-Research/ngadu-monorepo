-- +goose Up
-- Ratings are deliberately reset; make seed supplies the shared global baseline.
-- Profiles and their private photo object keys are preserved.
DROP TABLE player_ratings;
DROP TABLE memberships;
DROP TABLE communities;

CREATE TABLE player_ratings (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    rating integer NOT NULL DEFAULT 1200,
    wins integer NOT NULL DEFAULT 0 CHECK (wins >= 0),
    losses integer NOT NULL DEFAULT 0 CHECK (losses >= 0),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX player_ratings_ladder_idx ON player_ratings (rating DESC, user_id ASC);

-- +goose Down
-- The discarded scoped ratings and role assignments cannot be reconstructed.
-- +goose StatementBegin
DO $$ BEGIN
    RAISE EXCEPTION 'The global ladder reset is forward-only; restore a backup to roll back';
END $$;
-- +goose StatementEnd
