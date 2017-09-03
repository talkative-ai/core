package models

import (
	"database/sql"
	"time"

	"github.com/go-gorp/gorp"
)

// User model for the AUM User
type User struct {
	AumModel    `json:"-"`
	GivenName   string
	FamilyName  string
	Email       string
	Image       sql.NullString
	Password    sql.NullString `json:",omitempty" db:"-"`
	PasswordSHA sql.NullString `json:"-"`
	Salt        sql.NullString `json:"-"`
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

func (m *TeamMember) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}
