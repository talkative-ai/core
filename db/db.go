package db

import (
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

	DBMap.AddTableWithName(models.AumProject{}, "projects")
	DBMap.AddTableWithName(models.AumZone{}, "zones")
	DBMap.AddTableWithName(models.AumActor{}, "actors")
	DBMap.AddTableWithName(models.AumDialogNode{}, "dialogs")
	DBMap.AddTableWithName(models.AumNote{}, "notes")

	DBMap.AddTableWithName(models.User{}, "users")
	DBMap.AddTableWithName(models.UserLinkedAccount{}, "user_linked_accounts")
	DBMap.AddTableWithName(models.Team{}, "teams")

	return nil
}
