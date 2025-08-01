package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"context"
	"github.com/redis/go-redis/v9"
	_ "insider-auto-messaging/model"
	"insider-auto-messaging/repository"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MessageService struct {
	Repo       repository.MessageRepo
	Redis      *redis.Client
	HTTPClient HTTPClient
}

type Payload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type Response struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func (s *MessageService) SendPendingMessages() {
	msgs, err := s.Repo.GetUnsentMessages(2)
	if err != nil {
		fmt.Println("Error retrieving unsent messages:", err)
		return
	}

	for _, msg := range msgs {
		payload := Payload{To: msg.To, Content: msg.Content}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "https://webhook.site/c3f13233-1ed4-429e-9649-8133b3b9c9cd", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-ins-auth-key", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")

		resp, err := s.HTTPClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusAccepted {
			fmt.Println("Error sending message:", err)
			continue
		}

		var res Response
		json.NewDecoder(resp.Body).Decode(&res)
		s.Repo.MarkAsSent(msg.ID, res.MessageID)

		if s.Redis != nil {
			s.Redis.Set(context.Background(), res.MessageID, time.Now().Format(time.RFC3339), 0)
		}
	}
}
