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

-- Fictional leaderboard profiles. Insert only: reruns preserve edited names/photos.
INSERT INTO users (id, firebase_uid, display_name, profile_photo_key) VALUES
    ('00000000-0000-4000-8000-000000000201', 'seed-player-201', 'Adi Demo', 'profiles/coral.svg'),
    ('00000000-0000-4000-8000-000000000202', 'seed-player-202', 'Bella Demo', 'profiles/sage.svg'),
    ('00000000-0000-4000-8000-000000000203', 'seed-player-203', 'Citra Demo', 'profiles/blue.svg'),
    ('00000000-0000-4000-8000-000000000204', 'seed-player-204', 'Dani Demo', NULL),
    ('00000000-0000-4000-8000-000000000205', 'seed-player-205', 'Eka Demo', 'profiles/missing.svg'),
    ('00000000-0000-4000-8000-000000000206', 'seed-player-206', 'Farah Demo', 'profiles/blue.svg'),
    ('00000000-0000-4000-8000-000000000207', 'seed-player-207', 'Gilang Demo', 'profiles/coral.svg'),
    ('00000000-0000-4000-8000-000000000208', 'seed-player-208', 'Hana Demo', NULL),
    ('00000000-0000-4000-8000-000000000209', 'seed-player-209', 'Indra Demo', 'profiles/blue.svg'),
    ('00000000-0000-4000-8000-000000000210', 'seed-player-210', 'Jihan Demo', 'profiles/coral.svg'),
    ('00000000-0000-4000-8000-000000000211', 'seed-player-211', 'Kiki Demo', 'profiles/sage.svg'),
    ('00000000-0000-4000-8000-000000000212', 'seed-player-212', 'Laras Demo', NULL),
    ('00000000-0000-4000-8000-000000000213', 'seed-player-213', 'Miko Demo', 'profiles/coral.svg'),
    ('00000000-0000-4000-8000-000000000214', 'seed-player-214', 'Nina Demo', 'profiles/sage.svg'),
    ('00000000-0000-4000-8000-000000000215', 'seed-player-215', 'Oki Demo', 'profiles/blue.svg'),
    ('00000000-0000-4000-8000-000000000216', 'seed-player-216', 'Putri Demo', NULL),
    ('00000000-0000-4000-8000-000000000217', 'seed-player-217', 'Rafi Demo', 'profiles/sage.svg'),
    ('00000000-0000-4000-8000-000000000218', 'seed-player-218', 'Sari Demo', 'profiles/blue.svg'),
    ('00000000-0000-4000-8000-000000000219', 'seed-player-219', 'Tono Demo', 'profiles/coral.svg'),
    ('00000000-0000-4000-8000-000000000220', 'seed-player-220', 'Vina Demo', NULL)
ON CONFLICT (firebase_uid) DO NOTHING;

INSERT INTO memberships (user_id, community_id, role)
SELECT id, '00000000-0000-4000-8000-000000000101'::uuid, 'member'
FROM users WHERE firebase_uid IN (
    'seed-player-201',
    'seed-player-202',
    'seed-player-203',
    'seed-player-204',
    'seed-player-205',
    'seed-player-206',
    'seed-player-207',
    'seed-player-208',
    'seed-player-209',
    'seed-player-210',
    'seed-player-211',
    'seed-player-212',
    'seed-player-213',
    'seed-player-214',
    'seed-player-215',
    'seed-player-216',
    'seed-player-217',
    'seed-player-218',
    'seed-player-219',
    'seed-player-220')
ON CONFLICT (user_id, community_id) DO NOTHING;

-- Snapshot ratings: local demo match results do not update these records.
INSERT INTO player_ratings (community_id, user_id, rating, wins, losses)
SELECT '00000000-0000-4000-8000-000000000101'::uuid, u.id, f.rating, f.wins, f.losses
FROM (VALUES
    ('seed-user-avery', 1400, 15, 8),
    ('seed-user-blake', 1300, 12, 11),
    ('seed-user-casey', 1200, 10, 14),
    ('seed-player-201', 1750, 30, 5),
    ('seed-player-202', 1690, 29, 6),
    ('seed-player-203', 1690, 28, 7),
    ('seed-player-204', 1640, 27, 8),
    ('seed-player-205', 1610, 26, 9),
    ('seed-player-206', 1570, 25, 10),
    ('seed-player-207', 1530, 24, 11),
    ('seed-player-208', 1490, 23, 12),
    ('seed-player-209', 1450, 22, 13),
    ('seed-player-210', 1410, 21, 14),
    ('seed-player-211', 1370, 20, 15),
    ('seed-player-212', 1330, 19, 16),
    ('seed-player-213', 1290, 18, 17),
    ('seed-player-214', 1250, 17, 18),
    ('seed-player-215', 1210, 16, 19),
    ('seed-player-216', 1170, 15, 20),
    ('seed-player-217', 1130, 14, 21),
    ('seed-player-218', 1090, 13, 22),
    ('seed-player-219', 1050, 12, 23),
    ('seed-player-220', 1010, 11, 24)
) AS f(firebase_uid, rating, wins, losses)
JOIN users u USING (firebase_uid)
ON CONFLICT (community_id, user_id) DO NOTHING;

-- The second community demonstrates independent rankings for the same users.
INSERT INTO player_ratings (community_id, user_id, rating, wins, losses)
SELECT '00000000-0000-4000-8000-000000000102'::uuid, u.id, f.rating, 2, 1
FROM (VALUES ('seed-user-avery', 1250), ('seed-user-casey', 1280)) AS f(firebase_uid, rating)
JOIN users u USING (firebase_uid)
ON CONFLICT (community_id, user_id) DO NOTHING;

COMMIT;
