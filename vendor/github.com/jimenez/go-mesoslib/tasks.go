package mesoslib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jimenez/go-mesoslib/mesosproto"
)

type Volume struct {
	ContainerPath string `json:"container_path,omitempty"`
	HostPath      string `json:"host_path,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

type Task struct {
	ID      string
	Command []string
	Image   string
	Volumes []*Volume
}

func NewTask(image string, command []string) *Task {
	id := make([]byte, 6)
	n, err := rand.Read(id)
	if n != len(id) || err != nil {
		return nil
	}

	return &Task{
		ID:      hex.EncodeToString(id),
		Command: command,
		Image:   image,
	}
}

// Helper for task info object creation
func CreateTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *Task) *mesosproto.TaskInfo {
	taskInfo := mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("mesoscon-demo-task-%s", task.ID)),
		TaskId: &mesosproto.TaskID{
			Value: &task.ID,
		},
		AgentId:   offer.AgentId,
		Resources: resources,
		Command:   &mesosproto.CommandInfo{},
	}

	// Set value only if provided
	if task.Command[0] != "" {
		taskInfo.Command.Value = &task.Command[0]
	}

	// Set args only if they exist
	if len(task.Command) > 1 {
		taskInfo.Command.Arguments = task.Command[1:]
	}

	// Set the docker image if specified
	if task.Image != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &task.Image,
			},
		}

		for _, v := range task.Volumes {
			var (
				vv   = v
				mode = mesosproto.Volume_RW
			)

			if vv.Mode == "ro" {
				mode = mesosproto.Volume_RO
			}

			taskInfo.Container.Volumes = append(taskInfo.Container.Volumes, &mesosproto.Volume{
				ContainerPath: &vv.ContainerPath,
				HostPath:      &vv.HostPath,
				Mode:          &mode,
			})
		}

		taskInfo.Command.Shell = proto.Bool(false)
	}

	return &taskInfo
}
