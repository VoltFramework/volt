package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

type API struct {
	sync.RWMutex

	m   *mesoslib.MesosLib
	log *logrus.Logger

	tasks  []*Task
	states map[string]*mesosproto.TaskState
}

func NewAPI(m *mesoslib.MesosLib) *API {
	return &API{
		m:      m,
		tasks:  make([]*Task, 0),
		states: make(map[string]*mesosproto.TaskState, 0),
	}
}

// Simple _ping endpoint, returns OK
func (api *API) _ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

type Task struct {
	ID      string   `json:"id"`
	Command string   `json:"cmd"`
	Cpus    float64  `json:"cpus,string"`
	Disk    float64  `json:"disk,string"`
	Mem     float64  `json:"mem,string"`
	Files   []string `json:"files"`

	SlaveId *string               `json:"slave_id,string"`
	State   *mesosproto.TaskState `json:"state,string"`
}

func (api *API) writeError(w http.ResponseWriter, code int, message string) {
	api.m.Log.Warn(message)
	w.WriteHeader(code)
	data := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		code,
		message,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
	}
}

// Enpoint to call to add a new task
func (api *API) tasksAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		defaultState mesosproto.TaskState = mesosproto.TaskState_TASK_STAGING
		task                              = Task{State: &defaultState}
	)

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		api.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id := make([]byte, 6)
	n, err := rand.Read(id)
	if n != len(id) || err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.ID = hex.EncodeToString(id)
	api.Lock()
	api.tasks = append(api.tasks, &task)
	api.states[task.ID] = task.State
	api.Unlock()

	f := func() error {
		offer, resources, err := api.m.RequestOffer(task.Cpus, task.Mem, task.Disk)
		if err != nil {
			return err
		}
		if offer != nil {
			task.SlaveId = offer.SlaveId.Value
			return api.m.LaunchTask(offer, resources, task.Command+" > volt_stdout 2> volt_stderr", task.ID)
		}
		return fmt.Errorf("No offer available")
	}
	if len(task.Files) > 0 {
		if err := f(); err != nil {
			api.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		files, err := api.m.ReadFile(task.ID, task.Files...)
		if err != nil {
			api.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(files); err != nil {
			api.writeError(w, http.StatusInternalServerError, err.Error())
		}
	} else {
		go f()
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, "OK")
	}
}

// Endpoint to list all the tasks
func (api *API) tasksList(w http.ResponseWriter, r *http.Request) {
	api.RLock()
	data := struct {
		Size  int     `json:"size"`
		Tasks []*Task `json:"tasks"`
	}{
		len(api.tasks),
		api.tasks,
	}
	api.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
	}
}

// Endpoint to delete a task
func (api *API) tasksDelete(w http.ResponseWriter, r *http.Request) {
	var (
		vars   = mux.Vars(r)
		id     = vars["id"]
		tasks  = make([]*Task, 0)
		states = make(map[string]*mesosproto.TaskState, len(api.states)-1)
	)

	if err := api.m.KillTask(id); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Lock()
	for _, task := range api.tasks {
		if task != nil && task.ID != id {
			tasks = append(tasks, task)
			states[task.ID] = task.State
		}
	}
	api.tasks = tasks
	api.states = states
	api.Unlock()
	io.WriteString(w, "OK")
}

// Endpoint to kill a task
func (api *API) tasksKill(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		id   = vars["id"]
	)
	if err := api.m.KillTask(id); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
	} else {
		io.WriteString(w, "OK")
	}
}

func (api *API) getFile(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		id   = vars["id"]
		file = vars["file"]
	)

	files, err := api.m.ReadFile(id, []string{file}...)
	if err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	content, ok := files[file]
	if !ok {
		api.writeError(w, http.StatusNotFound, file+" not found")
		return
	}
	io.WriteString(w, content)
}

func (api *API) handleStates() {
	for {
		event := <-api.m.GetEvent(mesosproto.Event_UPDATE)
		ID := event.Update.Status.TaskId.GetValue()

		state, ok := api.states[ID]
		if !ok {
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Update received for unknown task.")
			continue
		}

		*state = *event.Update.Status.State
		switch *event.Update.Status.State {
		case mesosproto.TaskState_TASK_STAGING:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task was registered.")
		case mesosproto.TaskState_TASK_STARTING:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is starting.")
		case mesosproto.TaskState_TASK_RUNNING:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is running.")
		case mesosproto.TaskState_TASK_FINISHED:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is finished.")
		case mesosproto.TaskState_TASK_FAILED:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task has failed.")
		case mesosproto.TaskState_TASK_KILLED:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task was killed.")
		case mesosproto.TaskState_TASK_LOST:
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task was lost.")
		}
	}
}

// Register all the routes and then serve the API
func (api *API) ListenAndServe(port int) error {
	r := mux.NewRouter()
	api.m.Log.WithFields(logrus.Fields{"port": port}).Info("Starting API...")

	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request){
		"DELETE": {
			"/tasks/{id}": api.tasksDelete,
		},
		"GET": {
			"/_ping":                  api._ping,
			"/tasks/{id}/file/{file}": api.getFile,
			"/tasks":                  api.tasksList,
		},
		"POST": {
			"/tasks": api.tasksAdd,
		},
		"PUT": {
			"/tasks/{id}/kill": api.tasksKill,
		},
	}

	for method, routes := range endpoints {
		for route, fct := range routes {
			_route := route
			_fct := fct
			_method := method

			api.m.Log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering API route...")
			r.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				api.m.Log.WithFields(logrus.Fields{"from": r.RemoteAddr}).Infof("[%s] %s", _method, _route)
				_fct(w, r)
			})
		}
	}
	r.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "./static/"}))
	go api.handleStates()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
