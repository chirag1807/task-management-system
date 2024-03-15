package response

type SuccessResponse struct {
	Message string `json:"message"`
	ID      *int64 `json:"id,omitempty"`
}
