-- migrate:up
CREATE TYPE privacy AS ENUM ('PUBLIC', 'PRIVATE');

CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name VARCHAR(255) NOT NULL, last_name VARCHAR(255) NOT NULL, bio VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL, password VARCHAR(255), privacy privacy NOT NULL DEFAULT 'PUBLIC', UNIQUE (email));

CREATE TABLE IF NOT EXISTS otps (id SERIAL PRIMARY KEY, otp INT8 NOT NULL, otp_expire_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
email VARCHAR(255) NOT NULL, is_verified BOOLEAN DEFAULT false);

CREATE TABLE IF NOT EXISTS refresh_tokens (user_id INT64 NOT NULL REFERENCES users (id), refresh_token string NOT NULL);

CREATE TABLE IF NOT EXISTS teams (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, created_by INT64 NOT NULL REFERENCES users (id),
created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, team_privacy privacy NOT NULL DEFAULT 'PUBLIC');

CREATE TABLE IF NOT EXISTS team_members (team_id INT64 NOT NULL REFERENCES teams (id), member_id INT64 NOT NULL REFERENCES users (id));

CREATE TYPE taskstatus AS ENUM ('TO-DO', 'IN-PROGRESS', 'COMPLETED', 'CLOSED');

CREATE TYPE taskpriority AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'VERY HIGH');

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    description VARCHAR(1000),
    deadline TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    assignee_individual INT64 REFERENCES users (id),
    assignee_team INT64 REFERENCES teams (id),
    status taskstatus NOT NULL,
    priority taskpriority NOT NULL,
    created_by INT64 NOT NULL REFERENCES users (id),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INT64 REFERENCES users (id),
    updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS index_fetch_tasks ON tasks (title, description, status);

-- migrate:down
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TYPE IF EXISTS taskstatus;
DROP TYPE IF EXISTS taskpriority;
DROP TABLE IF EXISTS refresh_tokens;
DROP INDEX IF EXISTS index_fetch_tasks;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS privacy;
DROP TABLE IF EXISTS otps;