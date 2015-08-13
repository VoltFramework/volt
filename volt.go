package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/api"
	flag "github.com/dotcloud/docker/pkg/mflag"
	"github.com/jimenez/mesoscon-demo/lib"
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
	flag.IntVar(&port, []string{"p", "-port"}, 8080, "Port to listen on for the API")
	flag.StringVar(&master, []string{"m", "-master"}, "localhost:5050", "Master to connect to")
	flag.BoolVar(&debug, []string{"D", "-debug"}, false, "")
	flag.StringVar(&user, []string{"u", "-user"}, "root", "User to execute tasks as")
	flag.StringVar(&ip, []string{"-ip"}, "", "IP address to listen on [default: autodetect]")

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
	m := lib.New(master, "volt")

	// start the API
	api.ListenAndServe(m, port)

	// try to register against the master
	if err := m.Subscribe(); err != nil {
		log.Fatal(err)
	}

	waitForSignals()
}
