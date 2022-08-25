package migrator_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/klyngen/mini-migrator"
	_ "github.com/mattn/go-sqlite3"
)

var dbName = "testdb.db"

func getDatabaseContext() *sql.DB {
	db, _ := sql.Open("sqlite3", fmt.Sprintf("file:%v", dbName))
	return db
}

func removeDatabase() {
	os.Remove(dbName)
}

func TestMigrations_MigrateDatabase(t *testing.T) {
	db := getDatabaseContext()
	defer removeDatabase()

	migrations := []migrator.Migration{{
		Name:        "test1",
		Description: "must see that this tooling works",
		Script:      "CREATE TABLE TEST1 (id INTEGER, name VARHCAR(50))",
	}}

	m, err := migrator.NewMigrator(db, migrator.SQLiteDriver)

	if err != nil {
		t.Log("Could create migrator")
		t.Fail()
	}

	m.MigrateDatabase(migrations)
}

func TestMigrations_EnsureMigrationFailsWhenHashChanges(t *testing.T) {
	db := getDatabaseContext()
	defer removeDatabase()

	migrations := []migrator.Migration{{
		Name:        "test1",
		Description: "must see that this tooling works",
		Script:      "CREATE TABLE TEST1 (id INTEGER, name VARHCAR(50))",
	}}

	m, _ := migrator.NewMigrator(db, migrator.SQLiteDriver)

	err := m.MigrateDatabase(migrations)

	if err != nil {
		t.Log("Migration should not return an error")
		t.Fail()
	}

	migrations[0].Script = "CREATE TABLE TEST2 (id INTEGER, name VARHCAR(50))"

	err = m.MigrateDatabase(migrations)

	if err == nil {
		t.Log("Migration should fail")
		t.Fail()
	}

}
