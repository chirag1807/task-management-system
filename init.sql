CREATE DATABASE IF NOT EXISTS taskmanager;

USE taskmanager;

CREATE TYPE userprofile AS ENUM ('Public', 'Private');

CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name VARCHAR(255) NOT NULL, last_name VARCHAR(255) NOT NULL, bio VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL, password VARCHAR(255), profile userprofile NOT NULL DEFAULT 'Public', UNIQUE (email));

CREATE TABLE IF NOT EXISTS otps (id SERIAL PRIMARY KEY, otp INT8 NOT NULL, otp_expire_time TIMESTAMP WITHOUT TIME ZONE NOT NULL );

CREATE TABLE IF NOT EXISTS refresh_tokens (user_id INT64 NOT NULL REFERENCES users (id), refresh_token string NOT NULL);

CREATE TABLE IF NOT EXISTS teams (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, created_by INT64 NOT NULL REFERENCES users (id),
created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, team_profile userprofile NOT NULL DEFAULT 'Public');

CREATE TABLE IF NOT EXISTS team_members (team_id INT64 NOT NULL REFERENCES teams (id), member_id INT64 NOT NULL REFERENCES users (id));

CREATE TYPE taskstatus AS ENUM ('TO-DO', 'In-Progress', 'Completed', 'Closed');

CREATE TYPE taskpriority AS ENUM ('Low', 'Medium', 'High', 'Very High');

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
