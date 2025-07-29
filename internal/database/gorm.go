package database

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// DB wraps gorm.DB with additional helper methods
type DB struct {
	*gorm.DB
}

// NewDB creates a new DB wrapper
func NewDB(db *gorm.DB) *DB {
	return &DB{DB: db}
}

// PaginationResult represents paginated query result
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// QueryOptions represents common query options
type QueryOptions struct {
	Page     int                    `json:"page"`
	Limit    int                    `json:"limit"`
	OrderBy  string                 `json:"order_by"`
	Filters  map[string]interface{} `json:"filters"`
	Preloads []string               `json:"preloads"`
	Search   *SearchOptions         `json:"search"`
}

// SearchOptions represents search configuration
type SearchOptions struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}

// DefaultQueryOptions returns default query options
func DefaultQueryOptions() *QueryOptions {
	return &QueryOptions{
		Page:    1,
		Limit:   10,
		OrderBy: "created_at DESC",
		Filters: make(map[string]interface{}),
	}
}

// Create creates a new record
func (db *DB) Create(value interface{}) error {
	return db.DB.Create(value).Error
}

// GetByID retrieves a record by ID with optional preloads
func (db *DB) GetByID(dest interface{}, id interface{}, preloads ...string) error {
	query := db.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(dest, id).Error
}

// GetByField retrieves a record by a specific field
func (db *DB) GetByField(dest interface{}, field string, value interface{}, preloads ...string) error {
	query := db.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.Where(field+" = ?", value).First(dest).Error
}

// Update updates a record
func (db *DB) Update(value interface{}) error {
	return db.DB.Save(value).Error
}

// UpdateFields updates specific fields of a record
func (db *DB) UpdateFields(model interface{}, id interface{}, fields map[string]interface{}) error {
	return db.DB.Model(model).Where("id = ?", id).Updates(fields).Error
}

// Delete soft deletes a record by ID
func (db *DB) Delete(model interface{}, id interface{}) error {
	return db.DB.Delete(model, id).Error
}

// HardDelete permanently deletes a record
func (db *DB) HardDelete(model interface{}, id interface{}) error {
	return db.DB.Unscoped().Delete(model, id).Error
}

// Exists checks if a record exists
func (db *DB) Exists(model interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	query := db.DB.Model(model)
	
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// Count returns the count of records matching conditions
func (db *DB) Count(model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := db.DB.Model(model)
	
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// List retrieves records with pagination and filtering
func (db *DB) List(dest interface{}, options *QueryOptions) (*PaginationResult, error) {
	if options == nil {
		options = DefaultQueryOptions()
	}

	// Validate pagination parameters
	if options.Page < 1 {
		options.Page = 1
	}
	if options.Limit < 1 {
		options.Limit = 10
	}
	if options.Limit > 100 {
		options.Limit = 100 // Prevent excessive queries
	}

	// Get the model type for counting
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return nil, errors.New("dest must be a pointer to slice")
	}

	sliceType := destValue.Elem().Type()
	elementType := sliceType.Elem()
	
	// Create a new instance for counting
	modelInstance := reflect.New(elementType).Interface()

	// Build query
	query := db.DB.Model(modelInstance)

	// Apply filters
	for field, value := range options.Filters {
		query = query.Where(field+" = ?", value)
	}

	// Apply search if provided
	if options.Search != nil && options.Search.Query != "" {
		if len(options.Search.Fields) > 0 {
			searchQuery := ""
			searchArgs := make([]interface{}, 0)
			
			for i, field := range options.Search.Fields {
				if i > 0 {
					searchQuery += " OR "
				}
				searchQuery += field + " LIKE ?"
				searchArgs = append(searchArgs, "%"+options.Search.Query+"%")
			}
			
			query = query.Where(searchQuery, searchArgs...)
		}
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply preloads
	for _, preload := range options.Preloads {
		query = query.Preload(preload)
	}

	// Apply ordering
	if options.OrderBy != "" {
		query = query.Order(options.OrderBy)
	}

	// Apply pagination
	offset := (options.Page - 1) * options.Limit
	query = query.Offset(offset).Limit(options.Limit)

	// Execute query
	if err := query.Find(dest).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((total + int64(options.Limit) - 1) / int64(options.Limit))

	return &PaginationResult{
		Data:       dest,
		Total:      total,
		Page:       options.Page,
		Limit:      options.Limit,
		TotalPages: totalPages,
	}, nil
}

// FindWithConditions finds records with complex conditions
func (db *DB) FindWithConditions(dest interface{}, conditions map[string]interface{}, preloads ...string) error {
	query := db.DB
	
	// Apply preloads
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	// Apply conditions
	for field, value := range conditions {
		query = query.Where(field+" = ?", value)
	}
	
	return query.Find(dest).Error
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(fn func(*DB) error) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return fn(NewDB(tx))
	})
}

// BulkCreate creates multiple records in a single query
func (db *DB) BulkCreate(values interface{}, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100
	}
	return db.DB.CreateInBatches(values, batchSize).Error
}

// BulkUpdate updates multiple records
func (db *DB) BulkUpdate(model interface{}, updates map[string]interface{}, conditions ...interface{}) error {
	query := db.DB.Model(model)
	
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	return query.Updates(updates).Error
}

// BulkDelete deletes multiple records
func (db *DB) BulkDelete(model interface{}, conditions ...interface{}) error {
	query := db.DB
	
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	
	return query.Delete(model).Error
}

// Increment increments a numeric field
func (db *DB) Increment(model interface{}, id interface{}, field string, value interface{}) error {
	return db.DB.Model(model).Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" + ?", value)).Error
}

// Decrement decrements a numeric field
func (db *DB) Decrement(model interface{}, id interface{}, field string, value interface{}) error {
	return db.DB.Model(model).Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" - ?", value)).Error
}

// Raw executes a raw SQL query
func (db *DB) Raw(sql string, values ...interface{}) *gorm.DB {
	return db.DB.Raw(sql, values...)
}

// Exec executes a raw SQL command
func (db *DB) Exec(sql string, values ...interface{}) error {
	return db.DB.Exec(sql, values...).Error
}

// IsRecordNotFound checks if error is record not found
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicateEntry checks if error is duplicate entry
func IsDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	return contains(err.Error(), "Duplicate entry") || contains(err.Error(), "duplicate key")
}

// contains checks if string contains substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr || 
		     containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetDB returns the underlying gorm.DB instance
func (db *DB) GetDB() *gorm.DB {
	return db.DB
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetStats returns database connection statistics
func (db *DB) GetStats() (map[string]interface{}, error) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, err
	}
	
	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}