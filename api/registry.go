package api

import "github.com/VoltFramework/volt/task"

type Registry interface {
	// Register a new task by ID in the configured registry
	Register(string, *task.Task) error

	// Tasks returns all tasks in the registry
	Tasks() ([]*task.Task, error)

	// Delete removes the task by ID from the registry
	Delete(string) error

	// Fetch returns a specific task in the registry by ID
	Fetch(string) (*task.Task, error)

	// Update finds the task in the registry for the ID and updates it's data
	Update(string, *task.Task) error
}
