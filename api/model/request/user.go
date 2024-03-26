package request

// User model info
// @Description User information with first name, last name, bio, email, password and profile.
type User struct {
	FirstName   string `json:"firstName" db:"first_name" example:"Chirag"`
	LastName    string `json:"lastName" db:"last_name" example:"Makwana"`
	Bio         string `json:"bio" db:"bio" example:"Junior Software Engineer at ZURU TECH INDIA."`
	Email       string `json:"email" db:"email" example:"chiragmakwana@gmail.com"`
	Password    string `json:"password" db:"password" example:"Chirag123$"`
	NewPassword string `json:"newPassword" db:"password" example:"Chirag123$"`
	Profile     string `json:"profile" db:"profile" example:"Public"`
}

// UserQueryParams model info
// @Description used for retrieving users from database with pagination and search option.
type UserQueryParams struct {
	Limit  int    `json:"limit" example:"10"`
	Offset int    `json:"offset" example:"0"`
	Search string `json:"search" example:"Chirag"`
}
