package mesoslib

import (
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

type Resources []*mesosproto.Resource

func createScalarResource(name string, value float64) *mesosproto.Resource {
	return &mesosproto.Resource{
		Name:   &name,
		Type:   mesosproto.Value_SCALAR.Enum(),
		Scalar: &mesosproto.Value_Scalar{Value: &value},
	}
}

func (resources Resources) Add(otherResources []*mesosproto.Resource) {
	var found bool

	for _, otherResource := range otherResources {
		found = false
		for _, resource := range resources {
			if *resource.Name == *otherResource.Name {
				*resource.Scalar.Value = resource.Scalar.GetValue() + otherResource.Scalar.GetValue()
				found = true
			}
		}
		if !found {
			resources = append(resources, createScalarResource(*otherResource.Name, otherResource.Scalar.GetValue()))
		}
	}
}

func (resources Resources) Compare(resources, otherResources []*mesosproto.Resource) bool {
	var found bool

	for _, resource := range resources {
		found = false
		for _, otherResource := range otherResources {
			if *resource.Name == *otherResource.Name && otherResource.Scalar.GetValue() >= resource.Scalar.GetValue() {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func NewResources(cpus, mem, disk float64) Resources {
	m.Log.WithFields(logrus.Fields{"cpus": cpus, "mem": mem, "disk": disk}).Info("Building resources...")
	var resources Resources

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
