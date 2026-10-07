package handlers

import (
	"encoding/json"
	"go-todo/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTodoService struct {
	GetTodoByIDFunc func(id int) (*models.Todo, error)
}

func (m *MockTodoService) GetTodos() ([]models.Todo, error) {
	panic("not implemented")
}

func (m *MockTodoService) GetTodoByID(id int) (*models.Todo, error) {
	return m.GetTodoByIDFunc(id)
}

func (m *MockTodoService) CreateTodo(todo *models.Todo) error {
	panic("not implemented")
}

func (m *MockTodoService) UpdateTodo(id int, todo *models.Todo) error {
	panic("not implemented")
}

func (m *MockTodoService) DeleteTodo(id int) error {
	panic("not implemented")
}

func TestTodoHandler_GetTodoByID(t *testing.T) {
	expectedTodo := &models.Todo{
		ID:          1,
		Title:       "テストTodo",
		Description: "テスト詳細",
	}

	mockService := &MockTodoService{
		GetTodoByIDFunc: func(id int) (*models.Todo, error) {
			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/todos?id=1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var actualTodo models.Todo

	err := json.NewDecoder(rec.Body).Decode(&actualTodo)
	if err != nil {
		t.Fatalf("failed To decode responce: %v", err)
	}

	if actualTodo.ID != expectedTodo.ID {
		t.Errorf("expected ID %d, got %d", expectedTodo.ID, actualTodo.ID)
	}

	if actualTodo.Title != expectedTodo.Title {
		t.Errorf("expected title %s, got %s", expectedTodo.Title, actualTodo.Title)
	}

	if actualTodo.Description != expectedTodo.Description {
		t.Errorf(
			"expected description %s, got %s",
			expectedTodo.Description,
			actualTodo.Description,
		)
	}
}
