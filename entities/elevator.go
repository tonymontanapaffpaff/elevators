package entities

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// time between state changes
const elevatorsDelay int = 250

// Elevator States
const (
	StateIdle = iota
	StateCheckButton
	StateReady
	StateMoving
)

type Elevator struct {
	Title         string
	Position      float64
	Goal          int
	Level         int
	Valid         bool
	Speed         float64
	State         int
	Buttons       []bool
	ButtonMutex   *sync.Mutex
	PositionMutex *sync.Mutex
}

func NewElevator(title string, start int, floors int) *Elevator {
	e := Elevator{
		Title:         title,
		Position:      float64(start),
		Goal:          start,
		Level:         start,
		Valid:         true,
		Speed:         0.25,
		State:         StateIdle,
		Buttons:       make([]bool, floors),
		ButtonMutex:   &sync.Mutex{},
		PositionMutex: &sync.Mutex{},
	}
	return &e
}

func (e *Elevator) RunElevator() {
	for true {
		switch e.State {
		case StateIdle:
			if e.HasButtonPressed() {
				e.State = StateCheckButton
			}
			break
		case StateCheckButton:
			// loops through our buttons and finds closest destination
			e.Goal = e.GetGoalLevel()
			if e.Goal != e.Level {
				e.Valid = false
				fmt.Printf("- %s moving towards %d\n", e.Title, e.Goal)
				e.State = StateMoving
			} else {
				e.ResetButton(e.Goal)
				fmt.Printf("- %s opening doors\n", e.Title)
				e.State = StateReady
			}
			break
		case StateReady:
			// necessary waiting period at a level before going idle or moving
			time.Sleep(2 * time.Second)
			fmt.Printf("- %s closing doors\n", e.Title)
			e.State = StateIdle
			break
		case StateMoving:
			// when a button has been pressed
			if !e.Valid {
				e.Move()
			} else {
				fmt.Printf("- %s opening doors\n", e.Title)
				e.State = StateReady
			}
			break
		}
		time.Sleep(time.Duration(elevatorsDelay) * time.Millisecond)
	}
}

func (e *Elevator) Move() {
	currentPosition := math.Round(e.Position*100) / 100
	goal := math.Round(float64(e.Goal))
	// check where we move
	if currentPosition > goal {
		e.PositionMutex.Lock()
		e.Position -= e.Speed
		e.PositionMutex.Unlock()
		e.Valid = false
		fmt.Printf("- %s moving down\n", e.Title)
	} else if currentPosition < goal {
		e.PositionMutex.Lock()
		e.Position += e.Speed
		e.PositionMutex.Unlock()
		e.Valid = false
		fmt.Printf("- %s moving up\n", e.Title)
	} else {
		fmt.Printf("- %s arrived at %d\n", e.Title, e.Goal)
		e.PositionMutex.Lock()
		e.Position = goal
		e.PositionMutex.Unlock()
		e.Level = e.Goal
		e.Valid = true
		e.ResetButton(e.Goal)
	}
}

func (e *Elevator) GetGoalLevel() int {
	var closestDist = 99999.0
	var closestFloor = 0

	e.ButtonMutex.Lock()
	defer e.ButtonMutex.Unlock()

	for i := 0; i < len(e.Buttons); i++ {
		buttonFloor := i + 1
		dist := math.Abs(float64(e.Level - buttonFloor))
		if e.Buttons[i] && (closestFloor == 0 || dist < closestDist) {
			closestDist = dist
			closestFloor = buttonFloor
		}
	}

	if closestFloor > 0 {
		return closestFloor
	}

	return e.Level
}

func (e *Elevator) HasButtonPressed() bool {
	e.ButtonMutex.Lock()
	for i := 0; i < len(e.Buttons); i++ {
		if e.Buttons[i] {
			e.ButtonMutex.Unlock()
			return true
		}
	}
	e.ButtonMutex.Unlock()
	return false
}

func (e *Elevator) ReadyAtLevel(level int) bool {
	return e.State == StateReady && e.Level == level
}

func (e *Elevator) PushButton(level int) {
	fmt.Printf("- %s button requested for %d\n", e.Title, level)
	e.ButtonMutex.Lock()
	e.Buttons[level-1] = true
	e.ButtonMutex.Unlock()
}

func (e *Elevator) ResetButton(level int) {
	e.ButtonMutex.Lock()
	e.Buttons[level-1] = false
	e.ButtonMutex.Unlock()
}
