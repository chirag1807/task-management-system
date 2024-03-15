package repository

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type TeamRepository interface{
	CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (error)
	AddMembersToTeam(teamMembersToAdd request.TeamMembers) (error)
	RemoveMembersFromTeam(teamMembersToRemove request.TeamMembers) (error)
	GetAllTeams(userID int64, flag string) ([]response.Team, error)
	//flag is used for get my created teams and get teams in which i was added.
	GetTeamMembers(teamID int64) (response.TeamMembers, error)
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

func (t teamRepository) CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (error) {
	//see if teamMembers could be fit in team struct.
	return nil
}

func (t teamRepository) AddMembersToTeam(teamMembersToAdd request.TeamMembers) (error) {
	return nil
}

func (t teamRepository) RemoveMembersFromTeam(teamMembersToRemove request.TeamMembers) (error) {
	return nil
}

func (t teamRepository) GetAllTeams(userID int64, flag string) ([]response.Team, error) {
	return []response.Team{}, nil
}

func (t teamRepository) GetTeamMembers(teamID int64) (response.TeamMembers, error) {
	return response.TeamMembers{}, nil
}