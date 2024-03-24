-- migrate:up
CREATE FUNCTION update_assignee_trigger() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.assignee_team IS NOT NULL THEN
      OLD.assignee_individual := NULL;
    ELSIF NEW.assignee_team IS NOT NULL THEN
      OLD.assignee_individual := NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_assignee_trigger BEFORE UPDATE
    ON tasks FOR EACH ROW EXECUTE PROCEDURE update_assignee_trigger();


-- migrate:down
DROP TRIGGER IF EXISTS update_assignee_trigger ON tasks;
DROP FUNCTION IF EXISTS update_assignee_trigger
