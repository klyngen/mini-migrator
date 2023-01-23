# Basic golang database migrator

Simple idea. Just put the SQL migrations into your code.

### Why is this better than exsiting solutions

- No need to ship migration-files with your executable
- No CLI tools
- No extra pipelines executing migrations
- Easy to extend

## Usage

```golang
package main

import (
  "github.com/klyngen/mini-migrator",
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
)

migrations := []migrator.Migration{{
    Name:        "test1",
    Description: "must see that this tooling works",
    Script:      "CREATE TABLE TEST1 (id INTEGER, name VARHCAR(50))",
}}

func main() {
	db, _ := sql.Open("sqlite3", fmt.Sprintf("file:%v", dbName))
	m, _ := migrator.NewMigrator(db, migrator.SQLiteDriver, migrator.MigrationOptions{Strict: true})

	err := m.MigrateDatabase(migrations)

    if err != nil {
      // HANDLE this properly
    }
}

```

## Features

- 📦 Supports most common databases (postgres, mysql, sqlite)
- 📦 Easy to extend for most common relational databases

## How it works

Every migration is put into a table with a name, description, status and a hash. The status is to keep track of the executed migrations and the hash is to ensure that you don't change the queries retrospectively.

1. Validation of existing migrations. If earlier migrations were successful and hashes compare, we proceed.
2. Write the migration with status `PROGRESS`
3. Execute the migration
4. Update the migration status to either `FAILED` or `COMPLETE`
