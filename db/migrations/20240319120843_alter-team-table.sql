-- migrate:up
ALTER TABLE teams ADD COLUMN team_profile userprofile NOT NULL DEFAULT 'Public';

-- migrate:down
ALTER TABLE teams DROP COLUMN team_profile;
