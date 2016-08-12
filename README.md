# ![Volt logo](https://raw.githubusercontent.com/VoltFramework/volt/master/static/img/logo.png)

*volt* is a simple Mesos framework written in Go.

![build](https://travis-ci.org/VoltFramework/volt.svg?branch=master)

## Installation

The following steps describe how to get started with the Volt framework.

### From Source

First get the go dependencies:

```sh
go get github.com/VoltFramework/volt/...
```

Then you can compile `volt` with:

```sh
go install github.com/VoltFramework/volt
```

If `$GOPATH/bin` is in your `PATH`, you can invoke `volt` from the CLI.

### Latest Release

To get started with the latest release, run the following commands on a mesos
master node:

```
wget https://github.com/voltframework/volt/releases/download/v1.0.0-alpha/volt
chmod +x volt
./volt --master=localhost:5050
```

### API Requests

#### Run a container with data volumes
```json
{
    "cmd": "touch /data/volt",
    "cpus": "0.1",
    "mem": "32",
    "docker_image": "busybox",
    "volumes": [
        {
            "container_path":"/data",
            "host_path":"/volumes/volt"
        }
    ]
}
```

## Creators

**Victor Vieux**

- <http://twitter.com/vieux>
- <http://github.com/vieux>

**Isabel Jimenez**

- <http://twitter.com/ijimene>
- <http://github.com/jimenez>

## Thanks

Thanks to [@dhammon](http://github.com/dhammon) for his work on [gozer](http://github.com/twitter/gozer)

## Licensing

Volt is licensed under the Apache License, Version 2.0. See LICENSE for full license text.

