package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/VoltFramework/volt/mesosproto"
	"github.com/samuel/go-zookeeper/zk"
)

func getMasterFromZK() (string, error) {
	var (
		masterInfo mesosproto.MasterInfo
		parts      = strings.SplitN(master[5:], "/", 2)
		uris       = parts[0]
		path       = "/" + parts[1]
		data       []byte
	)

	c, _, err := zk.Connect(strings.Split(uris, ","), time.Second)
	if err != nil {
		return "", err
	}
	children, _, err := c.Children(path)
	if err != nil {
		return "", err
	}

	sort.Strings(children)

	for _, child := range children {
		if strings.HasPrefix(child, "info_") {
			data, _, err = c.Get(path + "/" + child)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if data == nil {
		return "", errors.New("Unable to get master from ZooKeeper")
	}

	if err := proto.Unmarshal(data, &masterInfo); err != nil {
		return "", err
	}

	parts = strings.Split(masterInfo.GetPid(), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("Received invalid PID from Zookeeper: %s", masterInfo.GetPid())
	}
	return parts[1], nil
}
