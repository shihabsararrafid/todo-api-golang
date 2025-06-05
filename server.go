package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todo-api/middlewares"
	"todo-api/models"
	"todo-api/storage"
	"todo-api/utils"
)

type TodoServer struct {
	store storage.Store
	http.Handler
}

func NewTodoServer(store storage.Store) *TodoServer {

	p := new(TodoServer)
	p.store = store

	router := http.NewServeMux()

	// You need to provide a valid handler function here, e.g., p.todoHandler
	router.Handle("/todo/", http.HandlerFunc(p.todoHandler))

	p.Handler = middlewares.Logging(router)

	return p
}

func (t *TodoServer) todoHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/todo")
	log.Print(path)
	switch {
	case path == "" || path == "/":

		switch r.Method {
		case http.MethodGet:
			t.GetTodos(w, r)
		case http.MethodPost:
			t.createTodo(w, r)
		default:
			utils.WriteError(w, http.StatusNotImplemented, "Method not allowed")
		}
	case strings.HasPrefix(path, "/"):
		idStr := strings.TrimPrefix(path, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid TODO ID")
			return
		}
		switch r.Method {
		case http.MethodGet:
			t.GetSingleTodo(w, r, id)
		case http.MethodPatch:
			t.updateTodo(w, r, id)
		case http.MethodDelete:
			t.DeleteTodo(w, r, id)
		default:
			utils.WriteError(w, http.StatusNotImplemented, "Method not allowed")

		}
	default:
		utils.WriteError(w, http.StatusNotImplemented, "Method not allowed")

	}

}
func (t *TodoServer) createTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal("Invalid Json", w, http.StatusBadRequest)
		return
	}
	todo, err := t.store.Create(req)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}
	utils.WriteJson(w, http.StatusCreated, todo)
}

func (t *TodoServer) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, _ := t.store.GetAll()
	utils.WriteJson(w, http.StatusOK, todos)
}
func (t *TodoServer) GetSingleTodo(w http.ResponseWriter, r *http.Request, id int) {
	todo, err := t.store.GetByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}
	utils.WriteJson(w, http.StatusOK, todo)
}
func (t *TodoServer) updateTodo(w http.ResponseWriter, r *http.Request, id int) {
	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal("Invalid Json", w, http.StatusBadRequest)
		return
	}
	todo, err := t.store.Update(id, req)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}
	utils.WriteJson(w, http.StatusCreated, todo)
}
func (t *TodoServer) DeleteTodo(w http.ResponseWriter, r *http.Request, id int) {
	err := t.store.Delete(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Todo not found")
		return
	}
	utils.WriteJson(w, http.StatusOK, "Todo Deleted")
}
