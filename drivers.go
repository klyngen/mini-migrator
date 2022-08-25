package migrator

var SQLiteDriver = RelationalDriver{
	CreationQuery: `
		CREATE TABLE IF NOT EXISTS migrationTable (
		id INTEGER NOT NULL PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		description TEXT,
		name VARCHAR(50),
		hash VARCHAR(36),
		status INTEGER NOT NULL)`,
	ExsistanseQuery:            `SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='migrationTable'`,
	FetchMigrationsQuery:       "SELECT id, name, status, hash FROM migrationTable",
	WriteMigrationQuery:        "INSERT INTO migrationTable  (id, timestamp, description, name, status, hash) VALUES(?, ?, ?, ?, ?, ?)",
	UpdateMigrationStatusQuery: "UPDATE migrationTable SET status = ? WHERE id = ?",
}

var MySqlDriver = RelationalDriver{
	CreationQuery: `
		CREATE TABLE IF NOT EXISTS migrationTable (
		id INTEGER NOT NULL PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		description TEXT,
		name VARCHAR(50),
		hash VARCHAR(36),
		status INTEGER NOT NULL)`,
	ExsistanseQuery:            `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = "migrationTable";`,
	FetchMigrationsQuery:       "SELECT id, name, status, hash FROM migrationTable",
	WriteMigrationQuery:        "INSERT INTO migrationTable  (id, timestamp, description, name, status, hash) VALUES(?, ?, ?, ?, ?, ?)",
	UpdateMigrationStatusQuery: "UPDATE migrationTable SET status = ? WHERE id = ?",
}

var PostgreSQLDriver = RelationalDriver{
	CreationQuery: `
		CREATE TABLE IF NOT EXISTS migrationTable (
		id INTEGER NOT NULL PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		description TEXT,
		name VARCHAR(50),
		hash VARCHAR(36),
		status INTEGER NOT NULL)`,
	ExsistanseQuery: `SELECT
    COUNT(table_name)
FROM
    information_schema.tables
WHERE
    table_schema LIKE 'public' AND
    table_type LIKE 'BASE TABLE' AND
	table_name = 'migrationTable';
`,
	FetchMigrationsQuery:       "SELECT id, name, status, hash FROM migrationTable",
	WriteMigrationQuery:        "INSERT INTO migrationTable  (id, timestamp, description, name, status, hash) VALUES(?, ?, ?, ?, ?, ?)",
	UpdateMigrationStatusQuery: "UPDATE migrationTable SET status = ? WHERE id = ?",
}
