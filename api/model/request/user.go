package request

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Profile   string `json:"profile"`
}
