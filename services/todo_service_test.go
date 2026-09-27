package services

import (
	"errors"
	"go-todo/models"
	"testing"
)

// nodetsと同じserviceのテストで呼ぶためのmockrepoを作成
// インターフェイスの要件を満たすための12-26行
type MockTodoRepository struct {
	todos []models.Todo
	todo  *models.Todo
	err   error
}

func (m *MockTodoRepository) GetTodos() ([]models.Todo, error) {
	return m.todos, m.err
}

func (m *MockTodoRepository) GetTodoByID(id int) (*models.Todo, error) {
	return m.todo, m.err
}

func (m *MockTodoRepository) CreateTodo(todo *models.Todo) error {
	return nil
}

func (m *MockTodoRepository) UpdateTodo(id int, todo *models.Todo) error {
	return m.err
}

func (m *MockTodoRepository) DeleteTodo(id int) error {
	return nil
}

func TestTodoService_GetTodos(t *testing.T) {
	todos := []models.Todo{
		{ID: 1, Title: "Todo1", Description: "Description1"},
		{ID: 2, Title: "Todo2", Description: "Description2"},
	}

	//MockTodoRepositoryにtodosには38のtodosを
	//errはnilを返すように
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

func TestTodoService_GetTodoByID(t *testing.T) {
	todo := &models.Todo{
		ID:          1,
		Title:       "todo1",
		Description: "description1",
	}

	mockRepo := &MockTodoRepository{
		todo: todo,
		err:  nil,
	}

	service := NewTodoService(mockRepo)

	result, err := service.GetTodoByID(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected todo ID 1, got %d", result.ID)
	}
}

func TestTodoService_GetTodoByID_Error(t *testing.T) {
	expectedErr := errors.New("repository error")

	//MockTodoRepositoryにtodoはnilを
	//err は89で作成したexpectedErrを返す
	mockRepo := &MockTodoRepository{
		todo: nil,
		err:  expectedErr,
	}

	service := NewTodoService(mockRepo)

	result, err := service.GetTodoByID(1)

	//87で作成したエラーを返す
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil todo, got %v", result)
	}
}

func TestTodoService_CreateTodo(t *testing.T) {
	todo := &models.Todo{
		Title:       "todo1",
		Description: "description1",
	}

	//MockTodoRepositoryには何も返さない 25でのreturn nilが返る(エラーなし作成成功!)
	mockRepo := &MockTodoRepository{}

	service := NewTodoService(mockRepo)

	err := service.CreateTodo(todo)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTodoService_UpdateTodo(t *testing.T) {
	todo := &models.Todo{
		Title:       "Updated Todo",
		Description: "Updated Description",
	}

	mockRepo := &MockTodoRepository{}

	service := NewTodoService(mockRepo)

	err := service.UpdateTodo(1, todo)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTodoService_UpdateTodo_Error(t *testing.T) {
	expectedErr := errors.New("repository error")

	todo := &models.Todo{
		Title:       "Updated Todo",
		Description: "Updated Description",
	}

	mockRepo := &MockTodoRepository{
		err: expectedErr,
	}

	service := NewTodoService(mockRepo)

	err := service.UpdateTodo(1, todo)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestTodoService_DeleteTodo(t *testing.T) {
	mockRepo := &MockTodoRepository{}

	service := NewTodoService(mockRepo)

	err := service.DeleteTodo(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
