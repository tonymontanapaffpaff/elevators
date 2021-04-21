package entities

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// time between state changes
const workerDelay int = 250

type Activity struct {
	Goal    int
	Seconds int
}

// Worker States
const (
	StateLazing = iota
	StateRequest
	StateWaiting
	StateRiding
	StateWorking
)

type Worker struct {
	Name      string
	Level     int
	Goal      int
	State     int
	Schedule  []*Activity
	Storage   *Storage
	Elevator  *Elevator
	WaitGroup *sync.WaitGroup
}

func NewWorker(name string, level int, b *Storage) *Worker {
	p := Worker{
		Name:     name,
		Level:    level,
		State:    StateLazing,
		Storage:  b,
		Elevator: nil,
	}
	return &p
}

func (p *Worker) RunWorker() {
	for {
		switch p.State {
		case StateLazing:
			if len(p.Schedule) > 0 {
				p.State = StateRequest
			} else {
				fmt.Printf("%s leaving\n", p.Name)
				return
			}
			break
		case StateRequest:
			fmt.Printf("%s requests a lift\n", p.Name)
			p.MakeRequest(p.Storage)
			p.State = StateWaiting
			break
		case StateWaiting:
			e := p.Storage.GetElevator(p.Level)
			if e != nil {
				p.Elevator = e
				goal := p.Schedule[0].Goal
				fmt.Printf("%s's elevator arrived at %d floor\n", p.Name, goal)
				p.Elevator.PushButton(goal)
				p.State = StateRiding
			}
			break
		case StateRiding:
			goal := p.Schedule[0].Goal
			if p.Elevator.ReadyAtLevel(goal) {
				fmt.Printf("%s's elevator ready at level\n", p.Name)
				p.Elevator = nil
				p.Level = goal
				p.State = StateWorking
			}
			break
		case StateWorking:
			// get duration worker should be on this level
			s := time.Duration(p.Schedule[0].Seconds)
			fmt.Printf("%s is working for %s seconds\n", p.Name, strconv.Itoa(p.Schedule[0].Seconds))
			time.Sleep(s * time.Second) // work for that time
			// Remove schedule
			p.Schedule = p.Schedule[1:]
			fmt.Printf("%s done his work, going idle\n", p.Name)
			p.State = StateLazing
			break
		}
		time.Sleep(time.Duration(workerDelay) * time.Millisecond)
	}
}

func (p *Worker) AddActivity(level int, duration int) error {
	o := Activity{Goal: level, Seconds: duration}
	if level > p.Storage.FloorCount || level < 1 {
		return fmt.Errorf("introduced non-existent floor %d", level)
	}
	p.Schedule = append(p.Schedule, &o)
	return nil
}

func (p *Worker) SetGoal(level int) {
	p.Goal = level
}

func (p *Worker) MakeRequest(b *Storage) {
	next := p.Schedule[0].Goal
	up := next > p.Level
	b.RequestLift(p.Level, up)
}
