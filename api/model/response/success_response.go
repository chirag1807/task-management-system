package response

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	ID      *int64      `json:"id,omitempty"`
}
