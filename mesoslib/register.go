package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func (m *MesosLib) RegisterFramework() error {
	m.Log.WithFields(logrus.Fields{"master": m.master}).Info("Registering framework...")

	return m.send(&mesosproto.RegisterFrameworkMessage{
		Framework: m.frameworkInfo,
	}, "mesos.internal.RegisterFrameworkMessage")
}

func (m *MesosLib) UnRegisterFramework() error {
	m.Log.WithFields(logrus.Fields{"master": m.master}).Info("Unregistering framework...")

	return m.send(&mesosproto.UnregisterFrameworkMessage{
		FrameworkId: m.frameworkInfo.Id,
	}, "mesos.internal.UnRegisterFrameworkMessage")
}
