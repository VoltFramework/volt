package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) LaunchTask(offer *mesosproto.Offer, command, ID string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID, "command": command, "offerId": offer.Id}).Info("Launching task...")

	if err := m.send(&mesosproto.LaunchTasksMessage{
		FrameworkId: m.frameworkInfo.Id,
		Tasks: []*mesosproto.TaskInfo{
			&mesosproto.TaskInfo{
				Name: &command,
				TaskId: &mesosproto.TaskID{
					Value: &ID,
				},
				SlaveId:   offer.SlaveId,
				Resources: offer.Resources,
				Command: &mesosproto.CommandInfo{
					Value: &command,
				},
			},
		},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage"); err != nil {
		return err
	}

	event := <-m.GetEvent(mesosproto.Event_UPDATE)

	if err := m.send(&mesosproto.StatusUpdateAcknowledgementMessage{
		FrameworkId: m.frameworkInfo.Id,
		SlaveId:     event.Update.Status.SlaveId,
		TaskId:      event.Update.Status.TaskId,
		Uuid:        event.Update.Uuid,
	}, "mesos.internal.StatusUpdateAcknowledgementMessage"); err != nil {
		return err
	}

	return nil
}
