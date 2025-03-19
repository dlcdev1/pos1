package entity

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID           string
	Title        string
	Description  string
	Status       string
	DateFinished time.Time
}

func NewTodo() *Todo {
	return &Todo{
		ID: uuid.New().String(),
	}
}

func (t *Todo) Done() {
	t.Status = "done"
	t.DateFinished = time.Now()
}

func (t *Todo) Undone() {
	t.Status = "pending"
	t.DateFinished = time.Time{}
}
