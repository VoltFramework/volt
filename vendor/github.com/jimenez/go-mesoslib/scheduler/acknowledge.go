package scheduler

import (
	"github.com/jimenez/go-mesoslib/mesosproto"
	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) Acknowledge(taskId *mesosproto.TaskID, AgentId *mesosproto.AgentID, UUID []byte) error {
	call := &schedulerproto.Call{
		FrameworkId: lib.frameworkID,
		Type:        schedulerproto.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &schedulerproto.Call_Acknowledge{
			AgentId: AgentId,
			TaskId:  taskId,
			Uuid:    UUID,
		},
	}

	_, err := lib.send(call, 202)
	return err
}
