package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/api"
	"github.com/jimenez/mesoscon-demo/mesoslib/scheduler"
)

var (
	port   int
	master string
	user   string
	ip     string
	debug  bool

	log             = logrus.New()
	frameworkName   = "volt"
	registerTimeout = 5 * time.Second
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on for the API")
	flag.StringVar(&master, "master", "localhost:5050", "Master to connect to")
	flag.BoolVar(&debug, "debug", false, "")
	flag.StringVar(&user, "user", "root", "User to execute tasks as")
	flag.StringVar(&ip, "ip", "", "IP address to listen on [default: autodetect]")

	flag.Parse()
}

func waitForSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		log.Debugf("received signal %s unregistering framework\n", sig)

		os.Exit(0)
	}
}

func setupLogger() error {
	if debug {
		log.Level = logrus.DebugLevel
	}

	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		hook, err := newSentryHook(dsn, map[string]string{
			"master": master,
			"ip":     ip,
			"user":   user,
		})

		if err != nil {
			return err
		}

		log.Hooks.Add(hook)
	}

	return nil
}

func main() {
	if err := setupLogger(); err != nil {
		log.Fatal(err)
	}

	// initialize MesosLib
	m := scheduler.New(master, "volt")

	// start the API
	api := api.ListenAndServe(m, port)

	// try to register against the master
	if err := m.Subscribe(api.HandleOffers); err != nil {
		log.Fatal(err)
	}

	waitForSignals()
}
