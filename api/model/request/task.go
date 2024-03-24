package request

import "time"

type Task struct {
	ID                 int64      `json:"id,omitempty" db:"id"`
	Title              string     `json:"title" db:"title"`
	Description        string     `json:"description" db:"description"`
	Deadline           time.Time  `json:"deadline" db:"deadline"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty" db:"assignee_individual"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty" db:"assignee_team"`
	Status             string     `json:"status" db:"status"`
	Priority           string     `json:"priority" db:"priority"`
	CreatedBy          int64      `json:"createdBy" db:"created_by"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty" db:"updated_by"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

type TaskQueryParams struct {
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Search       string `json:"search"`
	Status       string `json:"status"`
	SortByFilter bool   `json:"sortByFilter"`
}
