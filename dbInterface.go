package migrator

import "database/sql"

type DatabaseDriver interface {
	databaseExists() (bool, error)
	createMigrationTable() error
	fetchMigrations() ([]Migration, error)
	writeMigration(Migration) error
	updateStatus(status MigrationStatus, id int) error
}

type Driver struct {
	db *sql.DB
}
