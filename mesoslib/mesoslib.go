package mesoslib

import (
	"net"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

type MesosLib struct {
	master        string
	Log           *logrus.Logger
	ip            string
	port          int
	frameworkInfo *mesosproto.FrameworkInfo

	events events
}

func NewMesosLib(master string, log *logrus.Logger, frameworkInfo *mesosproto.FrameworkInfo, ip string) *MesosLib {
	m := &MesosLib{
		Log:           log,
		master:        master,
		port:          9091,
		frameworkInfo: frameworkInfo,
		ip:            ip,
		events: events{
			mesosproto.Event_REGISTERED: make(chan *mesosproto.Event),
			mesosproto.Event_OFFERS:     make(chan *mesosproto.Event),
			mesosproto.Event_UPDATE:     make(chan *mesosproto.Event),
		},
	}

	if m.ip == "" {
		name, err := os.Hostname()
		if err != nil {
			log.Fatalf("Failed to get hostname: %+v", err)
		}

		addrs, err := net.LookupHost(name)
		if err != nil {
			log.Fatalf("Failed to get address for hostname %q: %+v", name, err)
		}

		for _, addr := range addrs {
			if m.ip == "" || !strings.HasPrefix(addr, "127") {
				m.ip = addr
			}
		}
	}
	m.initAPI()
	return m
}
