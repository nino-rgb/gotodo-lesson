package handlers

import (
	"encoding/json"
	"go-todo/models"
	"go-todo/services"
	"net/http"
	"strconv"
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
		id := r.URL.Query().Get("id")

		if id != "" {
			h.getTodoByID(w, r)
			return
		}

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

func (h *TodoHandler) getTodoByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr) //ASCII to Integer の略ASCII(文字列型)をInteger(整数型)に変換している
	//これがないとリポジトリ層でid intで要求しているのに文字列型できてて通らないよ‐ってことにならない
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
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
