package response

// User model info
// @Description User information with id, first name, last name, bio, email, password and privacy.
type User struct {
	ID        int64  `json:"id" example:"974751326021189896"`
	FirstName string `json:"firstName" example:"Chirag"`
	LastName  string `json:"lastName" example:"Makwana"`
	Bio       string `json:"bio" example:"Junior Software Engineer at ZURU TECH INDIA."`
	Email     string `json:"email" example:"chiragmakwana@gmail.com"`
	Password  string `json:"password" example:"Chirag123$,omitempty"`
	Privacy   string `json:"privacy" example:"PUBLIC"`
}

// UserWithTokens model info
// @Description Send user info along with access token and refresh token to response.
type UserWithTokens struct {
	Code         string `json:"code" example:"200 OK"`
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
