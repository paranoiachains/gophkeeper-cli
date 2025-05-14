package models

// User model used by both client and server
type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
