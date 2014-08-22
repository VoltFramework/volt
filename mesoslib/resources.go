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

func (m *MesosLib) CombineResources(resourcesOld, resourcesNew []*mesosproto.Resource) []*mesosproto.Resource {
	if resourcesOld == nil {
		return resourcesNew
	}
	for _, resourceNew := range resourcesNew {
		found := false
		for _, resourceOld := range resourcesOld {
			if *resourceOld.Name == *resourceNew.Name {
				*resourceOld.Scalar.Value = *resourceOld.Scalar.Value + *resourceNew.Scalar.Value
				found = true
			}
		}
		if !found {
			resourcesOld = append(resourcesOld, resourceNew)
		}
	}
	return resourcesOld
}

func (m *MesosLib) IsEnough(resourcesRequired, resourcesAvailalble []*mesosproto.Resource) bool {
	for _, resourceRequired := range resourcesRequired {
		found := false
		for _, resourceAvailalble := range resourcesAvailalble {
			if *resourceAvailalble.Name == *resourceRequired.Name {
				if *resourceAvailalble.Scalar.Value < *resourceRequired.Scalar.Value {
					return false
				}
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
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
