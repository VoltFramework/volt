package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/api"
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
	flag "github.com/dotcloud/docker/pkg/mflag"
)

func main() {
	var (
		log             = logrus.New()
		port            = flag.Int([]string{"p", "-port"}, 8080, "Port to listen on for the API")
		master          = flag.String([]string{"m", "-master"}, "localhost:5050", "Master to connect to")
		debug           = flag.Bool([]string{"D", "-debug"}, false, "")
		user            = flag.String([]string{"u", "-user"}, "root", "User to execute tasks as")
		ip              = flag.String([]string{"-ip"}, "", "IP address to listen on [default: autodetect]")
		frameworkName   = "volt"
		registerTimeout = 5 * time.Second
		frameworkInfo   = &mesosproto.FrameworkInfo{Name: &frameworkName, User: user}
	)

	flag.Parse()

	if *debug {
		log.Level = logrus.DebugLevel
	}

	// initialize MesosLib
	m := mesoslib.NewMesosLib(*master, log, frameworkInfo, *ip)

	// try to register against the master
	if err := m.RegisterFramework(); err != nil {
		log.Fatal(err)
	}

	// wait for the registered event
	select {
	case event := <-m.GetEvent(mesosproto.Event_REGISTERED):
		log.WithFields(logrus.Fields{"FrameworkId": *event.Registered.FrameworkId.Value}).Info("Registration successful.")
	case <-time.After(registerTimeout):
		log.WithField("--ip", *ip).Fatal("Registration timed out. --ip must route to this host from the mesos-master.")
	}

	// once we are registered, start the API
	if err := api.NewAPI(m).ListenAndServe(*port); err != nil {
		log.Fatal(err)
	}

	//TODO catch signal to unregister cleanly
	if err := m.UnRegisterFramework(); err != nil {
		log.Fatal(err)
	}
}
