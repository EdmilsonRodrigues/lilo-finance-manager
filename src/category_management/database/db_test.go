//go:build unit

package database_test

import (
	"testing"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/database"
	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"gorm.io/driver/sqlite"
)

func TestOpenDBConnection(t *testing.T) {
	dsn := ":memory:"
	opener := sqlite.Open
	db := database.OpenDBConnection(dsn, opener)
	if db == nil {
		t.Error("Expected a non-nil database connection")
	}
}

func TestMakeMigrations(t *testing.T) {
    fakeMigrator := &fakeMigrator{}
    database.MakeMigrations(fakeMigrator)
    if len(fakeMigrator.CallValues) != 1 {
        t.Errorf("Expected 1 call to AutoMigrate, got %d", len(fakeMigrator.CallValues))
    }
    args := fakeMigrator.CallValues[0].([]interface{})
    if len(args) != 1 {
        t.Fatalf("Expected the slice argument to AutoMigrate to have 1 element, got %d", len(args))
    }
    if _, ok := args[0].(*models.Category); !ok {
        t.Errorf("Expected AutoMigrate to be called with a *models.Category, got %T", args[0])
    }
}

type fakeMigrator struct{
    Error error
    CallValues []interface{}
}

func (m *fakeMigrator) AutoMigrate(dst ...interface{}) error {
    m.CallValues = append(m.CallValues, dst)
    return m.Error
}
