package utils

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func InsertMockData(tx pgx.Tx) (pgx.Tx, error) {
	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO users (id, first_name, last_name, bio, email, password, privacy) VALUES(954488202459119617, 'Dhyey', 'Panchal', 'Junior Software Engineer at Rapidops INC.', 'dhyey@gmail.com', '$2a$14$iCdRt4r2bigHcBxxDgxr/OOsjylBNwVrmQgsOcWgVwdjlZuJxtFNa', 'PUBLIC');")
	batch.Queue("INSERT INTO users (id, first_name, last_name, bio, email, password, privacy) VALUES(954497896847212545, 'Ridham', 'Chauhan', 'Junior Software Engineer at RiverEdge.', 'ridham@gmail.com', '$2a$14$8K8gJCgpqWwRTM86q0/bP.cSrlFEVuiy.0KlDBKzK6wmBtEhgV5Me', 'PUBLIC');")
	batch.Queue("INSERT INTO users (first_name, last_name, bio, email, password, privacy) VALUES('Aashutosh', 'Gupta', 'Junior Software Engineer at ZURU TECH INDIA', 'guptaaahutosh354@gmail.com', '$2a$14$FhDiMSnCN8sJ7Tb0UDBXn.bbKVYF3b4ZVwEwPXfAzvDgXZlC3B1g2', 'PUBLIC');")
	batch.Queue("INSERT INTO users (id, first_name, last_name, bio, email, password, privacy) VALUES(954497896847212546, 'Aashutosh', 'Gupta', 'Junior Software Engineer at ZURU TECH INDIA', 'guptaaahutosh355@gmail.com', '$2a$14$FhDiMSnCN8sJ7Tb0UDBXn.bbKVYF3b4ZVwEwPXfAzvDgXZlC3B1g2', 'PRIVATE');")
	batch.Queue("INSERT INTO teams (id, name, created_by, created_at, team_privacy) VALUES(954507580144451585, 'Team A', 954488202459119617, current_timestamp(), 'PUBLIC');")
	batch.Queue("INSERT INTO public.team_members (team_id, member_id) VALUES(954507580144451585, 954488202459119617);")
	batch.Queue("INSERT INTO teams (id, name, created_by, created_at, team_privacy) VALUES(954507580144451586, 'Team B', 954488202459119617, current_timestamp(), 'PRIVATE');")
	batch.Queue("INSERT INTO public.team_members (team_id, member_id) VALUES(954507580144451586, 954488202459119617);")
	batch.Queue("INSERT INTO tasks (id, title, description, deadline, assignee_team, status, priority, created_by, created_at) VALUES(954511608047501313, 'task3', 'this is task3', current_timestamp(), 954507580144451585, 'TO-DO', 'VERY HIGH', 954488202459119617, current_timestamp());")
	batch.Queue("INSERT INTO tasks (id, title, description, deadline, assignee_team, status, priority, created_by, created_at) VALUES(954511608047501314, 'task4', 'this is task3', current_timestamp(), 954507580144451585, 'CLOSED', 'VERY HIGH', 954488202459119617, current_timestamp());")
	batch.Queue("INSERT INTO refresh_tokens (user_id, refresh_token) VALUES (954488202459119617, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc');")
	batch.Queue("INSERT INTO otps (id, otp, otp_expire_time, email, is_verified) VALUES (954537852771565569, 1099, 'infinity', 'dhyey@gmail.com', true);")
	batch.Queue("INSERT INTO tasks (title, description, deadline, assignee_team, status, priority, created_by, created_at) VALUES('task2', 'this is task2', '2024-03-30T22:59:59.000Z', 954507580144451585, 'TO-DO', 'VERY HIGH', 954488202459119617, current_timestamp());")
	results := tx.SendBatch(context.Background(), batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(context.Background())
		log.Fatal(err)
		return tx, err
	}
	return tx, nil
}

func ClearMockData(dbConn *pgx.Conn) error {
	// query := "DELETE FROM tasks WHERE id <> 954511608047501313;" +
	// 	"DELETE FROM team_members WHERE NOT (team_id = 954507580144451585 AND member_id = 954488202459119617);" +
	// 	"DELETE FROM teams WHERE id <> 954507580144451585;" +
	// 	"DELETE FROM refresh_tokens WHERE user_id <> 954488202459119617;" +
	// 	"DELETE FROM users WHERE id NOT IN (954488202459119617, 954497896847212545);" +
	// 	"DELETE FROM otps WHERE id <> 954537852771565569;"
	
	query := "DELETE FROM tasks;" + "DELETE FROM team_members;" + "DELETE FROM teams;" +
		"DELETE FROM refresh_tokens;" + "DELETE FROM users;" + "DELETE FROM otps;"

	_, err := dbConn.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
