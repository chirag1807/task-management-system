package response

// User model info
// @Description User information with id, first name, last name, bio, email, password and profile.
type User struct {
	ID        int64  `json:"id" example:"974751326021189896"`
	FirstName string `json:"firstName" example:"Chirag"`
	LastName  string `json:"lastName" example:"Makwana"`
	Bio       string `json:"bio" example:"Junior Software Engineer at ZURU TECH INDIA."`
	Email     string `json:"email" example:"chiragmakwana@gmail.com"`
	Password  string `json:"password" example:"Chirag123$"`
	Profile   string `json:"profile" example:"Public"`
}

// UserWithTokens model info
// @Description Send user info along with access token and refresh token to response.
type UserWithTokens struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

// Users model info
// @Description Send array of user to response.
type Users struct {
	Users []User `json:"users"`
}
