package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"todo-api/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (p *PostgresStore) Create(req models.CreateTodoRequest) (*models.TODO, error) {

	query := `
	INSERT INTO todos (title,description,completed,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id,title,description,completed,created_at,updated_at
	`

	now := time.Now()
	var todo models.TODO
	err := p.db.QueryRow(query, req.Title, req.Description, false, now, now).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	return &todo, nil

}

func (p *PostgresStore) Update(id int, req models.UpdateTodoRequest) (*models.TODO, error) {
	currentTodo, err := p.GetByID(id)
	if err != nil {
		return nil, err
	}
	setParts := []string{}
	args := []interface{}{}

	argPos := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argPos))
		args = append(args, req.Title)
		argPos++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argPos))
		args = append(args, req.Description)
		argPos++
	}
	if req.Completed != nil {
		setParts = append(setParts, fmt.Sprintf("completed = $%d", argPos))
		args = append(args, req.Completed)
		argPos++
	}
	if len(setParts) == 0 {
		// No fields to update, return current todo
		return currentTodo, nil
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	args = append(args, id)

	query := fmt.Sprintf(
		`
		UPDATE todos 
		SET %s
		WHERE id = $%d
		RETURNING id,title,description,completed,created_at,updated_at
		`,
		joinStrings(setParts, ", "), argPos,
	)
	var todo models.TODO
	err = p.db.QueryRow(query, args...).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	return &todo, nil
}
func (p *PostgresStore) GetAll() ([]*models.TODO, error) {
	query := `
  SELECT id,title,description,completed,created_at,updated_at
  FROM todos
  ORDER BY created_at DESC`

	rows, err := p.db.Query(query)

	if err != nil {
		return nil, errors.New("failed to get todos ")
	}
	defer rows.Close()
	var todos []*models.TODO
	for rows.Next() {
		var todo models.TODO
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}

		todos = append(todos, &todo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over todos: %w", err)
	}

	return todos, nil
}
func (p *PostgresStore) GetByID(id int) (*models.TODO, error) {
	query := `
  SELECT id,title,description,completed,created_at,updated_at
  FROM todos
  WHERE id = $1`
	var todo models.TODO

	err := p.db.QueryRow(query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("failed to get todos ")
	}

	return &todo, nil
}
func (p *PostgresStore) Delete(id int) error {
	query := `
  DELETE
  FROM todos
  WHERE id = $1`

	_, err := p.db.Exec(query, id)

	if err != nil {
		return errors.New("failed to delete todos ")
	}

	return nil
}
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
