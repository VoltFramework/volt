package mesoslib

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RequestOffer(cpus, mem float64) (*mesosproto.Offer, []*mesosproto.Resource, error) {
	m.Log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem}).Info("Requesting offers...")

	var (
		resources = []*mesosproto.Resource{
			&mesosproto.Resource{
				Name:   proto.String("cpus"),
				Type:   mesosproto.Value_SCALAR.Enum(),
				Scalar: &mesosproto.Value_Scalar{Value: &cpus},
			},
			&mesosproto.Resource{
				Name:   proto.String("mem"),
				Type:   mesosproto.Value_SCALAR.Enum(),
				Scalar: &mesosproto.Value_Scalar{Value: &mem},
			},
		}
		event *mesosproto.Event
	)

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
			return nil, nil, err
		}

		event = <-m.GetEvent(mesosproto.Event_OFFERS)
	}

	if len(event.Offers.Offers) > 0 {
		m.Log.WithFields(logrus.Fields{"Id": event.Offers.Offers[0].Id}).Info("Received offer.")
		return event.Offers.Offers[0], resources, nil
	}
	return nil, nil, nil
}
