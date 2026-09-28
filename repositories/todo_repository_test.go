package repositories

import (
	"go-todo/models"
	"go-todo/tests/utils/database"
	"testing"
)

func TestTodoRepository_CreateTodo(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	todo := &models.Todo{
		Title:       "Repo Test Todo",
		Description: "Repo Test Todo",
	}

	err = repo.CreateTodo(todo)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	//IDが0のままな訳ないから0ならエラーを返す
	if todo.ID == 0 {
		t.Errorf("expected new ID, got %d", todo.ID)
	}

	var savedTodo models.Todo

	err = db.QueryRow(
		"SELECT id, title, description FROM todos WHERE id = ?",
		todo.ID,
	).Scan(
		&savedTodo.ID,
		&savedTodo.Title,
		&savedTodo.Description,
	)

	if err != nil {
		t.Fatalf("failed to get created todo: %v", err)
	}
	//タイトルと詳細が正しく保存されたかのチェック
	if savedTodo.Title != todo.Title {
		t.Errorf("expected title %s, got %s", todo.Title, savedTodo.Title)
	}

	if savedTodo.Description != todo.Description {
		t.Errorf(
			"expected description %s, got %s",
			todo.Description,
			savedTodo.Description,
		)
	}
	//テスト終了時にテストに使ったデータだけを削除する
	defer func() {
		_, err := db.Exec("DELETE FROM todos WHERE id = ?", todo.ID)
		if err != nil {
			t.Errorf("failed to cleanup todo: %v", err)
		}
	}()
}
