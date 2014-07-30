package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/gorilla/mux"
)

type API struct {
	m   *mesoslib.MesosLib
	log *logrus.Logger

	tasks []*Task
}

func NewAPI(m *mesoslib.MesosLib) *API {
	return &API{
		m:     m,
		tasks: make([]*Task, 0),
	}
}

// Simple _ping endpoint, returns OK
func (api *API) _ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

var defaultState mesosproto.TaskState = mesosproto.TaskState_TASK_STAGING

type Task struct {
	ID      string  `json:"id"`
	Command string  `json:"cmd"`
	Cpus    float64 `json:"cpus,string"`
	Mem     float64 `json:"mem,string"`

	SlaveId *string               `json:"slave_id",string`
	State   *mesosproto.TaskState `json:"state,string"`
}

// Enpoint to call to add a new task
func (api *API) tasksAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.m.Log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var task = Task{State: &defaultState}
	err = json.Unmarshal(body, &task)
	if err != nil {
		api.m.Log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := make([]byte, 6)
	n, err := rand.Read(id)
	if n != len(id) || err != nil {
		api.m.Log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	task.ID = hex.EncodeToString(id)

	go func() {
		offer, err := api.m.RequestOffer(task.Cpus, task.Mem)
		if err != nil {
			api.m.Log.Warn(err)
		}
		if offer != nil {
			task.SlaveId = offer.SlaveId.Value
			api.m.LaunchTask(offer, task.Command, task.ID, &task.State)
		}
	}()

	api.tasks = append(api.tasks, &task)
	io.WriteString(w, "OK")
}

// Endpoint to list all the tasks
func (api *API) tasksList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Size  int     `json:"size"`
		Tasks []*Task `json:"tasks"`
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
	api.m.Log.WithFields(logrus.Fields{"port": port}).Info("Starting API...")

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

			api.m.Log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering API route...")
			r.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				api.m.Log.WithFields(logrus.Fields{"from": r.RemoteAddr}).Infof("[%s] %s", _method, _route)
				_fct(w, r)
			})
		}
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
