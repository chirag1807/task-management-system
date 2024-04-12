package response

// SuccessResponse model info
// @Description Send success response to client with corresponding message and id(if any).
type SuccessResponse struct {
	Code    string `json:"code" example:"200 OK"`
	Message string `json:"message,omitempty" example:"Task Created Successfully."`
	ID      *int64 `json:"id,omitempty" example:"974751326021189896"`
}
