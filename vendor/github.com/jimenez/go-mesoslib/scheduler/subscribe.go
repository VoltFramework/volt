package scheduler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) handleEvents(body io.ReadCloser, offerHandler OfferHandler, taskStatusHandler TaskStatusHandler) {
	dec := json.NewDecoder(body)
	for {
		var event schedulerproto.Event
		if err := dec.Decode(&event); err != nil || event.Type == nil {
			continue
		}
		if event.GetType() == schedulerproto.Event_UPDATE {
			taskStatus := event.GetUpdate().GetStatus()
			lib.tasks[taskStatus.GetTaskId().GetValue()] = taskStatus.GetAgentId()
			go taskStatusHandler(taskStatus)
			log.Println("Status for", taskStatus.GetTaskId().GetValue(), "on", taskStatus.GetAgentId().GetValue(), "is", taskStatus.GetState().String())
			if taskStatus.GetUuid() != nil {
				lib.Acknowledge(taskStatus.GetTaskId(), taskStatus.GetAgentId(), taskStatus.GetUuid())
			}
		}

		switch event.GetType() {
		case schedulerproto.Event_SUBSCRIBED:
			lib.frameworkID = event.GetSubscribed().GetFrameworkId()
			log.Println("framework", lib.name, "subscribed succesfully (", lib.frameworkID.String(), ")")
		case schedulerproto.Event_OFFERS:
			for _, offer := range event.GetOffers().GetOffers() {
				go offerHandler(offer)
			}
			log.Println("framework", lib.name, "received", len(event.GetOffers().GetOffers()), "offer(s)")
		}
	}
}

func (lib *SchedulerLib) Subscribe(offerHandler OfferHandler, taskStatusHandler TaskStatusHandler) error {
	call := &schedulerproto.Call{
		Type: schedulerproto.Call_SUBSCRIBE.Enum(),
		Subscribe: &schedulerproto.Call_Subscribe{
			FrameworkInfo: lib.frameworkInfo,
		},
	}
	f := func(r *http.Response) {
		lib.MesosStreamId = r.Header.Get("Mesos-Stream-Id")
	}
	body, err := lib.sendDetail(call, 200, f)
	if err != nil {
		return err
	}
	go lib.handleEvents(body, offerHandler, taskStatusHandler)
	return nil
}
