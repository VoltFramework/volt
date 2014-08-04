package mesoslib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
)

func (m *MesosLib) getSlavePidAndExecutorId(taskId string) (string, string, error) {
	resp, err := http.Get("http://" + m.master + "/master/state.json")
	if err != nil {
		return "", "", err
	}

	data := struct {
		Frameworks []struct {
			Tasks []struct {
				ExecutorId string `json:"executor_id"`
				Id         string
				SlaveId    string `json:"slave_id"`
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
			Pid string
			Id  string
		}
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", err
	}
	resp.Body.Close()

	var (
		executorId string
		slaveId    string
	)

found1:
	for _, framework := range data.Frameworks {
		if framework.Id != *m.frameworkInfo.Id.Value {
			continue
		}
		for _, task := range framework.Tasks {
			if task.Id == taskId {
				executorId = task.ExecutorId
				slaveId = task.SlaveId
				break found1
			}
		}
		for _, task := range framework.CompletedTasks {
			if task.Id == taskId {
				executorId = task.ExecutorId
				slaveId = task.SlaveId
				break found1
			}
		}
	}

found2:
	for _, framework := range data.CompletedFrameworks {
		if framework.Id != *m.frameworkInfo.Id.Value {
			continue
		}
		for _, task := range framework.CompletedTasks {
			if task.Id == taskId {
				executorId = task.ExecutorId
				slaveId = task.SlaveId
				break found2
			}
		}
	}

	for _, slave := range data.Slaves {
		if slave.Id == slaveId {
			return slave.Pid, executorId, nil
		}
	}

	return "", "", nil
}

func (m *MesosLib) getTaskDirectory(slavePid, executorId string) (string, error) {
	resp, err := http.Get("http://" + slavePid + "/state.json")
	if err != nil {
		return "", err
	}

	data := struct {
		Frameworks []struct {
			Executors []struct {
				Id        string
				Directory string
			}
			CompletedExecutors []struct {
				Id        string
				Directory string
			} `json:"completed_executors"`
			Id string
		}
		CompletedFrameworks []struct {
			CompletedExecutors []struct {
				Id        string
				Directory string
			} `json:"completed_executors"`
			Id string
		} `json:"completed_frameworks"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	for _, framework := range data.Frameworks {
		if framework.Id != *m.frameworkInfo.Id.Value {
			continue
		}
		for _, executor := range framework.Executors {
			if executor.Id == executorId {
				return executor.Directory, nil
			}
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.Id == executorId {
				return executor.Directory, nil
			}
		}
	}

	for _, framework := range data.CompletedFrameworks {
		if framework.Id != *m.frameworkInfo.Id.Value {
			continue
		}
		for _, executor := range framework.CompletedExecutors {
			if executor.Id == executorId {
				return executor.Directory, nil
			}
		}
	}

	return "", nil
}

func (m *MesosLib) readFile(slavePid, directory, filename string) (string, error) {
	v := url.Values{"path": []string{filepath.Join(directory, filename)}, "offset": []string{"0"}}

	resp, err := http.Get("http://" + slavePid + "/files/read.json?" + v.Encode())
	if err != nil {
		return "", err
	}

	data := struct {
		Data string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	resp.Body.Close()

	return data.Data, nil
}

func (m *MesosLib) ReadFile(taskId string, filenames ...string) (map[string]string, error) {
	slavePid, executorId, err := m.getSlavePidAndExecutorId(taskId)
	if err != nil {
		return nil, err
	}
	if slavePid == "" {
		return nil, fmt.Errorf("cannot get slave PID")
	}
	if executorId == "" {
		executorId = taskId
	}
	directory, err := m.getTaskDirectory(slavePid, executorId)
	if err != nil {
		return nil, err
	}

	var files = make(map[string]string)

	for _, filename := range filenames {
		file, err := m.readFile(slavePid, directory, filename)
		if err != nil {
			return nil, err
		}
		files[filename] = file
	}
	return files, nil
}
