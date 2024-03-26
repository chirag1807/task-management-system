package request

import "time"

// Task model info
// @Description Task information with title, description, deadline, assignee, status, priority.
type Task struct {
	ID                 int64      `json:"id,omitempty" db:"id" example:"974751326021189496"`
	Title              string     `json:"title" db:"title" example:"GoLang project: Task Manager"`
	Description        string     `json:"description" db:"description" example:"Create Task Manager Project with GoLang as Backend."`
	Deadline           time.Time  `json:"deadline" db:"deadline" example:"2024-03-25T22:59:59.000Z"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty" db:"assignee_individual" example:"974751326021189123"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty" db:"assignee_team" example:"974751326021189234"`
	Status             string     `json:"status" db:"status" example:"TO-DO"`
	Priority           string     `json:"priority" db:"priority" example:"High"`
	CreatedBy          int64      `json:"createdBy" db:"created_by" example:"974751326021189896"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at" example:"2024-03-25T22:59:59.000Z"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty" db:"updated_by" example:"974751326021189896"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" db:"updated_at" example:"2024-03-26T12:49:539.000Z"`
}

// TaskQueryParams model info
// @Description used for retrieving tasks from database with pagination, search, status and sorting option.
type TaskQueryParams struct {
	Limit        int    `json:"limit" example:"10"`
	Offset       int    `json:"offset" example:"0"`
	Search       string `json:"search" example:"GoLang Project"`
	Status       string `json:"status" example:"TO-DO"`
	SortByFilter bool   `json:"sortByFilter" example:"true"`
}
