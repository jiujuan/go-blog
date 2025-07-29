package database

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestModel for testing purposes
type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Age  int
}

func setupTestDB(t *testing.T) *DB {
	// Use in-memory SQLite for testing
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate test model
	if err := gormDB.AutoMigrate(&TestModel{}); err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	return NewDB(gormDB)
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)

	model := &TestModel{Name: "Test", Age: 25}
	err := db.Create(model)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if model.ID == 0 {
		t.Error("Expected ID to be set after create")
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	original := &TestModel{Name: "Test", Age: 25}
	db.Create(original)

	// Test GetByID
	var result TestModel
	err := db.GetByID(&result, original.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if result.Name != original.Name || result.Age != original.Age {
		t.Error("Retrieved data doesn't match original")
	}
}

func TestGetByField(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	original := &TestModel{Name: "UniqueTest", Age: 30}
	db.Create(original)

	// Test GetByField
	var result TestModel
	err := db.GetByField(&result, "name", "UniqueTest")
	if err != nil {
		t.Errorf("GetByField failed: %v", err)
	}

	if result.ID != original.ID {
		t.Error("Retrieved record doesn't match original")
	}
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	model := &TestModel{Name: "Original", Age: 25}
	db.Create(model)

	// Update
	model.Name = "Updated"
	model.Age = 30
	err := db.Update(model)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Verify update
	var result TestModel
	db.GetByID(&result, model.ID)
	if result.Name != "Updated" || result.Age != 30 {
		t.Error("Update didn't persist correctly")
	}
}

func TestUpdateFields(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	model := &TestModel{Name: "Original", Age: 25}
	db.Create(model)

	// Update specific fields
	fields := map[string]interface{}{
		"name": "FieldUpdated",
		"age":  35,
	}
	err := db.UpdateFields(&TestModel{}, model.ID, fields)
	if err != nil {
		t.Errorf("UpdateFields failed: %v", err)
	}

	// Verify update
	var result TestModel
	db.GetByID(&result, model.ID)
	if result.Name != "FieldUpdated" || result.Age != 35 {
		t.Error("UpdateFields didn't persist correctly")
	}
}

func TestExists(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	model := &TestModel{Name: "ExistsTest", Age: 25}
	db.Create(model)

	// Test exists - should return true
	exists, err := db.Exists(&TestModel{}, "name = ?", "ExistsTest")
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected record to exist")
	}

	// Test exists - should return false
	exists, err = db.Exists(&TestModel{}, "name = ?", "NonExistent")
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected record to not exist")
	}
}

func TestCount(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	models := []*TestModel{
		{Name: "Count1", Age: 25},
		{Name: "Count2", Age: 30},
		{Name: "Count3", Age: 25},
	}
	for _, model := range models {
		db.Create(model)
	}

	// Test count all
	count, err := db.Count(&TestModel{})
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// Test count with condition
	count, err = db.Count(&TestModel{}, "age = ?", 25)
	if err != nil {
		t.Errorf("Count with condition failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestList(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	models := []*TestModel{
		{Name: "List1", Age: 25},
		{Name: "List2", Age: 30},
		{Name: "List3", Age: 35},
		{Name: "List4", Age: 40},
		{Name: "List5", Age: 45},
	}
	for _, model := range models {
		db.Create(model)
	}

	// Test list with pagination
	var results []TestModel
	options := &QueryOptions{
		Page:    1,
		Limit:   3,
		OrderBy: "age ASC",
	}

	pagination, err := db.List(&results, options)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	if pagination.Total != 5 {
		t.Errorf("Expected total 5, got %d", pagination.Total)
	}

	if pagination.TotalPages != 2 {
		t.Errorf("Expected 2 total pages, got %d", pagination.TotalPages)
	}

	// Verify ordering
	if results[0].Age != 25 || results[1].Age != 30 || results[2].Age != 35 {
		t.Error("Results not properly ordered")
	}
}

func TestListWithFilters(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	models := []*TestModel{
		{Name: "Filter1", Age: 25},
		{Name: "Filter2", Age: 25},
		{Name: "Filter3", Age: 30},
	}
	for _, model := range models {
		db.Create(model)
	}

	// Test list with filters
	var results []TestModel
	options := &QueryOptions{
		Page:  1,
		Limit: 10,
		Filters: map[string]interface{}{
			"age": 25,
		},
	}

	pagination, err := db.List(&results, options)
	if err != nil {
		t.Errorf("List with filters failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 filtered results, got %d", len(results))
	}

	if pagination.Total != 2 {
		t.Errorf("Expected filtered total 2, got %d", pagination.Total)
	}
}

func TestTransaction(t *testing.T) {
	db := setupTestDB(t)

	err := db.Transaction(func(tx *DB) error {
		model1 := &TestModel{Name: "Transaction1", Age: 25}
		if err := tx.Create(model1); err != nil {
			return err
		}

		model2 := &TestModel{Name: "Transaction2", Age: 30}
		if err := tx.Create(model2); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}

	// Verify both records were created
	count, _ := db.Count(&TestModel{})
	if count != 2 {
		t.Errorf("Expected 2 records after transaction, got %d", count)
	}
}

func TestIncrement(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	model := &TestModel{Name: "IncrementTest", Age: 25}
	db.Create(model)

	// Test increment
	err := db.Increment(&TestModel{}, model.ID, "age", 5)
	if err != nil {
		t.Errorf("Increment failed: %v", err)
	}

	// Verify increment
	var result TestModel
	db.GetByID(&result, model.ID)
	if result.Age != 30 {
		t.Errorf("Expected age 30 after increment, got %d", result.Age)
	}
}

func TestDecrement(t *testing.T) {
	db := setupTestDB(t)

	// Create test data
	model := &TestModel{Name: "DecrementTest", Age: 25}
	db.Create(model)

	// Test decrement
	err := db.Decrement(&TestModel{}, model.ID, "age", 5)
	if err != nil {
		t.Errorf("Decrement failed: %v", err)
	}

	// Verify decrement
	var result TestModel
	db.GetByID(&result, model.ID)
	if result.Age != 20 {
		t.Errorf("Expected age 20 after decrement, got %d", result.Age)
	}
}

func TestIsRecordNotFound(t *testing.T) {
	db := setupTestDB(t)

	var result TestModel
	err := db.GetByID(&result, 999) // Non-existent ID

	if !IsRecordNotFound(err) {
		t.Error("Expected IsRecordNotFound to return true for non-existent record")
	}
}