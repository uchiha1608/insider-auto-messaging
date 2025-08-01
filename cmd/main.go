package main

import (
	"net/http"

	_ "github.com/lib/pq"
	"insider-auto-messaging/config"
	"insider-auto-messaging/controller"
	"insider-auto-messaging/repository"
	"insider-auto-messaging/scheduler"
	"insider-auto-messaging/service"
)

// @title Insider Messaging API
// @version 1.0
// @description Automatically sends messages every 2 minutes and logs them.
// @host localhost:8080
// @BasePath /
func main() {
	db := config.InitDB()
	rdb := config.InitRedis()

	repo := &repository.MessageRepository{DB: db}
	svc := &service.MessageService{Repo: repo, Redis: rdb}
	sched := &scheduler.Scheduler{Service: svc}
	ctrl := &controller.MessageController{Scheduler: sched, Repo: repo}

	http.HandleFunc("/start", ctrl.Start)
	http.HandleFunc("/stop", ctrl.Stop)
	http.HandleFunc("/sent", ctrl.SentMessages)

	sched.Start() // Start automatically on deploy

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
