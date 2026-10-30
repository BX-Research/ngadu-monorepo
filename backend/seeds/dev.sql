-- Local development fixtures only. Firebase UIDs are fictional placeholders.
BEGIN;

INSERT INTO users (id, firebase_uid, display_name) VALUES
    ('00000000-0000-4000-8000-000000000001', 'seed-user-avery', 'Avery Demo'),
    ('00000000-0000-4000-8000-000000000002', 'seed-user-blake', 'Blake Demo'),
    ('00000000-0000-4000-8000-000000000003', 'seed-user-casey', 'Casey Demo'),
    ('00000000-0000-4000-8000-000000000004', 'seed-user-drew', 'Drew Demo')
ON CONFLICT (firebase_uid) DO NOTHING;

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

-- Snapshot ratings: local demo match results do not update these records.
INSERT INTO player_ratings (user_id, rating, wins, losses)
SELECT u.id, f.rating, f.wins, f.losses
FROM (VALUES
    ('seed-user-avery', 1400, 15, 8),
    ('seed-user-blake', 1300, 12, 11),
    ('seed-user-casey', 1200, 10, 14),
    ('seed-user-drew', 1200, 0, 0),
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
ON CONFLICT (user_id) DO NOTHING;

COMMIT;
