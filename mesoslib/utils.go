package mesoslib

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/goprotobuf/proto"
	"github.com/Sirupsen/logrus"
	"github.com/vieux/volt/mesosproto"
)

type MesosLib struct {
	master string
	log    *logrus.Logger
	ip     string
	port   int

	events chan *mesosproto.Event
}

func NewMesosLib(master string, log *logrus.Logger) *MesosLib {
	m := &MesosLib{
		log:    log,
		master: master,
		port:   9091,
		events: make(chan *mesosproto.Event),
	}

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
	m.initAPI()
	return m
}

func (m *MesosLib) send(call *mesosproto.Call, path string) error {
	data, err := proto.Marshal(call)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/master/%s", m.master, path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-type", "application/octet-stream")
	req.Header.Add("Libprocess-From", fmt.Sprintf("mesoslib@%s:%d", m.ip, m.port))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code %d received while posting to: %s", resp.StatusCode, url)
	}
	return nil
}
