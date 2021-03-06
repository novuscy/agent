package services

import (
	"sync"
	"time"
)

const (
	timeout  = 3
	interval = 10000

	online  = "online"
	offline = "offline"

	service = "service"
	device  = "device"
)

type Service struct {
	Name     string
	LastSeen time.Time
	Status   string
	Type     string

	counter int
	done    chan bool
	ticker  *time.Ticker
	mu      sync.Mutex
}

func NewService(name, svctype string) *Service {
	ticker := time.NewTicker(interval * time.Millisecond)
	done := make(chan bool)
	s := Service{Name: name, Status: online, Type: svctype, done: done, counter: timeout, ticker: ticker}
	s.Listen()
	return &s
}

func (s *Service) Listen() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// TODO - we can disable ticker when the status gets OFFLINE
				// and on the next heartbeat enable it again
				s.mu.Lock()
				s.counter = s.counter - 1
				if s.counter == 0 {
					s.Status = offline
					s.counter = timeout
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Service) Update() {
	s.LastSeen = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = timeout
	s.Status = online
}
