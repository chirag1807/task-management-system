package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TeamRepository interface {
	CreateTeam(teamToCreate request.Team, teamMembers []int64) (int64, error)
	AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembersWithTeamID) error
	RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembersWithTeamID) error
	GetAllTeams(userID int64, queryParams request.TeamQueryParams) ([]response.Team, error)
	//flag is used for get my created teams and get teams in which i was added.
	GetTeamMembers(teamId int64, queryParams request.TeamQueryParams) ([]response.User, error)
	LeaveTeam(userID int64, teamId int64) error
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

func (t teamRepository) CreateTeam(teamToCreate request.Team, teamMembers []int64) (int64, error) {
	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	var teamId int64
	rows := tx.QueryRow(ctx, `INSERT INTO teams (name, team_privacy, created_by) VALUES ($1, $2, $3) RETURNING id`, teamToCreate.Name, teamToCreate.Privacy, teamToCreate.CreatedBy)
	err = rows.Scan(&teamId)
	if err != nil {
		tx.Rollback(ctx)
		return teamId, err
	}

	batch := &pgx.Batch{}
	for _, v := range teamMembers {
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

	for _, v := range teamMembers {
		t.redisClient.SAdd(ctx, "user:"+strconv.FormatInt(v, 10)+":teams", teamId)
	}

	return teamId, nil
}

func (t teamRepository) AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembersWithTeamID) error {
	var dbTeamCreatedBy int64
	rows := t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToAdd.TeamID)
	err := rows.Scan(&dbTeamCreatedBy)
	if err != nil {
		return err
	}

	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	var args []interface{}
	query := `SELECT privacy FROM users WHERE id IN (`
	for i, v := range teamMembersToAdd.MemberIDs {
		query += `$` + strconv.Itoa(i+1) + `, `
		args = append(args, v)
	}
	if len(teamMembersToAdd.MemberIDs) > 0 {
		query = query[:len(query)-2]
	}
	query += `)`

	users, err := t.dbConn.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}
	defer users.Close()

	var user response.User
	for users.Next() {
		if err := users.Scan(&user.Privacy); err != nil {
			return err
		}
		if user.Privacy == "PRIVATE" {
			return errorhandling.OnlyPublicMemberAllowed
		}
	}

	ctx := context.Background()
	tx, err := t.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, v := range teamMembersToAdd.MemberIDs {
		batch.Queue(`INSERT INTO team_members (team_id, member_id) VALUES ($1, $2)`, teamMembersToAdd.TeamID, v)
	}
	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(ctx)
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == constant.PG_Duplicate_Error_Code {
			return errorhandling.MemberExist
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == constant.PG_Duplicate_Error_Code {
			return errorhandling.MemberExist
		}
		return err
	}

	for _, v := range teamMembersToAdd.MemberIDs {
		t.redisClient.SAdd(ctx, "user:"+strconv.FormatInt(v, 10)+":teams", teamMembersToAdd.TeamID)
	}

	return nil
}

func (t teamRepository) RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembersWithTeamID) error {
	var dbTeamCreatedBy int64
	rows := t.dbConn.QueryRow(context.Background(), `SELECT created_by FROM teams WHERE id = $1`, teamMembersToRemove.TeamID)
	err := rows.Scan(&dbTeamCreatedBy)
	if err != nil {
		return err
	}
	if dbTeamCreatedBy != teamCreatedBy {
		return errorhandling.NotAllowed
	}

	var args []interface{}
	query := `DELETE FROM team_members WHERE member_id IN (`
	for i, v := range teamMembersToRemove.MemberIDs {
		query += `$` + strconv.Itoa(i+1) + `, `
		args = append(args, v)
	}
	if len(teamMembersToRemove.MemberIDs) > 0 {
		query = query[:len(query)-2]
	}
	query += `)`

	_, err = t.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	for _, v := range teamMembersToRemove.MemberIDs {
		t.redisClient.SRem(context.Background(), "user:"+strconv.FormatInt(v, 10)+":teams", teamMembersToRemove.TeamID)
	}

	return nil
}

func (t teamRepository) GetAllTeams(userID int64, queryParams request.TeamQueryParams) ([]response.Team, error) {
	//flag = 0 => created by me, flag = 1 => i am member
	var teams pgx.Rows
	var err error
	teamsSlice := make([]response.Team, 0)

	var query string
	if !queryParams.CreatedByMe {
		query = `SELECT * FROM teams WHERE created_by = $1 AND true`
		query = CreateQueryForParamsOfGetTeam(query, queryParams)
		teams, err = t.dbConn.Query(context.Background(), query, userID)
	}
	if queryParams.CreatedByMe {
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
		if err := teams.Scan(&team.ID, &team.Name, &team.CreatedBy, &team.CreatedAt, &team.TeamPrivacy); err != nil {
			return teamsSlice, err
		}
		teamsSlice = append(teamsSlice, team)
	}

	return teamsSlice, nil
}

func (t teamRepository) GetTeamMembers(teamId int64, queryParams request.TeamQueryParams) ([]response.User, error) {
	var teamMembers pgx.Rows
	var err error
	teamMembersSlice := make([]response.User, 0)

	query := `SELECT id, first_name, last_name, bio, email, privacy FROM users WHERE id IN (SELECT member_id from team_members where team_id = $1)`
	query = CreateQueryForParamsOfGetTeam(query, queryParams)
	teamMembers, err = t.dbConn.Query(context.Background(), query, teamId)

	if err != nil {
		return teamMembersSlice, err
	}
	defer teamMembers.Close()

	var teamMember response.User
	for teamMembers.Next() {
		if err := teamMembers.Scan(&teamMember.ID, &teamMember.FirstName, &teamMember.LastName, &teamMember.Bio, &teamMember.Email, &teamMember.Privacy); err != nil {
			return teamMembersSlice, err
		}
		teamMembersSlice = append(teamMembersSlice, teamMember)
	}

	return teamMembersSlice, nil
}

func CreateQueryForParamsOfGetTeam(query string, queryParams request.TeamQueryParams) string {
	if queryParams.Search != constant.EMPTY_STRING {
		query += fmt.Sprintf(" AND (name ILIKE '%%%s%%')", queryParams.Search)
	}
	if queryParams.SortByCreatedAt {
		query += " ORDER BY created_at"
	}
	query += fmt.Sprintf(" LIMIT %d", queryParams.Limit)
	query += fmt.Sprintf(" OFFSET %d", queryParams.Offset)
	return query
}

func (t teamRepository) LeaveTeam(userID int64, teamId int64) error {
	a, err := t.dbConn.Exec(context.Background(), "DELETE FROM team_members WHERE member_id = $1 AND team_id = $2", userID, teamId)
	if a.RowsAffected() == 0 {
		return errorhandling.NotAMember
	}
	if err != nil {
		return err
	}
	t.redisClient.SRem(context.Background(), "user:"+strconv.FormatInt(userID, 10)+":teams", teamId)
	return nil
}
