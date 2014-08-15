package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, command, ID string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID, "command": command, "offerId": offer.Id}).Info("Launching task...")

	return m.send(&mesosproto.LaunchTasksMessage{
		FrameworkId: m.frameworkInfo.Id,
		Tasks: []*mesosproto.TaskInfo{
			&mesosproto.TaskInfo{
				Name: &command,
				TaskId: &mesosproto.TaskID{
					Value: &ID,
				},
				SlaveId:   offer.SlaveId,
				Resources: resources,
				Command: &mesosproto.CommandInfo{
					Value: &command,
				},
			},
		},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}
