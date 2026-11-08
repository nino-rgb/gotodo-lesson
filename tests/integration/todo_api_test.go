package integration

import (
	"encoding/json"
	"fmt"
	"go-todo/handlers"
	"go-todo/models"
	"go-todo/repositories"
	"go-todo/services"
	"go-todo/tests/utils/database"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_GetTodos(t *testing.T) {
	todo := &models.Todo{
		Title:       "APIテスト用Todo",
		Description: "GET APIの統合テスト",
	}

	db, err := database.CreateDBConnection("../../.env.test")
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

func TestAPI_GetTodosByID(t *testing.T) {
	todo := &models.Todo{
		Title:       "APIテスト用Todo",
		Description: "GET APIの統合テスト",
	}

	db, err := database.CreateDBConnection("../../.env.test")
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

	err = repo.CreateTodo(todo)
	if err != nil {
		t.Fatalf("failed to prepare test todo: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/todos?id=%d", todo.ID), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var actualTodo models.Todo
	err = json.NewDecoder(rec.Body).Decode(&actualTodo)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if actualTodo.ID != todo.ID {
		t.Errorf("expected ID %d, got %d", todo.ID, actualTodo.ID)
	}

	if actualTodo.Title != todo.Title {
		t.Errorf("expected Title %s, got %s", todo.Title, actualTodo.Title)
	}

	if actualTodo.Description != todo.Description {
		t.Errorf("expected Description %s, got %s", todo.Description, actualTodo.Description)
	}
}

func TestAPI_GetTodoByID_NotFound(t *testing.T) {
	db, err := database.CreateDBConnection("../../.env.test")
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	repo := repositories.NewTodoRepository(db)
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/todos?id=999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestAPI_PostTodo(t *testing.T) {

	db, err := database.CreateDBConnection("../../.env.test")
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	var actualTodo models.Todo

	repo := repositories.NewTodoRepository(db)
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)

	body := strings.NewReader(`{"title":"POSTテスト用Todo","description":"POST APIの統合テスト"}`)

	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	err = json.NewDecoder(rec.Body).Decode(&actualTodo)
	if err != nil {
		t.Fatalf("failed to decode responce: %v", err)
	}
	if actualTodo.ID <= 0 {
		t.Fatalf("Todo ID == 0, cleanup skipped")
	}

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM todos WHERE id = ?", actualTodo.ID)
		if err != nil {
			t.Errorf("failed to clean up test todo: %v", err)
		}
	})

	if actualTodo.Title != "POSTテスト用Todo" {
		t.Errorf("expected title %s, got %s", "POSTテスト用Todo", actualTodo.Title)
	}

	if actualTodo.Description != "POST APIの統合テスト" {
		t.Errorf("expected description %s, got %s", "POST APIの統合テスト", actualTodo.Description)
	}

	dbTodo, err := repo.GetTodoByID(actualTodo.ID)

	if err != nil {
		t.Fatalf("failed to get todo from DB: %v", err)
	}
	if dbTodo.ID != actualTodo.ID {
		t.Errorf("expected ID %d, got %d", actualTodo.ID, dbTodo.ID)
	}

	if dbTodo.Title != actualTodo.Title {
		t.Errorf("expected Title %s, got %s", actualTodo.Title, dbTodo.Title)
	}

	if dbTodo.Description != actualTodo.Description {
		t.Errorf("expected Description %s, got %s", actualTodo.Description, dbTodo.Description)
	}
}

func TestAPI_PutTodo(t *testing.T) {

	db, err := database.CreateDBConnection("../../.env.test")
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	repo := repositories.NewTodoRepository(db)
	service := services.NewTodoService(repo)
	handler := handlers.NewTodoHandler(service)

	todo := &models.Todo{
		Title:       "更新前のタイトル",
		Description: "更新前の詳細",
	}

	err = repo.CreateTodo(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM todos WHERE id = ?", todo.ID)
		if err != nil {
			t.Errorf("failed to clean up test todo: %v", err)
		}
	})

	body := strings.NewReader(`{"title":"PUTテスト用Todo","description":"PUT APIの統合テスト"}`)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/todos?id=%d", todo.ID), body)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	dbTodo, err := repo.GetTodoByID(todo.ID)
	if err != nil {
		t.Fatalf("failed to get todo from DB: %v", err)
	}

	if dbTodo.Title != "PUTテスト用Todo" {
		t.Errorf("expected title %s, got %s", "PUTテスト用Todo", dbTodo.Title)
	}

	if dbTodo.Description != "PUT APIの統合テスト" {
		t.Errorf("expected description %s, got %s", "PUT APIの統合テスト", dbTodo.Description)
	}
}
