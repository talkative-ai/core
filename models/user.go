package models

// User model for the AUM User
type User struct {
	ID    uint64
	Email string
}

// UserLinkedAccount model for linking accounts to a User (e.g. Google)
type UserLinkedAccount struct {
	UserID   uint64
	Email    string
	Provider string
}

type Team struct {
	ID   uint64
	Name string
}
