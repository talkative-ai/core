package models

import (
	"database/sql"

	"github.com/go-gorp/gorp"
)

// User model for the AUM User
type User struct {
	AumModel
	GivenName   string
	FamilyName  string
	Email       string
	Image       string
	Password    string `db:"-"`
	PasswordSHA string `json:"-"`
	Salt        string `json:"-"`
}

// Team model relates multiple users under the same umbrella
// If the Name is null, then it's the user by themselves
type Team struct {
	AumModel
	Name sql.NullString
}

// TeamMember is the relationship bretween a user and a team
// and includes the user's role within the team
type TeamMember struct {
	UserID    uint64
	TeamID    uint64
	Role      int
	CreatedAt gorp.NullTime `json:"CreatedAt,omitempty"`
}
