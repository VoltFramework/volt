package main

import (
	"flag"
	"log"

	"github.com/vieux/gozer/src/mesos"
)

const (
	FRAMEWORK_NAME = "volt"
)

var (
	user = flag.String("user", "", "The user to register as")
)

func main() {
	flag.Parse()

	log.Printf("Registering...")
	err := mesos.Register(*user, FRAMEWORK_NAME)
	if err != nil {
		log.Fatal(err)
	}

}
