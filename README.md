# ![Volt logo](https://raw.githubusercontent.com/VoltFramework/volt/master/static/img/logo.png)

*volt* is a simple Mesos framework written in Go.

## Installation

First get the go dependencies:

```sh
go get github.com/VoltFramework/volt/...
```

Then you can compile `volt` with:

```sh
go install github.com/VoltFramework/volt
```

If `$GOPATH/bin` is in your `PATH`, you can invoke `volt` from the CLI.

## Running

To get started with the framework, run the following commands on a mesos
master node:

```
wget https://github.com/voltframework/volt/releases/downloads/v1.0.0-alpha/volt
chmod +x volt
./volt --master=localhost:5050
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

