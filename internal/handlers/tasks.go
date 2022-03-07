package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fylerx/fyler/internal/errors"
	"github.com/fylerx/fyler/internal/tasks"
)

type Tasks interface {
	GetAll() ([]*tasks.Task, error)
	Create(task *tasks.Task) (*tasks.Task, error)
}

type TasksHandler struct {
	TasksRepo Tasks
}

func (h *TasksHandler) Index(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TasksRepo.GetAll()
	if err != nil {
		errors.JsonError(w, http.StatusInternalServerError, "server error")
		return
	}

	resp, _ := json.Marshal(tasks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	// currentSession, err := session.SessionFromContext(r.Context())
	// if err != nil {
	// 	errors.JsonError(w, http.StatusUnauthorized, "unauthorized user")
	// 	return
	// }

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	newTask := tasks.Task{}
	err := json.Unmarshal(body, &newTask)
	if err != nil {
		errors.JsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	newTask.ProjectID = 1 //currentSession.ProjectID

	task, err := h.TasksRepo.Create(&newTask)
	if err != nil {
		errors.JsonError(w, http.StatusBadRequest, "cant create task")
		return
	}

	resp, _ := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
