package db

import (
	"database/sql"

	"github.com/artificial-universe-maker/go-utilities/models"
	"github.com/go-gorp/gorp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Required for sqlx postgres connections
)

// Instance is the PostgreSQL connection instance
var Instance *sqlx.DB
var DBMap *gorp.DbMap

// InitializeDB will setup the DB connection
func InitializeDB() error {
	var err error
	Instance, err = sqlx.Connect("postgres", "user=postgres dbname=postgres host=postgres sslmode=disable")
	if err != nil {
		return err
	}

	DBMap = &gorp.DbMap{Db: Instance.DB, Dialect: gorp.PostgresDialect{}}

	DBMap.AddTableWithName(models.AumProject{}, "workbench_projects")
	DBMap.AddTableWithName(models.AumZone{}, "workbench_zones")
	DBMap.AddTableWithName(models.AumActor{}, "workbench_actors")
	DBMap.AddTableWithName(models.AumZoneActor{}, "workbench_zones_actors")
	DBMap.AddTableWithName(models.AumDialogNode{}, "workbench_dialogs")
	DBMap.AddTableWithName(models.AumNote{}, "workbench_notes")

	DBMap.AddTableWithName(models.User{}, "users")
	DBMap.AddTableWithName(models.Team{}, "teams")
	DBMap.AddTableWithName(models.TeamMember{}, "team_members")

	return nil
}

func CreateAndSaveUser(user *models.User) error {
	err := DBMap.Insert(user)
	if err != nil {
		return err
	}
	team := &models.Team{Name: sql.NullString{Valid: false}}
	if err = DBMap.Insert(team); err != nil {
		return err
	}
	teamMember := &models.TeamMember{UserID: user.ID, TeamID: team.ID, Role: 1}
	if err = DBMap.Insert(teamMember); err != nil {
		return err
	}
	return nil
}
