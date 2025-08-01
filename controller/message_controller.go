package controller

import (
	"encoding/json"
	"net/http"

	"insider-auto-messaging/repository"
	"insider-auto-messaging/scheduler"
)

type MessageController struct {
	Scheduler scheduler.Controller
	Repo      repository.MessageRepo
}

// Start godoc
// @Summary Start auto message sending
// @Description Starts the background message-sending scheduler.
// @Tags Control
// @Success 200 {string} string "Scheduler started"
// @Router /start [get]
func (c *MessageController) Start(w http.ResponseWriter, r *http.Request) {
	c.Scheduler.Start()
	w.Write([]byte("Scheduler started"))
}

// Stop godoc
// @Summary Stop auto message sending
// @Description Stops the background message-sending scheduler.
// @Tags Control
// @Success 200 {string} string "Scheduler stopped"
// @Router /stop [get]
func (c *MessageController) Stop(w http.ResponseWriter, r *http.Request) {
	c.Scheduler.Stop()
	w.Write([]byte("Scheduler stopped"))
}

// SentMessages godoc
// @Summary List sent messages
// @Description Returns all messages that have been successfully sent.
// @Tags Messages
// @Produce json
// @Success 200 {array} model.Message
// @Router /sent [get]
func (c *MessageController) SentMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := c.Repo.GetAllSent()
	if err != nil {
		http.Error(w, "Error fetching messages", 500)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
