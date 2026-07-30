package services

import (
	"go-todo/models"
	"go-todo/repositories"
)

type TodoService struct {
	repo repositories.TodoRepositoryInterface
}

func NewTodoService(repo repositories.TodoRepositoryInterface) *TodoService {
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

func (s *TodoService) GetTodoByID(id int) (*models.Todo, error) {
	return s.repo.GetTodoByID(id)
}

func (s *TodoService) DeleteTodo(id int) error {
	return s.repo.DeleteTodo(id)
}

func (s *TodoService) UpdateTodo(id int, todo *models.Todo) error {
	return s.repo.UpdateTodo(id, todo)
}
