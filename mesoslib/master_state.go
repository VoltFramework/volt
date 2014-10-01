package mesoslib

import (
	"encoding/json"
	"net/http"
)

func (m *MesosLib) GetSlaveHostname(slaveId string) (string, error) {
	resp, err := http.Get("http://" + m.master + "/master/state.json")
	if err != nil {
		return "", err
	}

	data := struct {
		Slaves []struct {
			Id       string
			Hostname string
		}
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	for _, slave := range data.Slaves {
		if slave.Id == slaveId {
			return slave.Hostname, nil
		}

	}

	return "", nil
}
