package response

import (
	"time"
)

// Team model info
// @Description Team information with it's id, name, profile (Public or Private), id of user who created it, time when it was created and team members.
type Team struct {
	ID          int64       `json:"id" example:"954751326021189633"`
	Name        string      `json:"name" example:"Team Jupiter"`
	TeamProfile string     `json:"teamProfile" example:"Public"`
	CreatedBy   int64       `json:"createdBy" example:"954751326021189799"`
	CreatedAt   time.Time   `json:"createdAt" example:"2024-03-25T22:59:59.000Z"`
	TeamMembers TeamMembers `json:"teamMembers,omitempty"`
}

// TeamMembers model info
// @Description Send team's id and it's all members id to the response.
type TeamMembers struct {
	TeamID   int64   `json:"teamID" example:"954751326021189633"`
	MemberID []int64 `json:"memberID" example:"954751326021189800,954751326021189801"`
}

// TeamMemberDetails model info
// @Description Send array of user to response as team members.
type TeamMemberDetails struct {
	TeamMembers []User `json:"teamMembers"`
}

// Teams model info
// @Description Send array of team to response.
type Teams struct {
	Teams []Team `json:"team"`
}
