package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// RunMigrations executes database migrations using golang-migrate for PostgreSQL
func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	// Use golang-migrate for production-ready migration management
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

// RunMigrationsLegacy - fallback simple migration runner for development
func RunMigrationsLegacy(db *sqlx.DB, migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	// Sort files to ensure they run in order, only .up.sql files
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationsPath, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Execute the migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("error executing migration %s: %v", filename, err)
		}
	}

	return nil
}

// GetMigrationVersion returns current migration version
func GetMigrationVersion(db *sql.DB, migrationsPath string) (uint, bool, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %v", err)
	}

	return version, dirty, nil
}