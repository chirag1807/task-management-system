package response

import "time"

type Task struct {
	ID                 int64      `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Deadline           time.Time  `json:"deadline"`
	AssigneeIndividual *int64     `json:"assigneeIndividual,omitempty"`
	AssigneeTeam       *int64     `json:"assigneeTeam,omitempty"`
	Status             string     `json:"status"`
	Priority           string     `json:"priority"`
	CreatedBy          int64      `json:"createdBy"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedBy          *int64     `json:"updatedBy,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
}
