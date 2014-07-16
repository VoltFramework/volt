package mesoslib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/goprotobuf/proto"

	"github.com/gorilla/mux"
	"github.com/vieux/volt/mesosproto"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/{scheduler}/mesos.internal.FrameworkRegisteredMessage", FrameworkRegisteredMessage)
	http.Handle("/", r)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatalf("failed to start listening on port %d", port)
		}
	}()
}

func FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	message := new(mesosproto.FrameworkRegisteredMessage)
	if proto.Unmarshal(data, message) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		eventType := mesosproto.Event_REGISTERED
		events <- &mesosproto.Event{
			Type: &eventType,
			Registered: &mesosproto.Event_Registered{
				FrameworkId: message.FrameworkId,
				MasterInfo:  message.MasterInfo,
			},
		}
		w.WriteHeader(http.StatusOK)
	}
}
