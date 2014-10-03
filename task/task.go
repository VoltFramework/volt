package task

import (
	"time"

	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
)

type Task struct {
	ID          string   `json:"id"`
	Command     string   `json:"cmd"`
	Cpus        float64  `json:"cpus,string"`
	Disk        float64  `json:"disk,string"`
	Mem         float64  `json:"mem,string"`
	Files       []string `json:"files,omitempty"`
	DockerImage string   `json:"docker_image"`

	CreatedTime  time.Time             `json:"created_time"`
	FinishedTime *time.Time            `json:"finished_time,omitempty"`
	SlaveId      *string               `json:"slave_id,string"`
	State        *mesosproto.TaskState `json:"state,string"`
	Volumes      []*mesoslib.Volume    `json:"volumes,omitempty"`
}

type ByCreatedTime []*Task

func (a ByCreatedTime) Len() int           { return len(a) }
func (a ByCreatedTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedTime) Less(i, j int) bool { return a[i].CreatedTime.Before(a[j].CreatedTime) }
