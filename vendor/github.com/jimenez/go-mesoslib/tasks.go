package mesoslib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/jimenez/go-mesoslib/mesosproto"
)

type Volume struct {
	ContainerPath string `json:"container_path,omitempty"`
	HostPath      string `json:"host_path,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

type Task struct {
	ID       string
	Command  []string
	Image    string
	Volumes  []*Volume
	Executor *mesosproto.ExecutorInfo
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
		Name: proto.String(fmt.Sprintf("container-demo-task-%s", task.ID)),
		TaskId: &mesosproto.TaskID{
			Value: &task.ID,
		},
		AgentId:   offer.AgentId,
		Resources: resources,
		//		Command:   &mesosproto.CommandInfo{},
		Executor: task.Executor,
	}

	// // Set value only if provided
	// if task.Command[0] != "" {
	// 	taskInfo.Command.Value = &task.Command[0]
	// }

	// Set args only if they exist
	// if len(task.Command) > 1 {
	// 	taskInfo.Command.Arguments = task.Command
	// }

	// Set the docker image if specified
	cmd := "command"
	args := strings.Join(task.Command, "\", \"")
	if task.Image != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &task.Image,
				Parameters: []*mesosproto.Parameter{
					&mesosproto.Parameter{
						Key:   &cmd,
						Value: &args,
					},
				},
			},
		}

		// for _, v := range task.Volumes {
		// 	var (
		// 		vv   = v
		// 		mode = mesosproto.Volume_RW
		// 	)

		// 	if vv.Mode == "ro" {
		// 		mode = mesosproto.Volume_RO
		// 	}

		// 	taskInfo.Container.Volumes = append(taskInfo.Container.Volumes, &mesosproto.Volume{
		// 		ContainerPath: &vv.ContainerPath,
		// 		HostPath:      &vv.HostPath,
		// 		Mode:          &mode,
		// 	})
		// }

		//		taskInfo.Command.Shell = proto.Bool(false)
	}
	logrus.Infof("SENDING TASKINFO: %#v", taskInfo)
	spew.Sdump(taskInfo)
	return &taskInfo
}
