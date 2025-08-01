package scheduler

import (
	"insider-auto-messaging/service"
	"sync"
	"time"
)

type Scheduler struct {
	stopChan chan struct{}
	running  bool
	lock     sync.Mutex
	Service  *service.MessageService
}

func (s *Scheduler) Start() {
	s.lock.Lock()
	if s.running {
		s.lock.Unlock()
		return
	}
	s.running = true
	s.stopChan = make(chan struct{})
	s.lock.Unlock()

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.Service.SendPendingMessages()
				time.Sleep(2 * time.Minute)
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.running {
		close(s.stopChan)
		s.running = false
	}
}
