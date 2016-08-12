package mesoslib

import "github.com/VoltFramework/volt/mesosproto"

func (m *MesosLib) RequestOffers(resources []*mesosproto.Resource) ([]*mesosproto.Offer, error) {
	m.Log.Info("Requesting offers...")

	var event *mesosproto.Event

	select {
	case event = <-m.GetEvent(mesosproto.Event_OFFERS):
	}

	if event == nil {
		if err := m.send(&mesosproto.ResourceRequestMessage{
			FrameworkId: m.frameworkInfo.Id,
			Requests: []*mesosproto.Request{
				{
					Resources: resources,
				},
			},
		}, "mesos.internal.ResourceRequestMessage"); err != nil {
			return nil, err
		}

		event = <-m.GetEvent(mesosproto.Event_OFFERS)
	}

	m.Log.Infof("Received %d offer(s).", len(event.Offers.Offers))
	return event.Offers.Offers, nil
}
