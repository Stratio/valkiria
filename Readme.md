## Valkiria

Valkiria is an open source project based on http://principlesofchaos.org/ concepts an inspired by Netflix ChaosMonkey tool.
It is designed to be compatible with Stratio PAAS but can work with any cluster of Mesos.
Can kill processes randomly or selectively so it is a great tool for:
* Test high availability (HA)
* Test fault tolerance (FT)
* Prediction of error cases in the underlying technologies
* Benchmarking of HA and FT capabilities of different distributed architectures
* Sizing productive environments
* Understand environments


## Modules

* Valkiria Chaos
* Valkiria Admin - coming soon
* Valkiria Security - coming soon
* Valkiria Metrics - coming soon

## Dependencies
* Systemd init system
* DBus

## Architecture

Is based on orchestrator/agent architecture. The orchestrator has the responsibility to orchestrate Chaos sessions, 
discover the live agents, know the name of the tasks that are running and manage security (Enterprise Edition).
The agent has the responsibility to kill the process by task name. Agent kills daemons via systemd API from DBus, 
kills task via 'kill' and kills docker task via 'docker kill'.


## Installation

```console
$ go get -u github.com/stratio/valkiria
$ make install
```


## Usage

```console
$ valkiria orchestrator ip=10.200.27.1:9050 log=INFO
$ valkiria agent ip=10.200.27.2:9050 log=INFO
```


## Warning
Do not use in production, is in Alpha version. You need to run the agent process with root privileges. For tests of several packages is also necessary, else will be skipped.