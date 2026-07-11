package entity

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestTodoUsesSoftDelete(t *testing.T) {
	parsed, err := schema.Parse(&Todo{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("parse Todo schema: %v", err)
	}

	if len(parsed.QueryClauses) == 0 {
		t.Fatal("Todo schema does not filter soft-deleted records")
	}
	if len(parsed.DeleteClauses) == 0 {
		t.Fatal("Todo schema does not use soft delete")
	}
}
