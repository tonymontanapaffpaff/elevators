package server

import (
	"fmt"
	"github.com/tonymontanapaffpaff/elevators/entities"
)

type WorkerSchedulePair struct {
	Floor   int
	Seconds int
}
type WorkerRequest struct {
	Name     string
	Schedule []WorkerSchedulePair
}
type WorkerResponse struct {
	Message string
}

type Server struct {
	storage *entities.Storage
}

func New(floorCount int, elevatorCount int) *Server {
	// build server
	s := Server{entities.NewStorage(floorCount, elevatorCount)}
	// start the elevator storage
	go s.storage.RunStorage()

	return &s
}

func (s *Server) AddWorker(req WorkerRequest, res *WorkerResponse) error {
	// create the worker and his schedule
	worker := entities.NewWorker(req.Name, 1, s.storage)
	for _, schedule := range req.Schedule {
		err := worker.AddActivity(schedule.Floor, schedule.Seconds)
		if err != nil {
			return fmt.Errorf("failed to create activity, err: %w", err)
		}
	}
	go worker.RunWorker()
	res.Message = req.Name + " started doing his job"
	return nil
}
