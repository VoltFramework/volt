package mesoslib

import (
	"encoding/json"
	"net/http"
)

type Metrics struct {
	TotalCpus float64 `json:"total_cpus"`
	TotalMem  float64 `json:"total_mem"`
	TotalDisk float64 `json:"total_disk"`
	UsedCpus  float64 `json:"used_cpus"`
	UsedMem   float64 `json:"used_mem"`
	UsedDisk  float64 `json:"used_disk"`
}

func (m *MesosLib) Metrics() (*Metrics, error) {
	resp, err := http.Get("http://" + m.master + "/master/state.json")
	if err != nil {
		return nil, err
	}

	data := struct {
		Frameworks []struct {
			Tasks []struct {
				Resources struct {
					Cpus float64
					Mem  float64
					Disk float64
				}
			}
			Id string
		}
		Slaves []struct {
			Resources struct {
				Cpus float64
				Mem  float64
				Disk float64
			}
		}
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	resp.Body.Close()

	var (
		metrics Metrics
	)

	for _, framework := range data.Frameworks {
		for _, task := range framework.Tasks {
			metrics.UsedMem += task.Resources.Mem
			metrics.UsedCpus += task.Resources.Cpus
			metrics.UsedDisk += task.Resources.Disk
		}
	}

	for _, slave := range data.Slaves {
		metrics.TotalMem += slave.Resources.Mem
		metrics.TotalCpus += slave.Resources.Cpus
		metrics.TotalDisk += slave.Resources.Disk
	}

	return &metrics, nil
}
