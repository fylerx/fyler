package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fylerx/fyler/internal/projects"
	"github.com/fylerx/fyler/internal/tasks"
	u "github.com/fylerx/fyler/pkg/utils"
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
		u.RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	u.RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	currentProject, err := projects.ProjectFromContext(r.Context())
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "can't fetch current project")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	taskInput := tasks.Task{}
	err = json.Unmarshal(body, &taskInput)
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "can't unpack payload")
		return
	}

	taskInput.ProjectID = currentProject.ID

	task, err := h.TasksRepo.Create(&taskInput)
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "can't create task")
		return
	}

	u.RespondWithJSON(w, http.StatusOK, task)
}
