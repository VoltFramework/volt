package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	mesoslib "github.com/jimenez/go-mesoslib"
	"github.com/jimenez/go-mesoslib/mesosproto"
	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *mesoslib.Task) error {
	taskInfo := mesoslib.CreateTaskInfo(offer, resources, task)

	taskInfo.Executor.FrameworkId = lib.frameworkID
	taskInfo.Executor.Command.Value = proto.String("executor -framework_id " + lib.frameworkID.GetValue() + " -executor_id " + taskInfo.Executor.ExecutorId.GetValue())

	call := &schedulerproto.Call{
		FrameworkId: lib.frameworkID,
		Type:        schedulerproto.Call_ACCEPT.Enum(),
		Accept: &schedulerproto.Call_Accept{
			OfferIds: []*mesosproto.OfferID{
				offer.Id,
			},
			Operations: []*mesosproto.Offer_Operation{
				&mesosproto.Offer_Operation{
					Type: mesosproto.Offer_Operation_LAUNCH.Enum(),
					Launch: &mesosproto.Offer_Operation_Launch{
						TaskInfos: []*mesosproto.TaskInfo{taskInfo},
					},
				},
			},
		},
	}
	logrus.Infof("ARGUMENTS : %#v", taskInfo.Executor.Command.Arguments)
	_, err := lib.send(call, 202)
	return err
}
