-- Local development fixtures only. Firebase UIDs are fictional placeholders.
BEGIN;

INSERT INTO users (id, firebase_uid, display_name) VALUES
    ('00000000-0000-4000-8000-000000000001', 'seed-user-avery', 'Avery Demo'),
    ('00000000-0000-4000-8000-000000000002', 'seed-user-blake', 'Blake Demo'),
    ('00000000-0000-4000-8000-000000000003', 'seed-user-casey', 'Casey Demo'),
    ('00000000-0000-4000-8000-000000000004', 'seed-user-drew', 'Drew Demo')
ON CONFLICT (firebase_uid) DO NOTHING;

INSERT INTO communities (id, name) VALUES
    ('00000000-0000-4000-8000-000000000101', 'Demo Shuttle Club'),
    ('00000000-0000-4000-8000-000000000102', 'Demo Weekend Badminton')
ON CONFLICT (id) DO NOTHING;

-- Resolve stored user IDs by Firebase UID; Drew intentionally has no membership.
INSERT INTO memberships (user_id, community_id, role)
SELECT users.id, fixtures.community_id::uuid, fixtures.role
FROM (VALUES
    ('seed-user-avery', '00000000-0000-4000-8000-000000000101', 'organizer'),
    ('seed-user-blake', '00000000-0000-4000-8000-000000000101', 'member'),
    ('seed-user-casey', '00000000-0000-4000-8000-000000000101', 'member'),
    ('seed-user-avery', '00000000-0000-4000-8000-000000000102', 'member'),
    ('seed-user-casey', '00000000-0000-4000-8000-000000000102', 'organizer')
) AS fixtures(firebase_uid, community_id, role)
JOIN users USING (firebase_uid)
ON CONFLICT (user_id, community_id) DO UPDATE SET role = EXCLUDED.role
WHERE memberships.role IS DISTINCT FROM EXCLUDED.role;

COMMIT;
