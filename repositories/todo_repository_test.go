package repositories

import (
	"go-todo/tests/utils/database"
	"testing"
)

func TestDBConnection(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect databese: %v", err)
	}
	defer db.Close()
}
