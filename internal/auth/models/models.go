package models

type User struct {
	UID      string
	Name     string
	Email    string
	Password string
}

type UserRequest struct {
	Email    string
	Password string
}
