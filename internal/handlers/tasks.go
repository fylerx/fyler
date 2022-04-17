package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/fylerx/fyler/internal/enum"
	"github.com/fylerx/fyler/internal/projects"
	"github.com/fylerx/fyler/internal/tasks"
	u "github.com/fylerx/fyler/pkg/utils"
)

type Tasks interface {
	GetAll() ([]*tasks.Task, error)
	Create(task *tasks.Task) (*tasks.Task, error)
}

type Job struct {
	ID        uint64        `json:"id"`
	ProjectID uint32        `json:"project_id"`
	Status    enum.Status   `json:"status"`
	TaskType  enum.TaskType `json:"task_type"`
	FilePath  string        `json:"file_path"`
}
type TasksHandler struct {
	TasksRepo Tasks
	JM        *faktory.Client
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

	newJob := Job{
		ID:        task.ID,
		ProjectID: task.ProjectID,
		Status:    task.Status,
		TaskType:  task.TaskType,
		FilePath:  task.FilePath,
	}

	response, _ := json.Marshal(newJob)
	job := faktory.NewJob(newJob.TaskType.String(), response)
	job.Queue = "medium"

	if err = h.JM.Push(job); err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	fmt.Printf("Job pushed %v", job.Jid)

	u.RespondWithJSON(w, http.StatusOK, task)
}
