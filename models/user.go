package models

type User struct {
	Id             int
	Username       string
	Email          string
	Password_hash  string
	Created_at     string
	Updated_at     string
	Email_verified bool
}
