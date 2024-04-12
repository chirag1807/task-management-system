-- migrate:up
ALTER TABLE team_members ADD CONSTRAINT team_member_unique_constraint UNIQUE (team_id, member_id);

-- migrate:down
ALTER TABLE team_members DROP CONSTRAINT team_member_unique_constraint;
