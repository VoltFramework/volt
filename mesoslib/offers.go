package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RequestOffer(frameworkInfo *mesosproto.FrameworkInfo, cpus, mem float64) (*mesosproto.Offer, error) {
	m.log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem}).Info("Requesting offers...")

	var event *mesosproto.Event

	select {
	case event = <-m.events:
		if *event.Type != mesosproto.Event_OFFERS {
			event = nil
		}
	}

	if event == nil {
		callType := mesosproto.Call_REQUEST
		cpusName := "cpus"
		memoryName := "memory"
		scalar := mesosproto.Value_SCALAR

		requestCall := &mesosproto.Call{
			FrameworkInfo: frameworkInfo,
			Type:          &callType,
			Request: &mesosproto.Call_Request{
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
			},
		}

		if err := m.send(requestCall, "mesos.internal.ResourceRequestMessage"); err != nil {
			return nil, err
		}

		event = <-m.events
	}

	if len(event.Offers.Offers) > 0 {
		m.log.WithFields(logrus.Fields{"Id": event.Offers.Offers[0].Id}).Info("Received offer.")
		return event.Offers.Offers[0], nil
	}
	return nil, nil
}
