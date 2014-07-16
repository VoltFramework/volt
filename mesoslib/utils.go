package mesoslib

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/goprotobuf/proto"

	"github.com/vieux/volt/mesosproto"
)

var (
	master = flag.String("master", "localhost:5050", "Master to connect to")
	ip     string
	port   = 9091
)

func init() {
	name, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %+v", err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Fatalf("Failed to get address for hostname %q: %+v", name, err)
	}

	for _, addr := range addrs {
		if ip == "" || !strings.HasPrefix(addr, "127") {
			ip = addr
		}
	}
}

func send(call *mesosproto.Call, path string) error {
	data, err := proto.Marshal(call)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/master/%s", *master, path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-type", "application/octet-stream")
	req.Header.Add("Libprocess-From", fmt.Sprintf("mesoslib@%s:%d", ip, port))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while posting to %s: %v", url, err)
	}

	if resp != nil && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code %d received while posting to: %s", resp.StatusCode, url)
	}
	return nil
}
