package db_util

import (
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestItem struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"uniqueIndex"`
	Value     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestParent struct {
	Id       uint        `gorm:"primaryKey;autoIncrement"`
	Name     string
	Children []TestChild `gorm:"foreignKey:ParentId"`
}

type TestChild struct {
	Id       uint `gorm:"primaryKey;autoIncrement"`
	ParentId uint
	Label    string
}

type testDB struct {
	db *gorm.DB
}

func (t *testDB) Get() (*gorm.DB, error) {
	return t.db, nil
}

func setupTestDB(t *testing.T) *testDB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&TestItem{}, &TestParent{}, &TestChild{}); err != nil {
		t.Fatal(err)
	}
	return &testDB{db: db}
}

func TestCreate(t *testing.T) {
	tdb := setupTestDB(t)
	item := &TestItem{Name: "test1", Value: 42}
	if err := Create(tdb, item); err != nil {
		t.Fatal(err)
	}
	if item.Id == 0 {
		t.Fatal("expected id to be set after create")
	}
}

func TestFind(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "a", Value: 1})
	Create(tdb, &TestItem{Name: "b", Value: 2})

	var items []TestItem
	if err := Find(tdb, &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestFirst(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "first", Value: 10})
	Create(tdb, &TestItem{Name: "second", Value: 20})

	var item TestItem
	if err := First(tdb, &item); err != nil {
		t.Fatal(err)
	}
	if item.Name != "first" {
		t.Fatalf("expected 'first', got '%s'", item.Name)
	}
}

func TestLast(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "first", Value: 10})
	Create(tdb, &TestItem{Name: "second", Value: 20})

	var item TestItem
	if err := Last(tdb, &item); err != nil {
		t.Fatal(err)
	}
	if item.Name != "second" {
		t.Fatalf("expected 'second', got '%s'", item.Name)
	}
}

func TestSave(t *testing.T) {
	tdb := setupTestDB(t)
	item := &TestItem{Name: "original", Value: 1}
	Create(tdb, item)

	item.Value = 99
	if err := Save(tdb, item); err != nil {
		t.Fatal(err)
	}

	var found TestItem
	First(tdb, &found, "id = ?", item.Id)
	if found.Value != 99 {
		t.Fatalf("expected value 99, got %d", found.Value)
	}
}

func TestTimestampsOnCreateAndSave(t *testing.T) {
	tdb := setupTestDB(t)
	before := time.Now().Add(-time.Second)
	item := &TestItem{Name: "timestamps", Value: 1}
	Create(tdb, item)

	var created TestItem
	First(tdb, &created, "id = ?", item.Id)
	if created.CreatedAt.Before(before) {
		t.Fatal("expected created_at to be set on create")
	}
	if created.UpdatedAt.Before(before) {
		t.Fatal("expected updated_at to be set on create")
	}

	originalCreatedAt := created.CreatedAt
	originalUpdatedAt := created.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	created.Value = 99
	if err := Save(tdb, &created); err != nil {
		t.Fatal(err)
	}

	var updated TestItem
	First(tdb, &updated, "id = ?", item.Id)
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		t.Fatalf("expected created_at to be unchanged, was %v now %v", originalCreatedAt, updated.CreatedAt)
	}
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Fatalf("expected updated_at to advance on save, was %v now %v", originalUpdatedAt, updated.UpdatedAt)
	}
}

func TestDelete(t *testing.T) {
	tdb := setupTestDB(t)
	item := &TestItem{Name: "to_delete", Value: 1}
	Create(tdb, item)

	if err := Delete(tdb, item); err != nil {
		t.Fatal(err)
	}

	count, err := Count(tdb, &TestItem{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 items after delete, got %d", count)
	}
}

func TestCount(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "a", Value: 1})
	Create(tdb, &TestItem{Name: "b", Value: 2})
	Create(tdb, &TestItem{Name: "c", Value: 3})

	count, err := Count(tdb, &TestItem{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestFindPreload(t *testing.T) {
	tdb := setupTestDB(t)
	parent := &TestParent{Name: "parent1"}
	Create(tdb, parent)
	Create(tdb, &TestChild{ParentId: parent.Id, Label: "child1"})
	Create(tdb, &TestChild{ParentId: parent.Id, Label: "child2"})

	var parents []TestParent
	if err := FindPreload(tdb, []string{"Children"}, &parents); err != nil {
		t.Fatal(err)
	}
	if len(parents) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(parents))
	}
	if len(parents[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parents[0].Children))
	}
}

func TestRetryFind(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "retry_test", Value: 42})

	var items []TestItem
	if err := RetryFind(tdb.db, &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestCreateMany(t *testing.T) {
	tdb := setupTestDB(t)
	objects := []TestItem{
		{Name: "bulk1", Value: 1},
		{Name: "bulk2", Value: 2},
		{Name: "bulk3", Value: 3},
	}
	if err := CreateMany(tdb, &objects); err != nil {
		t.Fatal(err)
	}

	count, err := Count(tdb, &TestItem{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 items after bulk insert, got %d", count)
	}
}

func TestCreateManyTimestamps(t *testing.T) {
	tdb := setupTestDB(t)
	before := time.Now().Add(-time.Second)
	objects := []TestItem{
		{Name: "ts1", Value: 1},
	}
	if err := CreateMany(tdb, &objects); err != nil {
		t.Fatal(err)
	}

	var item TestItem
	First(tdb, &item)
	if item.CreatedAt.Before(before) {
		t.Fatal("expected created_at to be set on create")
	}
	if item.UpdatedAt.Before(before) {
		t.Fatal("expected updated_at to be set on create")
	}

	originalCreatedAt := item.CreatedAt
	originalUpdatedAt := item.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	item.Value = 99
	if err := Save(tdb, &item); err != nil {
		t.Fatal(err)
	}

	var updated TestItem
	First(tdb, &updated, "id = ?", item.Id)
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		t.Fatalf("expected created_at unchanged after save, was %v now %v", originalCreatedAt, updated.CreatedAt)
	}
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Fatalf("expected updated_at to advance on save, was %v now %v", originalUpdatedAt, updated.UpdatedAt)
	}
}

func TestIsRecordNotFoundError(t *testing.T) {
	if !IsRecordNotFoundError(gorm.ErrRecordNotFound) {
		t.Fatal("expected gorm.ErrRecordNotFound to be detected")
	}
	if !IsRecordNotFoundError(errors.New("record not found")) {
		t.Fatal("expected string match to be detected")
	}
	if IsRecordNotFoundError(errors.New("some other error")) {
		t.Fatal("expected false for unrelated error")
	}
}

func TestIsDuplicateEntryError(t *testing.T) {
	if !IsDuplicateEntryError(errors.New("Duplicate entry 'foo' for key 'name'")) {
		t.Fatal("expected duplicate entry to be detected")
	}
	if IsDuplicateEntryError(errors.New("some other error")) {
		t.Fatal("expected false for unrelated error")
	}
}

func TestFirstRecordNotFound(t *testing.T) {
	tdb := setupTestDB(t)
	var item TestItem
	err := First(tdb, &item)
	if err == nil {
		t.Fatal("expected error for empty table")
	}
	if !IsRecordNotFoundError(err) {
		t.Fatalf("expected record not found error, got: %v", err)
	}
}

func TestResult(t *testing.T) {
	tdb := setupTestDB(t)
	Create(tdb, &TestItem{Name: "result_test", Value: 1})

	db, _ := tdb.Get()
	var item TestItem
	r := db.First(&item)
	result, err := Result(r)
	if err != nil {
		t.Fatal(err)
	}
	if result.Error != nil {
		t.Fatal(result.Error)
	}
}

func TestResultError(t *testing.T) {
	tdb := setupTestDB(t)
	db, _ := tdb.Get()
	var item TestItem
	r := db.First(&item)
	_, err := Result(r)
	if err == nil {
		t.Fatal("expected error from Result on empty table")
	}
}
