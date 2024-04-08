package request

// User model info
// @Description User information with first name, last name, bio, email, password and profile.
type User struct {
	FirstName   string `json:"firstName" db:"first_name" example:"Chirag" validate:"required,alpha_with_spaces,min=2"`
	LastName    string `json:"lastName" db:"last_name" example:"Makwana" validate:"required,alpha_with_spaces,min=2"`
	Bio         string `json:"bio" db:"bio" example:"Junior Software Engineer at ZURU TECH INDIA." validate:"required,alphanum_with_spaces,min=6,max=96"`
	Email       string `json:"email" db:"email" example:"chiragmakwana@gmail.com" validate:"required,email"`
	Password    string `json:"password" db:"password" example:"Chirag123$" validate:"required,min=8"`
	NewPassword string `json:"newPassword" db:"password" example:"Chirag123$" validate:"omitempty,min=8"`
	Profile     string `json:"profile" db:"profile" example:"Public" validate:"required,alpha,oneof=Public Private"`
}

// UserCredentials model info
// @Description User credentials with email and password.
type UserCredentials struct {
	Email    string `json:"email" db:"email" example:"chiragmakwana@gmail.com" validate:"required,email"`
	Password string `json:"password" db:"password" example:"Chirag123$" validate:"required,min=8"`
}

type UserEmail struct {
	Email string `json:"email" db:"email" example:"chiragmakwana@gmail.com" validate:"required,email"`
}

// UpdateUser model info
// @Description User information with first name, last name, bio, email, password and profile.
type UpdateUser struct {
	FirstName   string `json:"firstName" db:"first_name" example:"Chirag" validate:"omitempty,alpha_with_spaces,min=2"`
	LastName    string `json:"lastName" db:"last_name" example:"Makwana" validate:"omitempty,alpha_with_spaces,min=2"`
	Bio         string `json:"bio" db:"bio" example:"Junior Software Engineer at ZURU TECH INDIA." validate:"omitempty,alphanum_with_spaces,min=6,max=96"`
	Email       string `json:"email" db:"email" example:"chiragmakwana@gmail.com" validate:"omitempty,email"`
	Password    string `json:"password" db:"password" example:"Chirag123$" validate:"omitempty,min=8"`
	NewPassword string `json:"newPassword" db:"password" example:"Chirag123$" validate:"omitempty,min=8"`
	Profile     string `json:"profile" db:"profile" example:"Public" validate:"omitempty,alpha,oneof=Public Private"`
}

// UserQueryParams model info
// @Description used for retrieving users from database with pagination and search option.
type UserQueryParams struct {
	Limit  int    `json:"limit" example:"10" validate:"number"`
	Offset int    `json:"offset" example:"0" validate:"number"`
	Search string `json:"search" example:"Chirag" validate:"omitempty,alphanum_with_spaces"`
}
