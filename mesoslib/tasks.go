package mesoslib

import (
	"fmt"
	"strings"

	"code.google.com/p/goprotobuf/proto"
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, command, ID, image string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID, "command": command, "offerId": offer.Id, "dockerImage": image}).Info("Launching task...")

	args := strings.Split(command, " ")

	taskInfo := &mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("volt-task-%s", ID)),
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
		SlaveId:   offer.SlaveId,
		Resources: resources,
		Command: &mesosproto.CommandInfo{
			Value:     &args[0],
			Arguments: args[1:],
			Shell:     proto.Bool(false),
		},
	}

	if image != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &image,
			},
		}
	}

	return m.send(&mesosproto.LaunchTasksMessage{
		FrameworkId: m.frameworkInfo.Id,
		Tasks:       []*mesosproto.TaskInfo{taskInfo},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}

func (m *MesosLib) KillTask(ID string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID}).Info("Killing task...")

	return m.send(&mesosproto.KillTaskMessage{
		FrameworkId: m.frameworkInfo.Id,
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
	}, "mesos.internal.KillTaskMessage")
}
