package task

import (
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
)

type Task struct {
	ID          string   `json:"id"`
	Command     string   `json:"cmd"`
	Cpus        float64  `json:"cpus,string"`
	Disk        float64  `json:"disk,string"`
	Mem         float64  `json:"mem,string"`
	Files       []string `json:"files"`
	DockerImage string   `json:"docker_image"`

	SlaveId *string               `json:"slave_id,string"`
	State   *mesosproto.TaskState `json:"state,string"`
	Volumes []*mesoslib.Volume    `json:"volumes,omitempty"`
}
