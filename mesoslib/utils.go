package mesoslib

import (
	"bytes"
	"fmt"
	"net/http"

	"code.google.com/p/goprotobuf/proto"
	"github.com/VoltFramework/volt/mesosproto"
)

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
