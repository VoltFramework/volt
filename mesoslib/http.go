package mesoslib

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"code.google.com/p/goprotobuf/proto"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/gorilla/mux"
)

func (m *MesosLib) initAPI() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Log.WithFields(logrus.Fields{"from": r.RemoteAddr}).Warnf("[%s] %s: Not implemented", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
	})
	m.Log.WithFields(logrus.Fields{"port": m.port}).Debug("Starting MesosLib-API...")
	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request, data []byte) error{
		"POST": {
			"/{scheduler}/mesos.internal.FrameworkRegisteredMessage": m.FrameworkRegisteredMessage,
			"/{scheduler}/mesos.internal.ResourceOffersMessage":      m.ResourceOffersMessage,
			"/{scheduler}/mesos.internal.StatusUpdateMessage":        m.StatusUpdateMessage,
		},
	}

	for method, routes := range endpoints {
		for route, fct := range routes {
			_route := route
			_fct := fct
			_method := method

			m.Log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering MesosLib-API route...")
			r.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				m.Log.WithFields(logrus.Fields{"from": r.RemoteAddr, "scheduler": mux.Vars(r)["scheduler"]}).Debugf("[%s] %s", _method, _route)

				// extract request body
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					m.Log.Warn(err)
					return
				}
				r.Body.Close()

				if err := _fct(w, r, data); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					m.Log.Warn(err)
				}
			})
		}
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), r); err != nil {
			m.Log.Fatalf("failed to start listening on port %d", m.port)
		}
	}()
}

// Endpoint called by the master upon registration
func (m *MesosLib) FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request, data []byte) error {
	message := new(mesosproto.FrameworkRegisteredMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	m.frameworkInfo.Id = message.FrameworkId

	eventType := mesosproto.Event_REGISTERED
	m.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

// Endpoint called by the master upon new offers
func (m *MesosLib) ResourceOffersMessage(w http.ResponseWriter, r *http.Request, data []byte) error {
	message := new(mesosproto.ResourceOffersMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}
	eventType := mesosproto.Event_OFFERS
	m.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

// Endpoint called by the master upon status update
func (m *MesosLib) StatusUpdateMessage(w http.ResponseWriter, r *http.Request, data []byte) error {
	message := new(mesosproto.StatusUpdateMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	if err := m.send(&mesosproto.StatusUpdateAcknowledgementMessage{
		FrameworkId: m.frameworkInfo.Id,
		SlaveId:     message.Update.Status.SlaveId,
		TaskId:      message.Update.Status.TaskId,
		Uuid:        message.Update.Uuid,
	}, "mesos.internal.StatusUpdateAcknowledgementMessage"); err != nil {
		return err
	}

	eventType := mesosproto.Event_UPDATE
	m.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Update: &mesosproto.Event_Update{
			Uuid:   message.Update.Uuid,
			Status: message.Update.Status,
		},
	})

	w.WriteHeader(http.StatusOK)
	return nil
}
