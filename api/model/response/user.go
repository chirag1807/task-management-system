package response

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Profile   string `json:"profile"`
}

type UserWithTokens struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
