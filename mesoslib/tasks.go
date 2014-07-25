package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) LaunchTask(frameworkInfo *mesosproto.FrameworkInfo, offer *mesosproto.Offer, command, ID string) error {
	m.log.WithFields(logrus.Fields{"ID": ID, "command": command}).Info("Launching task...")

	launchType := mesosproto.Call_LAUNCH

	launchCall := &mesosproto.Call{
		FrameworkInfo: frameworkInfo,
		Type:          &launchType,
		Launch: &mesosproto.Call_Launch{
			TaskInfos: []*mesosproto.TaskInfo{
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
		},
	}

	if err := m.send(launchCall, "mesos.internal.LaunchTasksMessage"); err != nil {
		return err
	}

	return nil
}
