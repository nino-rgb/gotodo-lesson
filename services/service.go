package services

import (
	"go-todo/models"
	"go-todo/repositories"
)

type TodoService struct {
	repo *repositories.TodoRepository
}

func NewTodoService(repo *repositories.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) GetTodos() ([]models.Todo, error) {
	return s.repo.GetTodos()
}

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.repo.CreateTodo(todo)
}
