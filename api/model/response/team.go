package response

import (
	"time"
)

// Team model info
// @Description Team information with it's id, name, privacy (PUBLIC or PRIVATE), id of user who created it, time when it was created and team members.
type Team struct {
	ID          int64       `json:"id" example:"954751326021189633"`
	Name        string      `json:"name" example:"Team Jupiter"`
	TeamPrivacy string     `json:"teamPrivacy" example:"PUBLIC"`
	CreatedBy   int64       `json:"createdBy" example:"954751326021189799"`
	CreatedAt   time.Time   `json:"createdAt" example:"2024-03-25T22:59:59.000Z"`
}

// TeamMembers model info
// @Description Send team's id and it's all members id to the response.
type TeamMembers struct {
	TeamID   int64   `json:"teamId" example:"954751326021189633"`
	MemberIDs []int64 `json:"memberIds" example:"954751326021189800,954751326021189801"`
}

