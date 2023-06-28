package models

// User represents a user in the application.
type User struct {
	*BaseModel

	Username string
	Email    string
	Password string
}

// NewUser creates a new instance of User model.
func NewUser(username, email, password string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
	}
}
