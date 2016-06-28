# Quobyte volume plug-in for Docker

Initial idea from the [quobyte docker plugin](https://github.com/quobyte/docker-volume). Tested with `CentOS 7.2` and `Docker 1.10.3`. This rewrite allows you to use [Quobyte](www.quobyte.com) with Docker without installing the Quobyte client on the host system (e.q. Rancher/CoreOS).

## Build

Get the code

```
$ go get -u github.com/quobyte/api
$ go get -u github.com/johscheuer/go-quobyte-docker
```

### Linux

```
$ go build -o quobyte-docker-plugin -o docker-quobyte-plugin .
$ cp quobyte-docker-plugin /usr/libexec/docker/docker-quobyte-plugin
```

### OSX/MacOS

```
$ GOOS=linux GOARCH=amd64 go build -o docker-quobyte-plugin
$ cp quobyte-docker-plugin /usr/libexec/docker/docker-quobyte-plugin
```

## Setup

- create a user in Quobyte for the plug-in:

  ```
  qmgmt -u <api-url> user config add docker <email>
  ```

- set mandatory configuration in environment

  ```
  export QUOBYTE_API_USER=docker
  export QUOBYTE_API_PASSWORD=...
  export QUOBYTE_API_URL=http://<host>:7860/
  # host[:port][,host:port] or SRV record name
  export QUOBYTE_REGISTRY=quobyte.corp
  ```

- Start the plug-in as root (with above environment)

  ```
  quobyte-docker-volume
  ```

Examples:

```
$ docker volume create --driver quobyte --name <volumename> --opt volume_config=MyConfig
$ docker volume create --driver quobyte --name <volumename>
$ docker volume rm <volumename>
$ docker run --volume-driver=quobyte -v <quobyte volumename>:path
```

- Install systemd files Set the variables in systemd/docker-quobyte.env.sample

```
$ cp systemd/docker-quobyte.env.sample /etc/quobyte/docker-quobyte.env
$ cp docker-quobyte-plugin /usr/libexec/docker/
$ cp systemd/* /lib/systemd/system

$ systemctl daemon-reload
$ systemctl start docker-quobyte-plugin
$ systemctl enable docker-quobyte-plugin
$ systemctl status docker-quobyte-plugin
```

## TODO

- [] Use OPTS to get user and group
