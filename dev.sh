#!/bin/sh

docker run -p 2181:2181 --name zk -d -p 2888:2888 -p 3888:3888 relateiq/zookeeper
docker run -p 5050:5050 --name master --link=zk:zk -d jimenez/mesos-dev:iwyu /mesos/build/bin/mesos-master.sh --zk=zk://zk:2181/mesos --quorum=1 --work_dir=/
docker run -p 5051:5051 --name slave --link=zk:zk -d -v /sys/fs/cgroup/:/sys/fs/cgroup/ -v /usr/bin/docker:/usr/bin/docker -v /var/run/docker.sock:/var/run/docker.sock jimenez/mesos-dev:iwyu /mesos/build/bin/mesos-slave.sh --master=zk://zk:2181/mesos --containerizers=docker,mesos --hostname="198.27.68.58"