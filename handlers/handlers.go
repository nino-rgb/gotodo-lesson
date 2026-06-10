package handlers

import (
	"encoding/json"
	"go-todo/services"
	"net/http"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// ここでどのメソッドかを判別tsでいうrouter.get()をまとめてしている
func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTodos(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
