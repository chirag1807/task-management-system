-- migrate:up
ALTER TABLE team_members ADD CONSTRAINT team_member_unique_constraint UNIQUE (team_id, member_id);

-- migrate:down
DROP INDEX team_member_unique_constraint CASCADE;