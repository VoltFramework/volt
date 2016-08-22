package scheduler

import (
	"github.com/jimenez/go-mesoslib/mesosproto"
	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) MessageTask(taskId, executorId, data string) error {
	call := &schedulerproto.Call{
		FrameworkId: lib.frameworkID,
		Type:        schedulerproto.Call_MESSAGE.Enum(),
		Message: &schedulerproto.Call_Message{
			AgentId:    lib.tasks[taskId],
			ExecutorId: &mesosproto.ExecutorID{Value: &executorId},
			Data:       []byte(data),
		},
	}

	_, err := lib.send(call, 202)
	return err
}
