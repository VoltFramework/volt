package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/VoltFramework/volt/registry"
	"github.com/VoltFramework/volt/registry/inmemory"
	"github.com/VoltFramework/volt/registry/zookeeper"
	"github.com/VoltFramework/volt/task"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

var defaultState = mesosproto.TaskState_TASK_STAGING

type API struct {
	m        *mesoslib.MesosLib
	registry registry.Registry
}

// Simple _ping endpoint, returns OK
func (api *API) _ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
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
	task := &task.Task{State: &defaultState}

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

	if err := api.registry.Register(task.ID, task); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	f := func() error {
		var resources = api.m.BuildResources(task.Cpus, task.Mem, task.Disk)
		offers, err := api.m.RequestOffers(resources)
		if err != nil {
			return err
		}

		if len(offers) > 0 {
			task.SlaveId = *offers[0].SlaveId.Value
			task.SlaveHostname, err = api.m.GetSlaveHostname(task.SlaveId)
			if err != nil {
				api.m.Log.Warnf("Error getting slave hostname: %v", err)
			}

			if err := api.registry.Update(task.ID, task); err != nil {
				return err
			}

			return api.m.LaunchTask(offers[0], resources, &mesoslib.Task{
				ID:      task.ID,
				Command: strings.Split(task.Command, " "),
				Image:   task.DockerImage,
				Volumes: task.Volumes,
			})
		}

		return fmt.Errorf("No offers available")
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
	tasks, err := api.registry.Tasks()
	if err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Size  int          `json:"size"`
		Tasks []*task.Task `json:"tasks"`
	}{
		len(tasks),
		tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
	}
}

// Endpoint to delete a task
func (api *API) tasksDelete(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		id   = vars["id"]
	)

	if err := api.m.KillTask(id); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := api.registry.Delete(id); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

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

func (api *API) metrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := api.m.Metrics()
	if err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		api.writeError(w, http.StatusInternalServerError, err.Error())
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
	for event := range api.m.GetEvent(mesosproto.Event_UPDATE) {
		ID := event.Update.Status.TaskId.GetValue()

		state := event.Update.Status.State

		task, err := api.registry.Fetch(ID)
		if err != nil {
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage(), "error": err}).Error("Fetch task in registry")
			continue
		}

		task.State = state
		if err := api.registry.Update(ID, task); err != nil {
			api.m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage(), "error": err}).Error("Update task state in registry")
			continue
		}

		switch *state {
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
func ListenAndServe(m *mesoslib.MesosLib, port int, zk string) {
	api := &API{
		m: m,
	}
	if zk == "" {
		api.registry = inmemory.New()
	} else {
		api.registry = zookeeper.New(zk)
	}

	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request){
		"DELETE": {
			"/tasks/{id}": api.tasksDelete,
		},
		"GET": {
			"/_ping":                  api._ping,
			"/tasks/{id}/file/{file}": api.getFile,
			"/tasks":                  api.tasksList,
			"/metrics":                api.metrics,
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

			m.Log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering Volt-API route...")
			m.Router.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				m.Log.WithFields(logrus.Fields{"from": r.RemoteAddr}).Infof("[%s] %s", _method, _route)
				_fct(w, r)
			})
		}
	}
	m.Router.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "./static/"}))
	go api.handleStates()
	m.Log.WithFields(logrus.Fields{"port": port}).Info("Starting API...")
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), m.Router); err != nil {
			m.Log.Fatal(err)
		}
	}()
}
