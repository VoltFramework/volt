package api

import "github.com/VoltFramework/volt/mesosproto"

type Task struct {
	ID          string   `json:"id"`
	Command     string   `json:"cmd"`
	Cpus        float64  `json:"cpus,string"`
	Disk        float64  `json:"disk,string"`
	Mem         float64  `json:"mem,string"`
	Files       []string `json:"files"`
	DockerImage string   `json:"docker_image"`

	SlaveId *string               `json:"slave_id,string"`
	State   *mesosproto.TaskState `json:"state,string"`
}

type Tasks []*Task

func (tasks Tasks) Slice(page, per_page int) Tasks {
	var end = per_page

	if page*per_page > len(tasks) {
		page, per_page, end = 0, 0, 0
	} else if page*per_page+per_page > len(tasks) {
		end = len(tasks) - page*per_page
	}
	return tasks[page*per_page : page*per_page+end]
}
