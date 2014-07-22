package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RequestOffer(frameworkInfo *mesosproto.FrameworkInfo, cpus, mem float64) error {
	m.log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem}).Info("Requesting offers...")

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
		return err
	}

	event := <-m.events

	for _, offer := range event.Offers.Offers {
		m.log.Warnln("Received offer: %#v", offer)
	}

	return nil
}
