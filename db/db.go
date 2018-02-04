package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/artificial-universe-maker/core/models"
	"github.com/go-gorp/gorp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Required for sqlx postgres connections
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
	Instance, err = sqlx.Connect("postgres", fmt.Sprintf("user=%v password=%v dbname=postgres host=postgres sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")))
	if err != nil {
		return err
	}

	DBMap = &gorp.DbMap{Db: Instance.DB, Dialect: gorp.PostgresDialect{}}

	DBMap.AddTableWithName(models.AumProject{}, "workbench_projects")
	DBMap.AddTableWithName(models.AumZone{}, "workbench_zones")
	DBMap.AddTableWithName(models.AumActor{}, "workbench_actors")
	DBMap.AddTableWithName(models.AumZoneActor{}, "workbench_zones_actors")
	DBMap.AddTableWithName(models.AumDialogNode{}, "workbench_dialog_nodes")
	DBMap.AddTableWithName(models.AumDialogRelation{}, "workbench_dialog_nodes_relations")
	DBMap.AddTableWithName(models.AumPrivateProjectGrants{}, "workbench_private_project_grants")
	DBMap.AddTableWithName(models.VersionedProjectJSONSafe{}, "static_published_projects_versioned")
	DBMap.AddTableWithName(models.AumNote{}, "workbench_notes")

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
