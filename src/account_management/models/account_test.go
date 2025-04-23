//go:build unit

package models_test

import (
	"testing"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/account_management/models"
)

func TestAccount(t *testing.T) {
	test_account := models.Account{
		ID: 1,
		OwnerID: 1,
		UserRoles: []models.UserRole{
			models.UserRole{UserID: 1, Role: "admin"},
			models.UserRole{UserID: 2, Role: "editor"},
			models.UserRole{UserID: 3, Role: "reader"},
		},
		Name: "Test Account",
		Balance: 100.00,
		EditorLimit: 200.00,
		Currency: "USD",
		CreatedAt: "2021-01-01",
		UpdatedAt: "2021-01-01",
		Active: true,
		Current: true,
	}

	if test_account.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", test_account.ID)
	}
	if test_account.OwnerID != 1 {
		t.Errorf("Expected OwnerID to be 1, got %d", test_account.OwnerID)
	}
	if len(test_account.UserRoles) != 3 {
		t.Errorf("Expected 3 UserRoles, got %d", len(test_account.UserRoles))
	}
	if test_account.Name != "Test Account" {
		t.Errorf("Expected Name to be Test Account, got %s", test_account.Name)
	}
	if test_account.Balance != 100.00 {
		t.Errorf("Expected Balance to be 100.00, got %f", test_account.Balance)
	}
	if test_account.EditorLimit != 200.00 {
		t.Errorf("Expected EditorLimit to be 200.00, got %f", test_account.EditorLimit)
	}
	if test_account.Currency != "USD" {
		t.Errorf("Expected Currency to be USD, got %s", test_account.Currency)
	}
	if test_account.CreatedAt != "2021-01-01" {
		t.Errorf("Expected CreatedAt to be 2021-01-01, got %s", test_account.CreatedAt)
	}
	if test_account.UpdatedAt != "2021-01-01" {
		t.Errorf("Expected UpdatedAt to be 2021-01-01, got %s", test_account.UpdatedAt)
	}
	if !test_account.Active {
		t.Errorf("Expected Active to be true, got false")
	}
	if !test_account.Current {
		t.Errorf("Expected Current to be true, got false")
	}

}
