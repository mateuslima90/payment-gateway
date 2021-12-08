package repository

import (
	"github.com/mateuslima90/payment-gateway/adapter/repository/fixture"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTransactionRepositoryDb_Insert(t *testing.T) {
	migrationsDir := os.DirFS("fixture/sql")
	db := fixture.Up(migrationsDir)
	defer fixture.Down(db, migrationsDir)

	repository := NewTransactionRepositoryDb(db)
	err := repository.Insert("1", "1", 2, "approved", "")
	assert.Nil(t, err)
}
