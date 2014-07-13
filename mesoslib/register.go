package mesoslib

import (
	"log"

	"github.com/vieux/volt/mesosproto"
)

func RegisterFramework(frameworkInfo *mesosproto.FrameworkInfo) error {
	log.Println("Registering framework:", *frameworkInfo.Name, "against:", *master)

	callType := mesosproto.Call_REGISTER
	registerCall := mesosproto.Call{
		Type:          &callType,
		FrameworkInfo: frameworkInfo,
	}
	return send(&registerCall, "mesos.internal.RegisterFrameworkMessage")
}

func UnRegisterFramework(frameworkInfo *mesosproto.FrameworkInfo) error {
	log.Println("UnRegistering framework:", *frameworkInfo.Name, "against:", *master)
	callType := mesosproto.Call_UNREGISTER
	unRegisterCall := mesosproto.Call{
		Type:          &callType,
		FrameworkInfo: frameworkInfo,
	}
	return send(&unRegisterCall, "mesos.internal.UnRegisterFrameworkMessage")
}
