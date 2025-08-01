package model

type Message struct {
	ID        int64  `db:"id"`
	To        string `db:"to"`
	Content   string `db:"content"`
	IsSent    bool   `db:"is_sent"`
	SentAt    string `db:"sent_at"`
	MessageID string `db:"message_id"`
}
