package storage

import (
	"errors"
	"sync"
	"time"
	"todo-api/models"
)

type MemoryStore struct {
	todos  map[int]*models.TODO
	nextID int
	mutex  sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		todos:  make(map[int]*models.TODO),
		nextID: 1,
	}
}

func (s *MemoryStore) Create(req models.CreateTodoRequest) *models.TODO {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	todo := &models.TODO{
		ID:          s.nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.todos[s.nextID] = todo
	s.nextID++
	return todo
}
func (s *MemoryStore) GetAll() []*models.TODO {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todos := make([]*models.TODO, len(s.todos))

	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

func (s *MemoryStore) GetByID(id int) (*models.TODO, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}
func (s *MemoryStore) Delete(id int) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.todos[id]
	if !ok {
		return errors.New("todo not found")
	}
	delete(s.todos, id)
	return nil
}
func (s *MemoryStore) Update(id int, req models.UpdateTodoRequest) (*models.TODO, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, errors.New("todo not found")
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	todo.UpdatedAt = time.Now()
	return todo, nil
}
