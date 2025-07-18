package database

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	// Sort files to ensure they run in order
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
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
			return err
		}
	}

	return nil
}