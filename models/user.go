package models

import (
	"database/sql"
)

// User model for the AUM User
type User struct {
	ID          uint64 `db:"id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Image       string `db:"image"`
	Password    string `db:"-"`
	Passwordsha string `json:"-"`
	Salt        string `json:"-"`
}

// Team model relates multiple users under the same umbrella
// If the Name is null, then it's the user by themselves
type Team struct {
	ID   uint64         `db:"id"`
	Name sql.NullString `db:"name"`
}

// TeamMember is the relationship bretween a user and a team
// and includes the user's role within the team
type TeamMember struct {
	UserID uint64 `db:"user_id"`
	TeamID uint64 `db:"team_id"`
	Role   int    `db:"role"`
}
