package storage

import "todo-api/models"

// Store interface that both MemoryStore and PostgresStore will implement
type Store interface {
	Create(req models.CreateTodoRequest) (*models.TODO, error)
	GetAll() ([]*models.TODO, error)
	GetByID(id int) (*models.TODO, error)
	Update(id int, req models.UpdateTodoRequest) (*models.TODO, error)
	Delete(id int) error
}
