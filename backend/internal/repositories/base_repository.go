package repositories

import (
	"fmt"
	"time"

	"flutter-bengkel/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	db *sqlx.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sqlx.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// SoftDelete performs a soft delete on any table
func (br *BaseRepository) SoftDelete(tableName string, id uuid.UUID, deletedBy uuid.UUID) error {
	return utils.SoftDelete(br.db, utils.SoftDeleteOptions{
		TableName:   tableName,
		WhereClause: "id = $3",
		Args:        []interface{}{id},
		DeletedBy:   deletedBy,
	})
}

// Restore restores a soft deleted record
func (br *BaseRepository) Restore(tableName string, id uuid.UUID) error {
	return utils.Restore(br.db, utils.RestoreOptions{
		TableName:   tableName,
		WhereClause: "id = $2",
		Args:        []interface{}{id},
	})
}

// PermanentDelete permanently deletes a record
func (br *BaseRepository) PermanentDelete(tableName string, id uuid.UUID) error {
	return utils.PermanentDelete(br.db, tableName, "id = $1", id)
}

// GenerateUUID generates a new UUID
func (br *BaseRepository) GenerateUUID() uuid.UUID {
	return uuid.New()
}

// SetTimestamps sets created_at and updated_at timestamps for new records
func (br *BaseRepository) SetCreateTimestamps(query string) string {
	return fmt.Sprintf("%s, created_at = NOW(), updated_at = NOW()", query)
}

// SetUpdateTimestamp sets updated_at timestamp for updates
func (br *BaseRepository) SetUpdateTimestamp(query string) string {
	return fmt.Sprintf("%s, updated_at = NOW()", query)
}

// BuildFilterQuery builds a query with soft delete filter
func (br *BaseRepository) BuildFilterQuery(baseQuery string, includeDeleted bool) string {
	if includeDeleted {
		return utils.BuildIncludeDeletedQuery(baseQuery)
	}
	return utils.BuildSoftDeleteQuery(baseQuery)
}