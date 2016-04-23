# marathon-consul

This program inspects running apps on a Marathon cluster and register each task within Consul.

## Naming conventions

Service are named in dashed reverse order, only with characters allowed in Consul. So, for instance, for a Marathon app in **/website/production/front** exposing two ports, two services will be registered in Consul through their _port index_ :

 * port0-front-production-website.service.consul
 * port1-front-production-website.service.consul

Other than that, it is a regular Consul service registration, where each Marathon task will get one entry.

## Usage

```
$ ./marathon-consul -help                                                                                                                                                                                                                              master:be067b1 ✱ ◼
Usage of ./marathon-consul:
  -consul string
        address to a consul agent (default "127.0.0.1")
  -healthcheck
        should we expose an healthcheck endpoint
  -healthcheck-port int
        port on which to expose the healthcheck endpoint (default 8080)
  -interval duration
        interval in seconds for marathon checks (default 10s)
  -marathon string
        zookeeper path to marathon (default "/marathon")
  -master string
        comma-separated list of zookeeper servers (default "127.0.0.1:2181")

$ ./marathon-consul -master 192.168.100.11:2181,192.168.100.12:2181 -interval 10s
```

## Docker usage

```
$ docker pull apognu/marathon-consul
$ docker run -master 192.168.100.11:2181,192.168.100.12:2181 -interval 10s -consul 172.17.0.1:8500
```
