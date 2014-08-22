package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func createScalarResource(name string, value float64) *mesosproto.Resource {
	return &mesosproto.Resource{
		Name:   &name,
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: &value},
	}
}

func (m *MesosLib) BuildResources(cpus, mem, disk float64) []*mesosproto.Resource {
	m.Log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem, "disk": disk}).Info("Building resources...")
	var resources = []*mesosproto.Resource{}

	if cpus > 0 {
		resources = append(resources, createScalarResource("cpus", cpus))
	}

	if mem > 0 {
		resources = append(resources, createScalarResource("mem", mem))
	}

	if disk > 0 {
		resources = append(resources, createScalarResource("disk", disk))
	}

	return resources
}
