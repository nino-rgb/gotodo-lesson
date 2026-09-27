package services

import (
	"go-todo/models"
	"testing"
)

// nodetsと同じserviceのテストで呼ぶためのmockrepoを作成
// インターフェイスの要件を満たすための12-26行
type MockTodoRepository struct {
	todos []models.Todo
	err   error
}

func (m *MockTodoRepository) GetTodos() ([]models.Todo, error) {
	return m.todos, m.err
}

func (m *MockTodoRepository) GetTodoByID(id int) (*models.Todo, error) {
	return nil, nil
}

func (m *MockTodoRepository) CreateTodo(todo *models.Todo) error {
	return nil
}

func (m *MockTodoRepository) UpdateTodo(id int, todo *models.Todo) error {
	return nil
}

func (m *MockTodoRepository) DeleteTodo(id int) error {
	return nil
}

func TestTodoService_GetTodos(t *testing.T) {
	todos := []models.Todo{
		{ID: 1, Title: "Todo1", Description: "Description1"},
		{ID: 2, Title: "Todo2", Description: "Description2"},
	}

	mockRepo := &MockTodoRepository{
		todos: todos,
		err:   nil,
	}

	service := NewTodoService(mockRepo)

	result, err := service.GetTodos()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 todos, got %d", len(result))
	}
}
