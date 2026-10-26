package handlers

import (
	"encoding/json"
	"go-todo/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockTodoService struct {
	GetTodoByIDFunc func(id int) (*models.Todo, error)
	UpdateTodoFunc  func(id int, todo *models.Todo) error
	CreateTodoFunc  func(todo *models.Todo) error
}

func (m *MockTodoService) GetTodos() ([]models.Todo, error) {
	panic("not implemented")
}

func (m *MockTodoService) GetTodoByID(id int) (*models.Todo, error) {
	return m.GetTodoByIDFunc(id)
}

func (m *MockTodoService) CreateTodo(todo *models.Todo) error {
	return m.CreateTodoFunc(todo)
}

func (m *MockTodoService) UpdateTodo(id int, todo *models.Todo) error {
	return m.UpdateTodoFunc(id, todo)
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

func TestTodoHandler_UpdateTodo_NotFound(t *testing.T) {

	mockService := &MockTodoService{
		UpdateTodoFunc: func(id int, todo *models.Todo) error {
			return models.ErrNotFound
		},
	}

	handler := NewTodoHandler(mockService)

	body := `{
    "title": "更新タイトル",
		"description": "更新詳細"
		}`
	req := httptest.NewRequest(http.MethodPut, "/todos?id=999", strings.NewReader(body))
	rec := httptest.NewRecorder()
	//strings.NewReader(body) 文字列を httpリクエストのbodyとして読み取れる形に変更するもの
	//httptest.NewRequest()の第三引数はio.Readerを返さなければならないために使っている
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandler_CreateTodo_InvalidJSON(t *testing.T) {

	mockService := &MockTodoService{
		CreateTodoFunc: func(todo *models.Todo) error {
			t.Fatalf("CreateTodo Should not be called")
			return nil
		},
	}

	handler := NewTodoHandler(mockService)

	body := `{"title"`

	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
