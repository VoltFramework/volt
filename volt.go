package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/api"
	"github.com/VoltFramework/volt/mesoslib"
	"github.com/VoltFramework/volt/mesosproto"
	flag "github.com/dotcloud/docker/pkg/mflag"
)

func init() {
}

func main() {
	var (
		log           = logrus.New()
		port          = flag.Int([]string{"p", "-port"}, 8080, "Port to listen on for the API")
		master        = flag.String([]string{"m", "-master"}, "localhost:5050", "Master to connect to")
		debug         = flag.Bool([]string{"D", "-debug"}, false, "")
		frameworkName = "volt"
		user          = ""
		frameworkInfo = &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	)

	flag.Parse()

	if *debug {
		log.Level = logrus.Debug
	}

	// initialize MesosLib
	m := mesoslib.NewMesosLib(*master, log)

	// try to register against the master
	if err := m.RegisterFramework(frameworkInfo); err != nil {
		log.Fatal(err)
	}

	// wait for the registered event
	event := <-m.GetEvent(mesosproto.Event_REGISTERED)

	log.WithFields(logrus.Fields{"FrameworkId": *event.Registered.FrameworkId.Value}).Info("Registration successful.")
	frameworkInfo.Id = event.Registered.FrameworkId

	// once we are registered, start the API
	if err := api.NewAPI(m, frameworkInfo, log).ListenAndServe(*port); err != nil {
		log.Fatal(err)
	}

	//TODO catch signal to unregister cleanly
	if err := m.UnRegisterFramework(frameworkInfo); err != nil {
		log.Fatal(err)
	}
}
