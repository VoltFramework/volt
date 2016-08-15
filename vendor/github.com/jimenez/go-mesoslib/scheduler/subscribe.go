package scheduler

import (
	"encoding/json"
	"io"
	"log"

	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) handleEvents(body io.ReadCloser, handler OfferHandler) {
	dec := json.NewDecoder(body)
	for {
		var event schedulerproto.Event
		if err := dec.Decode(&event); err != nil || event.Type == nil {
			continue
		}
		if event.GetType() == schedulerproto.Event_UPDATE {
			taskStatus := event.GetUpdate().GetStatus()
			lib.tasks[taskStatus.GetTaskId().GetValue()] = taskStatus.GetAgentId()
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
				handler(offer)
			}
			log.Println("framework", lib.name, "received", len(event.GetOffers().GetOffers()), "offer(s)")
		}
	}
}

func (lib *SchedulerLib) Subscribe(handler OfferHandler) error {
	call := &schedulerproto.Call{
		Type: schedulerproto.Call_SUBSCRIBE.Enum(),
		Subscribe: &schedulerproto.Call_Subscribe{
			FrameworkInfo: lib.frameworkInfo,
		},
	}

	body, err := lib.send(call, 200)
	if err != nil {
		return err
	}
	go lib.handleEvents(body, handler)
	return nil
}
