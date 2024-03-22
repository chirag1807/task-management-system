package response

import (
	"time"
)

type Team struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	TeamProfile string       `json:"teamProfile"`
	CreatedBy   int64        `json:"createdBy"`
	CreatedAt   time.Time    `json:"createdAt"`
	TeamMembers TeamMembers `json:"teamMembers,omitempty"`
}

type TeamMembers struct {
	TeamID   int64   `json:"teamID"`
	MemberID []int64 `json:"memberID"`
}

type TeamMemberDetails struct {
	TeamMembers []User `json:"teamMembers"`
}

type Teams struct {
	Teams []Team `json:"team"`
}
