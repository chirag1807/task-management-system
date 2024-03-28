package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TeamRepository interface {
	CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (int64, error)
	AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error
	RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error
	GetAllTeams(userID int64, flag int, queryParams request.TeamQueryParams) ([]response.Team, error)
	//flag is used for get my created teams and get teams in which i was added.
	GetTeamMembers(teamID int64, queryParams request.TeamQueryParams) ([]response.User, error)
	LeftTeam(userID int64, teamID int64) error
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

func (t teamRepository) CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (int64, error) {
	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	var teamId int64
	err = tx.QueryRow(ctx, `INSERT INTO teams (name, team_profile, created_by) VALUES ($1, $2, $3) RETURNING id`, teamToCreate.Name, teamToCreate.TeamProfile, teamToCreate.CreatedBy).Scan(&teamId)
	if err != nil {
		tx.Rollback(ctx)
		return teamId, err
	}

	batch := &pgx.Batch{}
	for _, v := range teamMembers.MemberID {
		batch.Queue(`INSERT INTO team_members (team_id, member_id) VALUES ($1, $2)`, teamId, v)
	}
	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(ctx)
		return teamId, err
	}
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return teamId, err
	}

	for _, v := range teamMembers.MemberID {
		t.redisClient.SAdd(ctx, "user:"+strconv.FormatInt(v, 10)+":teams", teamId)
	}

	return teamId, nil
}

func (t teamRepository) AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error {
	var dbTeamCreatedBy int64
	t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToAdd.TeamID).Scan(&dbTeamCreatedBy)

	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	var args []interface{}
	query := `SELECT profile FROM users WHERE id IN (`
	for i, v := range teamMembersToAdd.MemberID {
		query += `$` + strconv.Itoa(i+1) + `, `
		args = append(args, v)
	}
	query = query[:len(query)-2]
	query += `)`

	users, err := t.dbConn.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}
	defer users.Close()

	var user response.User
	for users.Next() {
		if err := users.Scan(&user.Profile); err != nil {
			return err
		}
		if user.Profile == "Private" {
			return errorhandling.OnlyPublicMemberAllowed
		}
	}

	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, v := range teamMembersToAdd.MemberID {
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

	for _, v := range teamMembersToAdd.MemberID {
		t.redisClient.SAdd(ctx, "user:"+strconv.FormatInt(v, 10)+":teams", teamMembersToAdd.TeamID)
	}

	return nil
}

func (t teamRepository) RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error {
	var dbTeamCreatedBy int64
	t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToRemove.TeamID).Scan(&dbTeamCreatedBy)

	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	var args []interface{}
	query := `DELETE FROM team_members WHERE member_id IN (`
	for i, v := range teamMembersToRemove.MemberID {
		query += `$` + strconv.Itoa(i+1) + `, `
		args = append(args, v)
	}
	query = query[:len(query)-2]
	query += `)`

	_, err := t.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	for _, v := range teamMembersToRemove.MemberID {
		t.redisClient.SRem(context.Background(), "user:"+strconv.FormatInt(v, 10)+":teams", teamMembersToRemove.TeamID)
	}

	return nil
}

func (t teamRepository) GetAllTeams(userID int64, flag int, queryParams request.TeamQueryParams) ([]response.Team, error) {
	//flag = 0 => created by me, flag = 1 => i am member
	var teams pgx.Rows
	var err error
	teamsSlice := make([]response.Team, 0)

	var query string
	if flag == 0 {
		query = `SELECT * FROM teams WHERE created_by = $1 AND true`
		query = CreateQueryForParamsOfGetTeam(query, queryParams)
		teams, err = t.dbConn.Query(context.Background(), query, userID)
	}
	if flag == 1 {
		query = `SELECT * FROM teams WHERE id IN (SELECT team_id from team_members where member_id = $1)`
		query = CreateQueryForParamsOfGetTeam(query, queryParams)
		teams, err = t.dbConn.Query(context.Background(), query, userID)
	}

	if err != nil {
		return teamsSlice, err
	}
	defer teams.Close()

	var team response.Team
	for teams.Next() {
		if err := teams.Scan(&team.ID, &team.Name, &team.CreatedBy, &team.CreatedAt, &team.TeamProfile); err != nil {
			return teamsSlice, err
		}
		teamsSlice = append(teamsSlice, team)
	}

	return teamsSlice, nil
}

func (t teamRepository) GetTeamMembers(teamID int64, queryParams request.TeamQueryParams) ([]response.User, error) {
	var teamMembers pgx.Rows
	var err error
	teamMembersSlice := make([]response.User, 0)

	query := `SELECT id, first_name, last_name, bio, email, profile FROM users WHERE id IN (SELECT member_id from team_members where team_id = $1)`
	query = CreateQueryForParamsOfGetTeam(query, queryParams)
	teamMembers, err = t.dbConn.Query(context.Background(), query, teamID)

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

func CreateQueryForParamsOfGetTeam(query string, queryParams request.TeamQueryParams) string {
	if queryParams.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE '%%%s%%')", queryParams.Search)
	}
	if queryParams.SortByCreateAt {
		query += " ORDER BY created_at"
	}
	query += fmt.Sprintf(" LIMIT %d", queryParams.Limit)
	query += fmt.Sprintf(" OFFSET %d", queryParams.Offset)
	return query
}

func (t teamRepository) LeftTeam(userID int64, teamID int64) error {
	a, err := t.dbConn.Exec(context.Background(), "DELETE FROM team_members WHERE member_id = $1 AND team_id = $2", userID, teamID)
	if a.RowsAffected() == 0 {
		return errorhandling.NotAMember
	}
	if err != nil {
		return err
	}
	t.redisClient.SRem(context.Background(), "user:"+strconv.FormatInt(userID, 10)+":teams", teamID)
	return nil
}
