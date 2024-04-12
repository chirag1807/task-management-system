package request

import "time"

// Task model info
// @Description Task information with title, description, deadline, assignee, status, priority.
type Task struct {
	ID                 int64      `json:"id,omitempty" db:"id" example:"974751326021189496" validate:"number"`
	Title              string     `json:"title" db:"title" example:"GoLang project: Task Manager" validate:"required,alphanum_with_spaces,min=4,max=48"`
	Description        string     `json:"description" db:"description" example:"Create Task Manager Project with GoLang as Backend." validate:"required,alphanum_with_spaces,min=12,max=196"`
	Deadline           time.Time  `json:"deadline" db:"deadline" example:"2024-03-25T22:59:59.000Z" validate:"required,time"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty" db:"assignee_individual" example:"974751326021189123" validate:"omitempty,number"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty" db:"assignee_team" example:"974751326021189234" validate:"omitempty,number"`
	Status             string     `json:"status" db:"status" example:"TO-DO" validate:"required,oneof=TO-DO In-PROGRESS COMPLETED CLOSED"`
	Priority           string     `json:"priority" db:"priority" example:"High" validate:"required,oneof=LOW MEDIUM HIGH 'VERY HIGH'"`
	CreatedBy          int64      `json:"createdBy" db:"created_by" example:"974751326021189896"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at" example:"2024-03-25T22:59:59.000Z"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty" db:"updated_by" example:"974751326021189896"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" db:"updated_at" example:"2024-03-26T12:49:539.000Z"`
}

type UpdateTask struct {
	ID                 int64      `json:"id,omitempty" db:"id" example:"974751326021189496" validate:"required,number"`
	Title              string     `json:"title" db:"title" example:"GoLang project: Task Manager" validate:"omitempty,alphanum_with_spaces,min=4,max=48"`
	Description        string     `json:"description" db:"description" example:"Create Task Manager Project with GoLang as Backend." validate:"omitempty,alphanum_with_spaces,min=12,max=196"`
	Deadline           time.Time  `json:"deadline" db:"deadline" example:"2024-03-25T22:59:59.000Z" validate:"omitempty,time"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty" db:"assignee_individual" example:"974751326021189123" validate:"omitempty,number"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty" db:"assignee_team" example:"974751326021189234" validate:"omitempty,number"`
	Status             string     `json:"status" db:"status" example:"TO-DO" validate:"omitempty,oneof=TO-DO In-PROGRESS COMPLETED CLOSED"`
	Priority           string     `json:"priority" db:"priority" example:"High" validate:"omitempty,oneof=LOW MEDIUM HIGH 'VERY HIGH'"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty" db:"updated_by" example:"974751326021189896"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" db:"updated_at" example:"2024-03-26T12:49:539.000Z"`
}

// TaskQueryParams model info
// @Description used for retrieving tasks from database with pagination, search, status and sorting option.
type TaskQueryParams struct {
	CreatedByMe  bool   `json:"createdByMe" example:"true" validate:"boolean"`
	Limit        int    `json:"limit" example:"10" validate:"number,gte=0,max=50"`
	Offset       int    `json:"offset" example:"0" validate:"number"`
	Search       string `json:"search" example:"GoLang Project" validate:"omitempty,alphanum_with_spaces"`
	Status       string `json:"status" example:"TO-DO" validate:"omitempty,oneof=TO-DO In-PROGRESS COMPLETED CLOSED"`
	SortByFilter bool   `json:"sortByFilter" example:"true" validate:"boolean"`
}
