-- +goose Up
-- Avatar keys are derived as avatars/{userId}.jpg by the application.
ALTER TABLE users DROP COLUMN profile_photo_key;

-- +goose Down
-- Former object-key values cannot be recovered; restored references start empty.
ALTER TABLE users ADD COLUMN profile_photo_key text;
