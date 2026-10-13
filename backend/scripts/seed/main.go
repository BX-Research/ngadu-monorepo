package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"ngadu/backend/internal/devdb"
)

type seedUser struct {
	ID, FirebaseUID, DisplayName string
}

type seedCommunity struct {
	ID, Name string
}

type seedMembership struct {
	FirebaseUID, CommunityID, Role string
}

// FirebaseUID values are fictional. Replace them with development identities
// to test authenticated requests; this command does not create Firebase users.
var seedUsers = []seedUser{
	{"00000000-0000-4000-8000-000000000001", "seed-user-avery", "Avery Demo"},
	{"00000000-0000-4000-8000-000000000002", "seed-user-blake", "Blake Demo"},
	{"00000000-0000-4000-8000-000000000003", "seed-user-casey", "Casey Demo"},
	{"00000000-0000-4000-8000-000000000004", "seed-user-drew", "Drew Demo"},
}

var seedCommunities = []seedCommunity{
	{"00000000-0000-4000-8000-000000000101", "Demo Shuttle Club"},
	{"00000000-0000-4000-8000-000000000102", "Demo Weekend Badminton"},
}

// Drew intentionally has no memberships, for testing the nonmember state.
var seedMemberships = []seedMembership{
	{seedUsers[0].FirebaseUID, "00000000-0000-4000-8000-000000000101", "organizer"},
	{seedUsers[1].FirebaseUID, "00000000-0000-4000-8000-000000000101", "member"},
	{seedUsers[2].FirebaseUID, "00000000-0000-4000-8000-000000000101", "member"},
	{seedUsers[0].FirebaseUID, "00000000-0000-4000-8000-000000000102", "member"},
	{seedUsers[2].FirebaseUID, "00000000-0000-4000-8000-000000000102", "organizer"},
}

func seed(ctx context.Context, conn *pgx.Conn, output io.Writer) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin seed: %w", err)
	}
	defer tx.Rollback(context.Background())
	var messages bytes.Buffer
	userIDs := make(map[string]string, len(seedUsers))
	for _, user := range seedUsers {
		var id string
		err := tx.QueryRow(ctx, `INSERT INTO users (id, firebase_uid, display_name)
			VALUES ($1, $2, $3) ON CONFLICT (firebase_uid) DO NOTHING
			RETURNING id::text`, user.ID, user.FirebaseUID, user.DisplayName).Scan(&id)
		status := "Created"
		if errors.Is(err, pgx.ErrNoRows) {
			status = "Skip (exists)"
			err = tx.QueryRow(ctx, "SELECT id::text FROM users WHERE firebase_uid = $1", user.FirebaseUID).Scan(&id)
		}
		if err != nil {
			return fmt.Errorf("seed user %s: %w", user.FirebaseUID, err)
		}
		userIDs[user.FirebaseUID] = id
		fmt.Fprintf(&messages, "%s: user %s\n", status, user.FirebaseUID)
	}
	for _, community := range seedCommunities {
		result, err := tx.Exec(ctx, `INSERT INTO communities (id, name)
			VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`, community.ID, community.Name)
		if err != nil {
			return fmt.Errorf("seed community %s: %w", community.ID, err)
		}
		status := "Created"
		if result.RowsAffected() == 0 {
			status = "Skip (exists)"
		}
		fmt.Fprintf(&messages, "%s: community %s\n", status, community.Name)
	}
	for _, membership := range seedMemberships {
		result, err := tx.Exec(ctx, `INSERT INTO memberships (user_id, community_id, role)
			VALUES ($1, $2, $3) ON CONFLICT (user_id, community_id)
			DO UPDATE SET role = EXCLUDED.role
			WHERE memberships.role IS DISTINCT FROM EXCLUDED.role`,
			userIDs[membership.FirebaseUID], membership.CommunityID, membership.Role)
		if err != nil {
			return fmt.Errorf("seed membership for %s: %w", membership.FirebaseUID, err)
		}
		status := "Created/updated"
		if result.RowsAffected() == 0 {
			status = "Skip (unchanged)"
		}
		fmt.Fprintf(&messages, "%s: membership %s in %s → %s\n", status,
			membership.FirebaseUID, membership.CommunityID, membership.Role)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed: %w", err)
	}
	_, err = io.Copy(output, &messages)
	return err
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	connectionURL, err := devdb.URL()
	if err != nil {
		return err
	}
	conn, err := pgx.Connect(ctx, connectionURL)
	if err != nil {
		// Avoid printing a connection URL that may include database credentials.
		return errors.New("local development database unavailable; run make db-migrate")
	}
	defer conn.Close(context.Background())
	return seed(ctx, conn, os.Stdout)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Seed failed:", err)
		os.Exit(1)
	}
}
