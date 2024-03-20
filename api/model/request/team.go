package request

import (
	"time"
)

type Team struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name"`
	TeamProfile *string   `json:"teamProfile,omitempty"`
	CreatedBy   int64     `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type TeamMembers struct {
	TeamID   int64   `json:"teamID,omitempty"`
	MemberID []int64 `json:"memberID"`
}

type CreateTeam struct {
	TeamDetails Team        `json:"teamDetails"`
	TeamMembers TeamMembers `json:"teamMembers"`
}

type TeamQueryParams struct {
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Search         string `json:"search"`
	SortByCreateAt bool   `json:"sortByCreatedAt"`
}
