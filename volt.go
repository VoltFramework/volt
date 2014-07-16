package main

import (
	"github.com/Sirupsen/logrus"
	flag "github.com/dotcloud/docker/pkg/mflag"
	"github.com/vieux/volt/api"
	"github.com/vieux/volt/mesoslib"
	"github.com/vieux/volt/mesosproto"
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
	)

	flag.Parse()

	if *debug {
		log.Level = logrus.Debug
	}

	// initialize MesosLib
	m := mesoslib.NewMesosLib(*master, log)

	log.Infof("Starting %s...", frameworkName)
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}

	// try to register against the master
	if err := m.RegisterFramework(frameworkInfo); err != nil {
		log.Fatal(err)
	}

	// wait for the registered event
	event := m.GetEvent()
	if *event.Type != mesosproto.Event_REGISTERED {
		log.Fatalln("Unsuccessful registration.")
	}

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
