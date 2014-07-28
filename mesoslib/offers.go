package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RequestOffer(cpus, mem float64) (*mesosproto.Offer, error) {
	m.Log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem}).Info("Requesting offers...")

	var event *mesosproto.Event

	select {
	case event = <-m.GetEvent(mesosproto.Event_OFFERS):
	}

	if event == nil {
		cpusName := "cpus"
		memoryName := "memory"
		scalar := mesosproto.Value_SCALAR

		if err := m.send(&mesosproto.ResourceRequestMessage{
			FrameworkId: m.frameworkInfo.Id,
			Requests: []*mesosproto.Request{
				&mesosproto.Request{
					Resources: []*mesosproto.Resource{
						&mesosproto.Resource{
							Name: &cpusName,
							Type: &scalar,
							Scalar: &mesosproto.Value_Scalar{
								Value: &cpus,
							},
						},
						&mesosproto.Resource{
							Name: &memoryName,
							Type: &scalar,
							Scalar: &mesosproto.Value_Scalar{
								Value: &mem,
							},
						},
					},
				},
			},
		}, "mesos.internal.ResourceRequestMessage"); err != nil {
			return nil, err
		}

		event = <-m.GetEvent(mesosproto.Event_OFFERS)
	}

	if len(event.Offers.Offers) > 0 {
		m.Log.WithFields(logrus.Fields{"Id": event.Offers.Offers[0].Id}).Info("Received offer.")
		return event.Offers.Offers[0], nil
	}
	return nil, nil
}
