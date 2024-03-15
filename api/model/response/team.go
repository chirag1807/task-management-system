package response

import (
	"time"
)

type Team struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	CreatedBy   int64       `json:"createdBy"`
	CreatedAt   time.Time   `json:"createdAt"`
	TeamMembers TeamMembers `json:"teamMembers,omitempty"`
}

type TeamMembers struct {
	TeamID   int64   `json:"teamID"`
	MemberID []int64 `json:"memberID"`
}
