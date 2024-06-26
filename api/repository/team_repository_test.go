package repository

import (
	"testing"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		TeamDetails  request.Team
		TeamMembers  []int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Created Successfully",
			TeamDetails: request.Team{
				Name:        "Team Stairs",
				Privacy: func() *string { team_privacy := string("PUBLIC"); return &team_privacy }(),
				CreatedBy:   954488202459119617,
				CreatedAt:   time.Now(),
			},
			TeamMembers: []int64{954488202459119617},
			Expected:   nil,
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			_, err := NewTeamRepo(dbConn, redisClient).CreateTeam(v.TeamDetails, v.TeamMembers)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestAddMembersToTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName  string
		TeamCreatedBy int64
		TeamMembers   request.TeamMembersWithTeamID
		Expected      interface{}
		StatusCode    int
	}{
		{
			TestCaseName:  "Team Members Added Successfully",
			TeamCreatedBy: 954488202459119617,
			TeamMembers: request.TeamMembersWithTeamID{
				TeamID:   954507580144451585,
				MemberIDs: []int64{954497896847212545},
			},
			Expected:   nil,
			StatusCode: 200,
		},
		{
			TestCaseName:  "Team Member Already Exist",
			TeamCreatedBy: 954488202459119617,
			TeamMembers: request.TeamMembersWithTeamID{
				TeamID:   954507580144451585,
				MemberIDs: []int64{954497896847212545},
			},
			Expected:   errorhandling.MemberExist,
			StatusCode: 409,
		},
		{
			TestCaseName:  "Not Allowed to Add Member",
			TeamCreatedBy: 954488202459119618,
			TeamMembers: request.TeamMembersWithTeamID{
				TeamID:   954507580144451585,
				MemberIDs: []int64{954497896847212545},
			},
			Expected:   errorhandling.NotAllowed,
			StatusCode: 401,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			err := NewTeamRepo(dbConn, redisClient).AddMembersToTeam(v.TeamCreatedBy, v.TeamMembers)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestRemoveMembersFromTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName  string
		TeamCreatedBy int64
		TeamMembers   request.TeamMembersWithTeamID
		Expected      interface{}
		StatusCode    int
	}{
		{
			TestCaseName:  "Team Members Added Successfully",
			TeamCreatedBy: 954488202459119617,
			TeamMembers: request.TeamMembersWithTeamID{
				TeamID:   954507580144451585,
				MemberIDs: []int64{954497896847212545},
			},
			Expected:   nil,
			StatusCode: 200,
		},
		{
			TestCaseName:  "Not Allowed to Removed Member",
			TeamCreatedBy: 954488202459119618,
			TeamMembers: request.TeamMembersWithTeamID{
				TeamID:   954507580144451585,
				MemberIDs: []int64{954497896847212545},
			},
			Expected:   errorhandling.NotAllowed,
			StatusCode: 401,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			err := NewTeamRepo(dbConn, redisClient).RemoveMembersFromTeam(v.TeamCreatedBy, v.TeamMembers)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestGetAllTeams(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		UserId       int64
		QueryParams  request.TeamQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Created By Me - Success",
			UserId:       954488202459119617,
			QueryParams: request.TeamQueryParams{
				CreatedByMe: false,
				Limit:  1,
				Offset: 0,
				Search: "",
			},
			Expected:   nil,
			StatusCode: 200,
		},
		{
			TestCaseName: "Team in Which I were Added - Success",
			UserId:       954488202459119617,
			QueryParams: request.TeamQueryParams{
				CreatedByMe: true,
				Limit:  1,
				Offset: 0,
				Search: "",
			},
			Expected:   nil,
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			_, err := NewTeamRepo(dbConn, redisClient).GetAllTeams(v.UserId, v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestGetTeamMembers(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		TeamID       int64
		QueryParams  request.TeamQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Created By Me - Success",
			TeamID:       954507580144451585,
			QueryParams: request.TeamQueryParams{
				Limit:  1,
				Offset: 0,
				Search: "",
			},
			Expected:   nil,
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			_, err := NewTeamRepo(dbConn, redisClient).GetTeamMembers(v.TeamID, v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestLeftTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		UserID       int64
		TeamID       int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Left Successfully",
			TeamID:       954507580144451585,
			UserID:       954497896847212545,
			Expected:     errorhandling.NotAMember,
			StatusCode:   401,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			err := NewTeamRepo(dbConn, redisClient).LeaveTeam(v.UserID, v.TeamID)
			assert.Equal(t, v.Expected, err)
		})
	}
}
