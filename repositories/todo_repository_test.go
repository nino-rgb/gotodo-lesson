package repositories

import (
	"database/sql"
	"go-todo/models"
	"go-todo/tests/utils/database"
	"testing"
)

// t.Helper() テスト本体ではなくテスト補助のためのヘルパー関数とGoに伝えるもの
func cleanTodos(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec("DELETE FROM todos")
	if err != nil {
		t.Fatalf("failed to clean todos: %v", err)
	}
}

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

func TestTodoRepository_CreateTodo_Error(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	repo := NewTodoRepository(db)

	//nodetsの時とは違いここでdb接続を切ってエラーを発生させる
	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close database: %v", err)
	}

	todo := &models.Todo{
		Title:       "Repo Test Todo",
		Description: "Repo Test Todo",
	}

	//db接続切ってるためdb.Exec()が失敗する
	err = repo.CreateTodo(todo)
	//エラーならテスト成功!
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
}

func TestTodoRepository_GetTodos(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	//一番上にtodosを空にするのを作成しているのでここで呼ぶだけでいい
	//nodetsのBeforEachを擬似的に再現している
	cleanTodos(t, db)

	//createの実装に依存したテストにならないためにrepo.CreateTodo(...)
	//ではなくdb.Exxec()を使ってインサートする
	_, err = db.Exec(
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		"ダミー1",
		"詳細1",
	)
	if err != nil {
		t.Fatalf("failed to insert todo: %v", err)
	}

	_, err = db.Exec(
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		"ダミー2",
		"詳細2",
	)
	if err != nil {
		t.Fatalf("failed to insert todo: %v", err)
	}

	repo := NewTodoRepository(db)

	todos, err := repo.GetTodos()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	//todoがあるかの確認用マップ
	titles := map[string]bool{}

	//SQLにORDER BYつけるんだるかったから作成順を気にしない仕様
	for _, todo := range todos {
		titles[todo.Title] = true
	}

	//ダミー1が取得できてなかったらテスト失敗
	if !titles["ダミー1"] {
		t.Errorf("expected ダミー1 to be included")
	}
	//ダミー2が取得できてなかったらテスト失敗
	if !titles["ダミー2"] {
		t.Errorf("expected ダミー2 to be included")
	}
}

func TestTodoRepository_GetTodos_Error(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	repo := NewTodoRepository(db)

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close database: %v", err)
	}

	_, err = repo.GetTodos()

	if err == nil {
		t.Fatalf("expect error, got nil")
	}
}

func TestTodoRepository_GetTodoByID(t *testing.T) {
	db, err := database.CreateDBConnection()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	defer db.Close()
	cleanTodos(t, db)

	result, err := db.Exec(
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		"ダミー1",
		"詳細1",
	)

	if err != nil {
		t.Fatalf("failed to insert todo: %v", err)
	}

	//ここで最後に生成されたidを取得して､resultに渡してる
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get inserted id: %v", err)
	}

	repo := NewTodoRepository(db)

	todo, err := repo.GetTodoByID(int(id))

	if todo.ID != int(id) {
		t.Errorf("expected id %d, got %d", id, todo.ID)
	}

	if todo.Title != "ダミー1" {
		t.Errorf("expected title ダミー1､ got %s", todo.Title)
	}

	if todo.Description != "詳細1" {
		t.Errorf("expected description 詳細1, got %s", todo.Description)
	}
}
