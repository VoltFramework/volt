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
				&mesosproto.Request{
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

func (m *MesosLib) DeclineOffers(offers []*mesosproto.Offer) error {
	m.Log.Infof("Declining %d offers.", len(offers))
	var offerIds = []*mesosproto.OfferID{}

	for _, offer := range offers {
		offerIds = append(offerIds, offer.Id)
	}

	return m.send(&mesosproto.LaunchTasksMessage{
		FrameworkId: m.frameworkInfo.Id,
		OfferIds:    offerIds,
		Filters:     &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}
