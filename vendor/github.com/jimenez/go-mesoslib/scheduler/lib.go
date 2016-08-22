package scheduler

import (
	"github.com/golang/protobuf/proto"
	"github.com/jimenez/go-mesoslib/mesosproto"
)

const ENDPOINT = "/master/api/v1/scheduler"

type SchedulerLib struct {
	name          string
	master        string
	frameworkInfo *mesosproto.FrameworkInfo
	frameworkID   *mesosproto.FrameworkID
	MesosStreamId string
	tasks         map[string]*mesosproto.AgentID
}

func New(master, name string) *SchedulerLib {
	return &SchedulerLib{
		name:          name,
		master:        master,
		frameworkInfo: &mesosproto.FrameworkInfo{Name: &name, User: proto.String("root")},
		tasks:         make(map[string]*mesosproto.AgentID),
	}
}

type OfferHandler func(offer *mesosproto.Offer)
type TaskStatusHandler func(taskStatus *mesosproto.TaskStatus)
