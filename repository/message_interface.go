package repository

import "insider-auto-messaging/model"

type MessageRepo interface {
	GetUnsentMessages(limit int) ([]model.Message, error)
	MarkAsSent(id int64, messageId string) error
	GetAllSent() ([]model.Message, error)
}
