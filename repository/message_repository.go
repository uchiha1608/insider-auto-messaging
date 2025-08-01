package repository

import (
	"database/sql"
	"insider-auto-messaging/model"
	"time"
)

type MessageRepository struct {
	DB *sql.DB
}

func (r *MessageRepository) GetUnsentMessages(limit int) ([]model.Message, error) {
	rows, err := r.DB.Query("SELECT id, to, content, is_sent FROM messages WHERE is_sent = FALSE LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.To, &msg.Content, &msg.IsSent); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepository) MarkAsSent(id int64, messageId string) error {
	_, err := r.DB.Exec("UPDATE messages SET is_sent = TRUE, message_id = $1, sent_at = $2 WHERE id = $3", messageId, time.Now(), id)
	return err
}

func (r *MessageRepository) GetAllSent() ([]model.Message, error) {
	rows, err := r.DB.Query("SELECT id, to, content, message_id, sent_at FROM messages WHERE is_sent = TRUE ORDER BY sent_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.To, &msg.Content, &msg.MessageID, &msg.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
