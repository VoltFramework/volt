package scheduler

import (
	"encoding/json"
	"net/http"
)

type masterState struct {
	Frameworks []struct {
		Tasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
			Resources  struct {
				Cpus float64
				Mem  float64
				Disk float64
			}
		}
		CompletedTasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
		} `json:"completed_tasks"`
		Id string
	}
	CompletedFrameworks []struct {
		CompletedTasks []struct {
			ExecutorId string `json:"executor_id"`
			Id         string
			SlaveId    string `json:"slave_id"`
		} `json:"completed_tasks"`
		Id string
	} `json:"completed_frameworks"`
	Slaves []struct {
		Id        string
		Pid       string
		Hostname  string
		Resources struct {
			Cpus float64
			Mem  float64
			Disk float64
		}
	}
}

func (m *SchedulerLib) getMasterState() (*masterState, error) {
	resp, err := http.Get("http://" + m.master + "/master/state.json")
	if err != nil {
		return nil, err
	}

	data := new(masterState)

	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}
