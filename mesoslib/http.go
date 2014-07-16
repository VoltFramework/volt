package mesoslib

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"code.google.com/p/goprotobuf/proto"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/vieux/volt/mesosproto"
)

func (m *MesosLib) initAPI() {
	r := mux.NewRouter()

	m.log.WithFields(logrus.Fields{"port": m.port}).Debug("Starting MesosLib-API...")
	endpoints := map[string]map[string]func(w http.ResponseWriter, r *http.Request, data []byte) error{
		"POST": {
			"/{scheduler}/mesos.internal.FrameworkRegisteredMessage": m.FrameworkRegisteredMessage,
			"/{scheduler}/mesos.internal.ResourceOffersMessage":      m.ResourceOffersMessage,
		},
	}

	for method, routes := range endpoints {
		for route, fct := range routes {
			_route := route
			_fct := fct
			_method := method

			m.log.WithFields(logrus.Fields{"method": _method, "route": _route}).Debug("Registering MesosLib-API route...")
			r.Path(_route).Methods(_method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				m.log.WithFields(logrus.Fields{"from": r.RemoteAddr, "scheduler": mux.Vars(r)["scheduler"]}).Debugf("[%s] %s", _method, _route)

				// extract request body
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					m.log.Warn(err)
					return
				}
				r.Body.Close()

				if err := _fct(w, r, data); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					m.log.Warn(err)
				}
			})
		}
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), r); err != nil {
			m.log.Fatalf("failed to start listening on port %d", m.port)
		}
	}()
}

// Endpoint called by the master upon registration
func (m *MesosLib) FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request, data []byte) error {
	message := new(mesosproto.FrameworkRegisteredMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}
	eventType := mesosproto.Event_REGISTERED
	m.events <- &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	}
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
	m.events <- &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
