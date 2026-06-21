package handlers

import (
	"encoding/json"
	"go-todo/models"
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
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// w http.ResponseWriter = tsのres, r *http.Request = tsのreq
func (h *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	// jsonを返すという宣言
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)     //HTTP 200 OKを返すtsのres.status(200)
	json.NewEncoder(w).Encode(todos) //todosをjson化してレスポンスに書く res.json(todos)と同じ
}

func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo) //DecodeでjsonをGoの構造体にしてGoが読めるようにする
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateTodo(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
