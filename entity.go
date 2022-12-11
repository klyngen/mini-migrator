package migrator

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type MigrationStatus int16

const (
	STARTED  MigrationStatus = 1
	FAILED   MigrationStatus = 2
	COMPLETE MigrationStatus = 3
)

type Migration struct {
	Script      string
	Name        string
	Description string
	order       int
	status      MigrationStatus
	timestamp   time.Time
	hash        string
}

func (m *Migration) createHash() string {
	bytes := md5.Sum([]byte(m.Script))
	return fmt.Sprintf("%x", bytes)
}

type migrator struct {
	driver DatabaseDriver
	db     *sql.DB
}

func NewMigrator(db *sql.DB, driver RelationalDriver) (*migrator, error) {
	driver.DB = db

	return &migrator{
		driver: &driver,
		db:     db,
	}, nil
}

func (m *migrator) MigrateDatabase(migrations []Migration) error {

	exists, err := m.driver.databaseExists()

	if err != nil {
		return errors.New(fmt.Sprintf("Could not see if database exists. I give up: %v\n", err))
	}

	if !exists {
		err = m.driver.createMigrationTable()

		if err != nil {
			return errors.New(fmt.Sprintf("Unable to create a migration table: %v\n", err))
		}
	}

	exsitingMigrations, err := m.driver.fetchMigrations()

	for i, migration := range exsitingMigrations {
		if migration.status == FAILED {
			return errors.New("Unable to do migration since the last migration failed for some reason. Please do some manual work on the database")
		}

		if migration.hash != migrations[i].createHash() {
			return errors.New(fmt.Sprintf("Do not dare to migrate the database. Migration named %s has changed after migration was performed. Migration-script will not do anything. Either stop changing the migration script or update the database with a valid MD5-hash", migrations[i].Name))
		}
	}

	newMigrations := len(migrations) - len(exsitingMigrations)
	if newMigrations > 0 {
		for i := len(exsitingMigrations); i < len(migrations); i++ {
			migrations[i].order = i + 1
			err = m.driver.writeMigration(migrations[i])
			if err != nil {
				log.Fatalf("Unable to write migration: %v", err)
			}

			_, err := m.db.Exec(migrations[i].Script)

			if err != nil {
				m.driver.updateStatus(FAILED, i+1)
				log.Fatalf("Unable to write migration: %v", err)
			}

			err = m.driver.updateStatus(COMPLETE, i+1)

			if err != nil {
				log.Printf("Migration complete but unable to update status: %v", err)
			}
		}
	}
	return nil
}
