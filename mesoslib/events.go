package mesoslib

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

type events map[mesosproto.Event_Type]chan *mesosproto.Event

func (m *MesosLib) AddEvent(eventType mesosproto.Event_Type, event *mesosproto.Event) error {
	m.Log.WithFields(logrus.Fields{"type": eventType}).Debug("Received event from master.")
	if c, ok := m.events[eventType]; ok {
		c <- event
		return nil
	}
	return fmt.Errorf("unknown event type: %v", eventType)
}

func (m *MesosLib) GetEvent(kind mesosproto.Event_Type) chan *mesosproto.Event {
	if c, ok := m.events[kind]; ok {
		return c
	} else {
		return nil
	}
}
