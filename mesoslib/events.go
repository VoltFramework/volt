package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/vieux/volt/mesosproto"
)

func (m *MesosLib) GetEvent() *mesosproto.Event {
	e := <-m.events
	m.log.WithFields(logrus.Fields{"type": e.Type}).Debug("Received event from master.")
	return e
}
