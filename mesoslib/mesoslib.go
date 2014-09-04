package mesoslib

import (
	"net"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/gorilla/mux"
)

type MesosLib struct {
	master        string
	Log           *logrus.Logger
	ip            string
	port          int
	frameworkInfo *mesosproto.FrameworkInfo
	Router        *mux.Router

	events events
}

func NewMesosLib(master string, log *logrus.Logger, frameworkInfo *mesosproto.FrameworkInfo, ip string, port int) *MesosLib {
	m := &MesosLib{
		Log:           log,
		master:        master,
		frameworkInfo: frameworkInfo,
		ip:            ip,
		port:          port,
		Router:        mux.NewRouter(),
		events: events{
			mesosproto.Event_REGISTERED: make(chan *mesosproto.Event, 64),
			mesosproto.Event_OFFERS:     make(chan *mesosproto.Event, 64),
			mesosproto.Event_UPDATE:     make(chan *mesosproto.Event, 64),
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
