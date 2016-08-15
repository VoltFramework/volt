package scheduler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/jimenez/go-mesoslib/mesosproto/schedulerproto"
)

func (lib *SchedulerLib) send(call *schedulerproto.Call, statusExpected int) (io.ReadCloser, error) {
	return lib.sendDetail(call, statusExpected, nil)
}

func (lib *SchedulerLib) sendDetail(call *schedulerproto.Call, statusExpected int, f func(r *http.Response)) (io.ReadCloser, error) {
	body, err := proto.Marshal(call)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://"+lib.master+ENDPOINT, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Accept", "application/json")
	if lib.MesosStreamId != "" {
		req.Header.Set("Mesos-Stream-Id", lib.MesosStreamId)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != statusExpected {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", body)
	}

	if f != nil {
		f(resp)
	}
	return resp.Body, nil

}
