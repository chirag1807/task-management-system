package request

import (
	"time"
)

// Team model info
// @Description Team information with it's id, name, profile (Public or Private), id of user who created it and time when it was created.
type Team struct {
	ID          int64     `json:"id,omitempty" example:"954751326021189633"`
	Name        string    `json:"name" example:"Team Jupiter"`
	TeamProfile *string   `json:"teamProfile,omitempty" example:"Public"`
	CreatedBy   int64     `json:"createdBy" example:"954751326021189799"`
	CreatedAt   time.Time `json:"createdAt,omitempty" example:"2024-03-25T22:59:59.000Z"`
}

// TeamMembers model info
// @Description Team's id and it's all members id.
type TeamMembers struct {
	TeamID   int64   `json:"teamID,omitempty" example:"954751326021189633"`
	MemberID []int64 `json:"memberID" example:"954751326021189800,954751326021189801"`
}

// CreateTeam model info
// @Description Team Details such as name, profile along with members id.
type CreateTeam struct {
	TeamDetails Team        `json:"teamDetails"`
	TeamMembers TeamMembers `json:"teamMembers"`
}

// TeamQueryParams model info
// @Description used for retrieving teams from database with pagination, search and sorting(based on date created) option.
type TeamQueryParams struct {
	Limit          int    `json:"limit" example:"10"`
	Offset         int    `json:"offset" example:"0"`
	Search         string `json:"search" example:"Jupiter"`
	SortByCreateAt bool   `json:"sortByCreatedAt" example:"true"`
}
