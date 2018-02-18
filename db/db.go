package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-gorp/gorp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Required for sqlx postgres connections
	"github.com/talkative-ai/core/models"
)

// Instance is the PostgreSQL connection instance
var Instance *sqlx.DB
var DBMap *gorp.DbMap

func GetMaxProjects() int {
	return 3
}

// InitializeDB will setup the DB connection
func InitializeDB() error {
	var err error
	port := os.Getenv("POSTGRES_HOST_PORT")
	if port == "" {
		port = "5432"
	}
	host := os.Getenv("POSTGRES_HOST_ADDR")
	if host == "" {
		host = "postgres"
	}
	Instance, err = sqlx.Connect("postgres", fmt.Sprintf("user=%v password=%v dbname=postgres host=%v port=%v sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		host,
		port))
	if err != nil {
		return err
	}

	DBMap = &gorp.DbMap{Db: Instance.DB, Dialect: gorp.PostgresDialect{}}

	DBMap.AddTableWithName(models.Project{}, "workbench_projects")
	DBMap.AddTableWithName(models.Zone{}, "workbench_zones")
	DBMap.AddTableWithName(models.Actor{}, "workbench_actors")
	DBMap.AddTableWithName(models.ZoneActor{}, "workbench_zones_actors")
	DBMap.AddTableWithName(models.DialogNode{}, "workbench_dialog_nodes")
	DBMap.AddTableWithName(models.DialogRelation{}, "workbench_dialog_nodes_relations")
	DBMap.AddTableWithName(models.PrivateProjectGrants{}, "workbench_private_project_grants")
	DBMap.AddTableWithName(models.VersionedProject{}, "static_published_projects_versioned")
	DBMap.AddTableWithName(models.Note{}, "workbench_notes")

	DBMap.AddTableWithName(models.User{}, "users")
	DBMap.AddTableWithName(models.Team{}, "teams")
	DBMap.AddTableWithName(models.UpgradeItem{}, "upgrade_item")
	DBMap.AddTableWithName(models.TeamMember{}, "team_members")

	DBMap.AddTableWithName(models.EventUserActon{}, "event_user_action")
	DBMap.AddTableWithName(models.EventStateChange{}, "event_state_change")

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
