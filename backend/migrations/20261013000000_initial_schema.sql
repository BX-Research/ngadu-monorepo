-- +goose Up
-- IF NOT EXISTS preserves databases created by the earlier development seeder.
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    firebase_uid text NOT NULL UNIQUE,
    display_name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS communities (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS memberships (
    user_id uuid NOT NULL REFERENCES users(id),
    community_id uuid NOT NULL REFERENCES communities(id),
    role text NOT NULL CHECK (role IN ('member', 'organizer')),
    joined_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, community_id)
);

-- +goose Down
DROP TABLE memberships;
DROP TABLE communities;
DROP TABLE users;
