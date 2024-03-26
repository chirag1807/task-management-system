package response

import "time"

// Task model info
// @Description Task information with title, description, deadline, assignee, status, priority.
type Task struct {
	ID                 int64      `json:"id,omitempty" example:"974751326021189496"`
	Title              string     `json:"title" example:"GoLang project: Task Manager"`
	Description        string     `json:"description" example:"Create Task Manager Project with GoLang as Backend."`
	Deadline           time.Time  `json:"deadline" example:"2024-03-25T22:59:59.000Z"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty" example:"974751326021189123"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty" example:"974751326021189234"`
	Status             string     `json:"status" example:"TO-DO"`
	Priority           string     `json:"priority" example:"High"`
	CreatedBy          int64      `json:"createdBy" example:"974751326021189896"`
	CreatedAt          time.Time  `json:"createdAt" example:"2024-03-25T22:59:59.000Z"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty" example:"974751326021189896"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" example:"2024-03-26T12:49:539.000Z"`
}

// Task model info
// @Description Send array of tasks to response.
type Tasks struct {
	Tasks []Task `json:"tasks"`
}
