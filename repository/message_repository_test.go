package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"

	"insider-auto-messaging/repository"
)

func TestGetUnsentMessages(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := &repository.MessageRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "to", "content", "is_sent"}).
		AddRow(1, "+1234567890", "Test message", false)

	mock.ExpectQuery("SELECT (.+) FROM messages WHERE is_sent = FALSE").
		WillReturnRows(rows)

	messages, err := repo.GetUnsentMessages(2)

	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "+1234567890", messages[0].To)
}

func TestMarkAsSent(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := &repository.MessageRepository{DB: db}

	mock.ExpectExec("UPDATE messages SET is_sent = TRUE").
		WithArgs("abc-123", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.MarkAsSent(1, "abc-123")
	assert.NoError(t, err)
}
