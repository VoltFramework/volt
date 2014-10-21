#!/bin/sh

docker run -p 5050:5050 --name master -d jimenez/mesos-master --ip=0.0.0.0 --work_dir=/ --registry=in_memory
docker run -p 5051:5051 --name slave --link=master:master -v /sys/fs/cgroup/:/sys/fs/cgroup/ -v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock -d jimenez/mesos-slave --master=master:5050 --containerizers=docker,mesos --hostname="198.27.68.58"