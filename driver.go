package migrator

import (
	"database/sql"
	"time"
)

type RelationalDriver struct {
	/* Should create a table containing id INT, timestampt DATETIME, description VARCHAR(300), name VARCHAR(30), status INT
	 ```sql
	CREATE TABLE IF NOT EXISTS migrationTable (
	id INTEGER NOT NULL PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	description TEXT,
	/name VARCHAR(50),
	status INTEGER NOT NULL)
	 ```*/
	CreationQuery string
	// ExsistanseQuery should select a single number. 1 if the table exists 0 if not
	// `SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name=?`
	ExsistanseQuery string
	// FetchMigrationsQuery should get status, id and name fro mthe database
	// `SELECT id, name, status, hash FROM migrationTable`
	FetchMigrationsQuery string
	// WriteMigrationQuery inserts 5 values. id, timestamp, description, name, status
	// `INSERT INTO migrationTable  (id, timestamp, description, name, status) VALUES(?, ?, ?, ?, ?)`
	WriteMigrationQuery string
	// UpdateMigrationStatusQuery set the status of a row based on the ID / index. Takes in id int and status int
	// `UPDATE migrationTable SET status = ? WHERE id = ?`
	UpdateMigrationStatusQuery string
	DB                         *sql.DB
}

func (r *RelationalDriver) databaseExists() (bool, error) {
	var count int
	row := r.DB.QueryRow(r.ExsistanseQuery)
	if row.Err() != nil {
		return false, row.Err()
	}

	err := row.Scan(&count)
	return count == 1, err
}

func (r *RelationalDriver) createMigrationTable() error {
	_, err := r.DB.Exec(r.CreationQuery)
	return err
}

func (r *RelationalDriver) fetchMigrations() ([]Migration, error) {
	rows, err := r.DB.Query(r.FetchMigrationsQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	resultSet := make([]Migration, 0)

	var id int
	var name string
	var status int
	var hash string

	for rows.Next() {
		if err = rows.Scan(&id, &name, &status, &hash); err != nil {
			continue
		}
		resultSet = append(resultSet, Migration{
			Name:   name,
			status: MigrationStatus(status),
			order:  id,
			hash:   hash,
		})

	}

	return resultSet, nil
}

func (r *RelationalDriver) writeMigration(migration Migration) error {
	_, err := r.DB.Exec(r.WriteMigrationQuery, migration.order, time.Now(), migration.Description, migration.Name, STARTED, migration.createHash())
	return err
}

func (r *RelationalDriver) updateStatus(status MigrationStatus, id int) error {
	_, err := r.DB.Exec(r.UpdateMigrationStatusQuery, status, id)
	return err
}
