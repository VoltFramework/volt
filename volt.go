package main

import (
	"flag"
	"log"

	"github.com/vieux/volt/mesoslib"
	"github.com/vieux/volt/mesosproto"
)

func main() {
	var (
		//		port          = flag.Int("port", 4343, "Port to listen on for HTTP endpoint")
		user          = flag.String("user", "", "User to execute commands as")
		frameworkName = "volt"
	)

	flag.Parse()

	log.Println("Starting", frameworkName)
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: user}

	if err := mesoslib.RegisterFramework(frameworkInfo); err != nil {
		log.Fatal(err)
	}

	event := mesoslib.GetEvent()
	log.Println("Received ID:", *event.Registered.FrameworkId.Value)
	if err := mesoslib.UnRegisterFramework(frameworkInfo); err != nil {
		log.Fatal(err)
	}
}
