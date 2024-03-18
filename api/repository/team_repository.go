package repository

import (
	"context"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TeamRepository interface {
	CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) error
	AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error
	RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error
	GetAllTeams(userID int64, flag int) ([]response.Team, error)
	//flag is used for get my created teams and get teams in which i was added.
	GetTeamMembers(teamID int64) ([]response.User, error)
	LeftTeam(userID int64, teamID int64) (error)
}

type teamRepository struct {
	dbConn      *pgx.Conn
	redisClient *redis.Client
}

func NewTeamRepo(dbConn *pgx.Conn, redisClient *redis.Client) TeamRepository {
	return teamRepository{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (t teamRepository) CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) error {
	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	var teamId int64
	err = tx.QueryRow(ctx, `INSERT INTO teams (name, created_by) VALUES ($1, $2) RETURNING id`, teamToCreate.Name, teamToCreate.CreatedBy).Scan(&teamId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	batch := &pgx.Batch{}
	for _, v := range teamMembers.MemberID {
		batch.Queue(`INSERT INTO team_members (team_id, member_id) VALUES ($1, $2)`, teamId, v)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (t teamRepository) AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error {
	var dbTeamCreatedBy int64
	t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToAdd.TeamID).Scan(&dbTeamCreatedBy)

	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	var memberProfile string
	for _, v := range teamMembersToAdd.MemberID {
		t.dbConn.QueryRow(context.Background(), `SELECT profile FROM users WHERE id = $1`, v).Scan(&memberProfile)
		if memberProfile == "Private"{
			return errorhandling.OnlyPublicMemberAllowed
		}
		batch.Queue(`INSERT INTO team_members (team_id, member_id) VALUES ($1, $2)`, teamMembersToAdd.TeamID, v)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(ctx)
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return errorhandling.MemberExist
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return errorhandling.MemberExist
		}
		return err
	}

	return nil
}

func (t teamRepository) RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error {
	var dbTeamCreatedBy int64
	t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToRemove.TeamID).Scan(&dbTeamCreatedBy)

	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	counter := 1
	var args []interface{}
	query := `DELETE FROM team_members WHERE member_id IN (`
	for _, v := range teamMembersToRemove.MemberID {
		query += `$` + strconv.Itoa(counter) + `, `
		counter++
		args = append(args, v)
	}
	query += `)`

	_, err := t.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t teamRepository) GetAllTeams(userID int64, flag int) ([]response.Team, error) {
	//flag = 0 => created by me, flag = 1 => i am member
	var teams pgx.Rows
	var err error
	teamsSlice := make([]response.Team, 0)

	if flag == 0 {
		teams, err = t.dbConn.Query(context.Background(), `SELECT * FROM teams WHERE created_by = $1`, userID)
	} else if flag == 1 {
		teams, err = t.dbConn.Query(context.Background(), `SELECT * FROM teams WHERE id IN (SELECT team_id from team_members where member_id = $1)`, userID)
	} else {
		return teamsSlice, errorhandling.FlagNotProvided
	}

	if err != nil {
		return teamsSlice, err
	}
	defer teams.Close()

	var team response.Team
	for teams.Next() {
		if err := teams.Scan(&team.ID, &team.Name, &team.CreatedBy, &team.CreatedAt); err != nil {
			return teamsSlice, err
		}
		teamsSlice = append(teamsSlice, team)
	}

	return teamsSlice, nil
}

func (t teamRepository) GetTeamMembers(teamID int64) ([]response.User, error) {
	var teamMembers pgx.Rows
	var err error
	teamMembersSlice := make([]response.User, 0)

	teamMembers, err = t.dbConn.Query(context.Background(), `SELECT id, first_name, last_name, bio, email, profile FROM users WHERE id IN (SELECT member_id from team_members where team_id = $1)`, teamID)

	if err != nil {
		return teamMembersSlice, err
	}
	defer teamMembers.Close()

	var teamMember response.User
	for teamMembers.Next() {
		if err := teamMembers.Scan(&teamMember.ID, &teamMember.FirstName, &teamMember.LastName, &teamMember.Bio, &teamMember.Email, &teamMember.Profile); err != nil {
			return teamMembersSlice, err
		}
		teamMembersSlice = append(teamMembersSlice, teamMember)
	}

	return teamMembersSlice, nil
}

func (t teamRepository) LeftTeam(userID int64, teamID int64) (error) {
	a, err := t.dbConn.Exec(context.Background(), "DELETE FROM team_members WHERE member_id = $1 AND team_id = $2", userID, teamID)
	if a.RowsAffected() == 0 {
		return errorhandling.NotAMember
	}
	if err != nil {
		return err
	}
	return nil
}