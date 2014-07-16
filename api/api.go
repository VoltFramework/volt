package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/vieux/volt/mesoslib"
	"github.com/vieux/volt/mesosproto"
)

type API struct {
	frameworkInfo *mesosproto.FrameworkInfo
	m             *mesoslib.MesosLib
	log           *logrus.Logger

	tasks []Task
}

func NewAPI(m *mesoslib.MesosLib, frameworkInfo *mesosproto.FrameworkInfo, log *logrus.Logger) *API {
	return &API{
		frameworkInfo: frameworkInfo,
		log:           log,
		m:             m,
		tasks:         []Task{},
	}
}

// Simple _ping endpoint, returns OK
func (api *API) _ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

type Task struct {
	Command string  `json:"cmd"`
	Cpus    float64 `json:"cpus,string"`
	Mem     float64 `json:"mem,string"`
}

// Enpoint to call to add a new task
func (api *API) tasksAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		api.log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		if err := api.m.RequestOffer(api.frameworkInfo, task.Cpus, task.Mem); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	api.tasks = append(api.tasks, task)
	io.WriteString(w, "OK")
}

// Endpoint to list all the tasks
func (api *API) tasksList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Size  int    `json:"size"`
		Tasks []Task `json:"tasks"`
	}{
		len(api.tasks),
		api.tasks,
	}
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Register all the routes and then serve the API
func (api *API) ListenAndServe(port int) error {
	r := mux.NewRouter()
	api.log.WithFields(logrus.Fields{"port": port}).Info("Starting API...")

	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request){
		"GET": {
			"/_ping": api._ping,
			"/tasks": api.tasksList,
		},
		"POST": {
			"/tasks": api.tasksAdd,
		},
	}

	for method, routes := range endpoints {
		for route, fct := range routes {
			_route := route
			_fct := fct
			_method := method

			api.log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering API route...")
			r.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				api.log.WithFields(logrus.Fields{"from": r.RemoteAddr}).Infof("[%s] %s", _method, _route)
				_fct(w, r)
			})
		}
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
