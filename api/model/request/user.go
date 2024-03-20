package request

type User struct {
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	Bio         string `json:"bio" db:"bio"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	NewPassword string `json:"newPassword" db:"password"`
	Profile     string `json:"profile" db:"profile"`
}

type UserQueryParams struct {
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Search       string `json:"search"`
}
