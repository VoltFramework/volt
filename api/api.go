package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VoltFramework/volt/inmemory"
	"github.com/VoltFramework/volt/task"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/jimenez/mesoscon-demo/lib"
	"github.com/jimenez/mesoscon-demo/lib/mesosproto"
)

var defaultState = mesosproto.TaskState_TASK_STAGING

type API struct {
	handler  *mux.Router
	m        *lib.DemoLib
	registry Registry
	OffersCH chan *mesosproto.Offer
}

func (api *API) HandleOffers(offer *mesosproto.Offer) {
	api.OffersCH <- offer
}

// Simple _ping endpoint, returns OK
func (api *API) _ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func (api *API) writeError(w http.ResponseWriter, code int, message string) {
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

	var resources = lib.BuildResources(task.Cpus, task.Mem, task.Disk)

	api.m.LaunchTask(<-api.OffersCH, resources, &lib.Task{
		ID:      task.ID,
		Command: strings.Split(task.Command, " "),
		Image:   task.DockerImage,
		Volumes: task.Volumes,
	})

	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, "OK")
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
	io.WriteString(w, "OK")
}

// Endpoint to kill a task
func (api *API) tasksKill(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

// Register all the routes and then serve the API
func ListenAndServe(m *lib.DemoLib, port int) *API {
	api := &API{
		m:        m,
		registry: inmemory.New(),
		handler:  mux.NewRouter(),
		OffersCH: make(chan *mesosproto.Offer),
	}

	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request){
		"DELETE": {
			"/tasks/{id}": api.tasksDelete,
		},
		"GET": {
			"/_ping": api._ping,
			"/tasks": api.tasksList,
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

			api.handler.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				_fct(w, r)
			})
		}
	}
	api.handler.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "./static/"}))
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", 9999), api.handler)
	}()

	return api
}
