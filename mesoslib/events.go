package mesoslib

import "github.com/vieux/volt/mesosproto"

var (
	events = make(chan *mesosproto.Event)
)

func GetEvent() *mesosproto.Event {
	return <-events
}
