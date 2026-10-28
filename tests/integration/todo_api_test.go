package integration

import (
	"encoding/json"
	"go-todo/handlers"
	"go-todo/models"
	"go-todo/repositories"
	"go-todo/services"
	"go-todo/tests/utils/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_GetTodos(t *testing.T) {
	todo := &models.Todo{
		Title:       "APIテスト用Todo",
		Description: "GET APIの統合テスト",
	}

	db, err := database.CreateDBConnection("../../.env")
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	t.Cleanup(func() {

		_, err := db.Exec("DELETE FROM todos WHERE id = ?", todo.ID)
		if err != nil {
			t.Errorf("failed to clean up test todo: %v", err)
		}
		db.Close()
	})
	repo := repositories.NewTodoRepository(db)
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()

	err = repo.CreateTodo(todo)
	if err != nil {
		t.Fatalf("failed to prepare test todo: %v", err)
	}

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var actualTodos []models.Todo

	err = json.NewDecoder(rec.Body).Decode(&actualTodos)
	if err != nil {
		t.Fatalf("failed to decode responce: %v", err)
	}
	t.Logf("取得したTodoの件数: %d", len(actualTodos))

	for _, todo := range actualTodos {
		t.Logf("ID: %d, Title: %s", todo.ID, todo.Title)
	}

	found := false
	for _, actualTodo := range actualTodos {
		if actualTodo.ID == todo.ID {
			found = true

			if actualTodo.Title != todo.Title {
				t.Errorf("expect title %s, got %s", todo.Title, actualTodo.Title)
			}

			if actualTodo.Description != todo.Description {
				t.Errorf("expect description %s, got %s", todo.Description, actualTodo.Description)
			}
			break
		}
	}

	if !found {
		t.Errorf("Todo with ID %d was not found", todo.ID)
	}
}
