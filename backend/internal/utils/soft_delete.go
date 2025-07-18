package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SoftDeleteOptions defines options for soft delete operations
type SoftDeleteOptions struct {
	TableName   string
	WhereClause string
	Args        []interface{}
	DeletedBy   uuid.UUID
}

// RestoreOptions defines options for restore operations
type RestoreOptions struct {
	TableName   string
	WhereClause string
	Args        []interface{}
}

// SoftDelete performs a soft delete operation on a table
func SoftDelete(db *sqlx.DB, opts SoftDeleteOptions) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET deleted_at = $1, deleted_by = $2, updated_at = $1
		WHERE %s AND deleted_at IS NULL
	`, opts.TableName, opts.WhereClause)

	args := []interface{}{time.Now(), opts.DeletedBy}
	args = append(args, opts.Args...)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to soft delete from %s: %w", opts.TableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows found to delete in %s", opts.TableName)
	}

	return nil
}

// Restore performs a restore operation on soft-deleted records
func Restore(db *sqlx.DB, opts RestoreOptions) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET deleted_at = NULL, deleted_by = NULL, updated_at = $1
		WHERE %s AND deleted_at IS NOT NULL
	`, opts.TableName, opts.WhereClause)

	args := []interface{}{time.Now()}
	args = append(args, opts.Args...)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore from %s: %w", opts.TableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no deleted rows found to restore in %s", opts.TableName)
	}

	return nil
}

// PermanentDelete performs a permanent delete operation
func PermanentDelete(db *sqlx.DB, tableName, whereClause string, args ...interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to permanently delete from %s: %w", tableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows found to delete in %s", tableName)
	}

	return nil
}

// BuildSoftDeleteQuery builds a query that excludes soft-deleted records
func BuildSoftDeleteQuery(baseQuery string) string {
	return fmt.Sprintf("%s AND deleted_at IS NULL", baseQuery)
}

// BuildIncludeDeletedQuery builds a query that includes soft-deleted records
func BuildIncludeDeletedQuery(baseQuery string) string {
	return baseQuery // No modification needed
}

// BuildOnlyDeletedQuery builds a query that only returns soft-deleted records
func BuildOnlyDeletedQuery(baseQuery string) string {
	return fmt.Sprintf("%s AND deleted_at IS NOT NULL", baseQuery)
}

// QueryFilter represents different query filter types
type QueryFilter int

const (
	FilterActiveOnly QueryFilter = iota
	FilterIncludeDeleted
	FilterDeletedOnly
)

// ApplyQueryFilter applies the appropriate filter to a base query
func ApplyQueryFilter(baseQuery string, filter QueryFilter) string {
	switch filter {
	case FilterActiveOnly:
		return BuildSoftDeleteQuery(baseQuery)
	case FilterIncludeDeleted:
		return BuildIncludeDeletedQuery(baseQuery)
	case FilterDeletedOnly:
		return BuildOnlyDeletedQuery(baseQuery)
	default:
		return BuildSoftDeleteQuery(baseQuery) // Default to active only
	}
}