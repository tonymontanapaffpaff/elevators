package entities

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// time between state changes
const storageDelay int = 250

type MoveRequest struct {
	Level int
	Up    bool
}

const (
	StateStart = iota
	StateRunning
)

type Storage struct {
	Elevators         []*Elevator
	Queue             []*MoveRequest
	FloorCount        int
	State             int
	QueueMutex        *sync.Mutex
}

func NewStorage(floors int, count int) *Storage {
	b := Storage{
		Elevators:         make([]*Elevator, 0),
		Queue:             make([]*MoveRequest, 0),
		FloorCount:        floors,
		State:             StateStart,
		QueueMutex:        &sync.Mutex{},
	}
	// create elevators on the first floor
	for i := 0; i < count; i++ {
		b.Elevators = append(b.Elevators, NewElevator("Elevator"+strconv.Itoa(i+1), 1, floors))
	}
	return &b
}

func (s *Storage) RunStorage() {
	for {
		switch s.State {
		case StateStart:
			fmt.Println("Storage started")
			for i := 0; i < len(s.Elevators); i++ {
				e := s.Elevators[i]
				go e.RunElevator()
			}
			s.State = StateRunning
			break
		case StateRunning:
			// if we have a queued request see can we assign an idle elevator
			if s.HasQueue() {
				r := s.PeekRequest()
				if r != nil {
					e := s.GetIdleElevatorClosest(r.Level)
					if e != nil {
						q := s.DequeueRequest()
						fmt.Printf("Storage sending elevator to %d floor\n", q.Level)
						e.PushButton(q.Level)
					}
				}
			}
			break
		}
		time.Sleep(time.Duration(storageDelay) * time.Millisecond)
	}
}

func (s *Storage) GetIdleElevator() *Elevator {
	for i := 0; i < len(s.Elevators); i++ {
		if s.Elevators[i].State == StateIdle {
			return s.Elevators[i]
		}
	}
	return nil
}

func (s *Storage) GetIdleElevatorClosest(requestedLevel int) *Elevator {
	// get closest elevator to requested level
	var elevators []*Elevator
	var closestDist = 0.0
	for _, e := range s.Elevators {
		dist := math.Abs(float64(e.Level - requestedLevel))
		// if there is an elevator on this floor return it immediately
		if dist == 0 && e.State == StateReady {
			return e
		}
		if e.State == StateIdle {
			if len(elevators) == 0 {
				closestDist = dist
				elevators = append(elevators, e)
			} else if dist < closestDist {
				closestDist = dist
				elevators = nil
				elevators = append(elevators, e)
			} else if dist == closestDist {
				elevators = append(elevators, e)
			}
		}
	}
	if len(elevators) == 0 {
		return nil
	}
	r := rand.Intn(len(elevators))
	return elevators[r]
}

func (s *Storage) GetElevator(level int) *Elevator {
	// look through elevators and see if one of this is at the requested floor
	var rl []*Elevator
	for _, e := range s.Elevators {
		if e.ReadyAtLevel(level) {
			rl = append(rl, e)
		}
	}
	if len(rl) == 0 {
		return nil
	}
	r := rand.Intn(len(rl))
	return rl[r]
}

func (s *Storage) RequestLift(level int, up bool) {
	// dont' queue if elevator at floor and ready
	if s.HasElevatorReady(level) {
		return
	}
	// queue requesting floor
	s.QueueMutex.Lock()
	s.Queue = append(s.Queue, &MoveRequest{
		Level: level,
		Up:    up,
	})
	s.QueueMutex.Unlock()
}

func (s *Storage) HasElevatorReady(level int) bool {
	for _, e := range s.Elevators {
		if e.State == StateReady && e.Level == level {
			return true
		}
	}
	return false
}

func (s *Storage) InQueue(level int) bool {
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()
	for _, l := range s.Queue {
		if l.Level == level {
			return true
		}
	}
	return false
}

func (s *Storage) HasQueue() bool {
	s.QueueMutex.Lock()
	r := len(s.Queue) > 0
	s.QueueMutex.Unlock()
	return r
}

func (s *Storage) PeekRequest() *MoveRequest {
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()
	if len(s.Queue) > 0 {
		return s.Queue[len(s.Queue)-1]
	}
	return nil
}

func (s *Storage) DequeueRequest() *MoveRequest {
	s.QueueMutex.Lock()
	defer s.QueueMutex.Unlock()
	if len(s.Queue) > 0 {
		q := s.Queue[0]
		s.Queue = s.Queue[1:]
		return q
	}
	return nil
}
