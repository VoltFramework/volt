package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RegisterFramework(frameworkInfo *mesosproto.FrameworkInfo) error {
	m.log.WithFields(logrus.Fields{"master": m.master}).Info("Registering framework...")

	callType := mesosproto.Call_REGISTER
	registerCall := mesosproto.Call{
		Type:          &callType,
		FrameworkInfo: frameworkInfo,
	}
	return m.send(&registerCall, "mesos.internal.RegisterFrameworkMessage")
}

func (m *MesosLib) UnRegisterFramework(frameworkInfo *mesosproto.FrameworkInfo) error {
	m.log.WithFields(logrus.Fields{"master": m.master}).Info("Unregistering framework...")

	callType := mesosproto.Call_UNREGISTER
	unRegisterCall := mesosproto.Call{
		Type:          &callType,
		FrameworkInfo: frameworkInfo,
	}
	return m.send(&unRegisterCall, "mesos.internal.UnRegisterFrameworkMessage")
}
