package entity

type User struct {
	ID          uint
	PhoneNumber string
	Name        string
	// Password always keep hashed password.
	Password string
	Role     Role
}
