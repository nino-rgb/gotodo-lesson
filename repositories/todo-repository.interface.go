package repositories

import "go-todo/models"

type TodoRepositoryInterface interface {
	GetTodos() ([]models.Todo, error)
	GetTodoByID(id int) (*models.Todo, error)
	CreateTodo(todo *models.Todo) error
	UpdateTodo(id int, todo *models.Todo) error
	DeleteTodo(id int) error
}
