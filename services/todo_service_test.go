package services

import "go-todo/models"

type MockTodoRepository struct {
	todos []models.Todo
	err   error
}
