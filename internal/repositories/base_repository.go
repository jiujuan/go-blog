package repositories

import (
	"go-blog/internal/database"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *database.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *database.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the database instance
func (r *BaseRepository) GetDB() *database.DB {
	return r.db
}

// Create creates a new record
func (r *BaseRepository) Create(value interface{}) error {
	return r.db.Create(value)
}

// GetByID retrieves a record by ID
func (r *BaseRepository) GetByID(dest interface{}, id interface{}, preloads ...string) error {
	return r.db.GetByID(dest, id, preloads...)
}

// Update updates a record
func (r *BaseRepository) Update(value interface{}) error {
	return r.db.Update(value)
}

// Delete deletes a record
func (r *BaseRepository) Delete(model interface{}, id interface{}) error {
	return r.db.Delete(model, id)
}

// List retrieves records with pagination
func (r *BaseRepository) List(dest interface{}, options *database.QueryOptions) (*database.PaginationResult, error) {
	return r.db.List(dest, options)
}

// Exists checks if a record exists
func (r *BaseRepository) Exists(model interface{}, conditions ...interface{}) (bool, error) {
	return r.db.Exists(model, conditions...)
}

// Count returns the count of records
func (r *BaseRepository) Count(model interface{}, conditions ...interface{}) (int64, error) {
	return r.db.Count(model, conditions...)
}