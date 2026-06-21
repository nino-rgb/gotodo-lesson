package repositories

import (
	"database/sql"
	"go-todo/models"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (t *TodoRepository) GetTodos() ([]models.Todo, error) {
	query := "SELECT id, title, description, created_at, updated_at FROM todos"
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *TodoRepository) CreateTodo(todo *models.Todo) error {
	query := "INSERT INTO todos (title, description) VALUES(?, ?)"

	result, err := t.db.Exec(query, todo.Title, todo.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	todo.ID = int(id)

	return nil
}
