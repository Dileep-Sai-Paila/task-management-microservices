package users

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"-"` // preventing password from being marshalled into JSON response.
	PasswordHash string `json:"-"`
}
