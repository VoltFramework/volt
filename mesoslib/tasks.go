package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

var none = 0.0

func (m *MesosLib) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, command, ID string, state *mesosproto.TaskState) error {
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
				Resources: resources,
				Command: &mesosproto.CommandInfo{
					Value: &command,
				},
			},
		},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{RefuseSeconds: &none},
	}, "mesos.internal.LaunchTasksMessage"); err != nil {
		return err
	}

	for {
		event := <-m.GetEvent(mesosproto.Event_UPDATE)
		*state = *event.Update.Status.State
		switch *event.Update.Status.State {
		case mesosproto.TaskState_TASK_STARTING:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is starting.")
		case mesosproto.TaskState_TASK_RUNNING:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is running.")
		case mesosproto.TaskState_TASK_FINISHED:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Info("Task is finished.")
			return nil
		case mesosproto.TaskState_TASK_FAILED:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task has failed.")
			return nil
		case mesosproto.TaskState_TASK_KILLED:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task was killed.")
			return nil
		case mesosproto.TaskState_TASK_LOST:
			m.Log.WithFields(logrus.Fields{"ID": ID, "message": event.Update.Status.GetMessage()}).Warn("Task was lost.")
			return nil
		}
	}
	return nil
}
