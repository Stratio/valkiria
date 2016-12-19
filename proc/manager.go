package proc

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"
)

const (
	mesosMaster         = "dcos-mesos-master.service"
	mesosAgentPublic    = "dcos-mesos-slave-public.service"
	mesosAgent          = "dcos-mesos-slave.service"
	dcosMarathon        = "dcos-marathon.service"
	dcosZookeeper       = "dcos-exhibitor.service"
	mesosAgentLogrotate = "mesos-logrotate"
	mesosDockerExecutor = "mesos-docker-ex"
)

const (
	searchTypeEnum = iota
	daemonEnum
	dockerEnum
	serviceEnum
)

const (
	all = iota
	onlyTask
	onlyExecutor
)

type Process interface {
	Kill() error
}

type Manager struct {
	daemonConfigString  []string
	daemonListForChildServices []string
	dockerConfigPattern string
	blackListServices   []string
	Daemons             []Daemon
	Dockers             []Docker
	Services            []Service
}

func (m *Manager) ConfigManager() {
	//TODO: plugin with config - mesos config default
	m.daemonConfigString = []string{mesosMaster, mesosAgentPublic, mesosAgent, dcosMarathon, dcosZookeeper}
	m.daemonListForChildServices = []string{mesosAgent, mesosAgentPublic}
	m.dockerConfigPattern = "^\\/mesos-.*"
	m.blackListServices = []string{mesosAgentLogrotate, mesosDockerExecutor}
}

// Load processes from SO
func (m *Manager) LoadProcesses() (err error) {
	m.Daemons, err = ReadAllDaemons(m.daemonConfigString)
	m.Dockers, err = ReadAllDockers(m.dockerConfigPattern, FunctionToAddDockerContainerMesosCluster)
	m.Services, err = ReadAllChildProcess(m.Daemons, m.daemonListForChildServices, m.blackListServices)
	return
}

// Shooter is a method that kills tasks by order
// name task
// killExecutor kill executor task too else only service task
// serviceType 1 - daemon; 2 - docker; 3 - service; 0 -search in all; n default case
// true, nil -> ok
// false, nil -> empty slice for docker and/or service
// false, error -> error in kill call
func (m *Manager) Shooter(name string, serviceType int, killExecutor int) (proc []Process, err error) {
	var res = Session{}
	res.Id = int(time.Now().Unix())
	res.Start = int(time.Now().Unix())
	res.SessionType = SHOOTER
	log.Debugf("proc.processes.Shooter - Kill task '%v' type '%v' in session: '%v' '%v' '%v'", name, serviceType, res.Id, res.Start, res.SessionType)
	log.Debugf("proc.processes.Shooter - len(daemon): '%v' len(docker): '%v' len(service): '%v')", len(m.Daemons), len(m.Dockers), len(m.Services))
	switch serviceType {
	case daemonEnum:
		proc, err = daemonsFor(name, m.Daemons)
	case dockerEnum:
		proc, err = dockerFor(name, m.Dockers)

	case serviceEnum:
		proc, err = serviceFor(name, m.Services, killExecutor)

	case searchTypeEnum:
		proc, err = daemonsFor(name, m.Daemons)
		if len(proc) < 1 && err == nil{
			proc, err = dockerFor(name, m.Dockers)
		}
		if len(proc) < 1 && err == nil {
			proc, err = serviceFor(name, m.Services, killExecutor)
		}

	default:
		err = errors.New("SERVICE TYPE NOT FOUND")
	}
	res.Finish = int(time.Now().Unix())
	Sessions = append(Sessions, res)
	return
}

// daemonsFor kill daemon
func daemonsFor(name string, daemons []Daemon) (proc []Process, err error) {
	for _, d := range daemons {
		if strings.Compare(name, d.Name) == 0 {
			if err = d.Kill(); err == nil {
				proc = append(proc, &Daemon{Pid: d.Pid, Name: d.Name, Path: d.Path, ChaosTimeStamp: d.ChaosTimeStamp})
			}
		}
	}
	return
}

// serviceFor kill task or executor
func serviceFor(name string, services []Service, killExecutor int) (proc []Process, err error) {
	for _, d := range services {
		if strings.Compare(name, d.TaskName) == 0 {
			switch killExecutor{
				case all:
					if err = d.Kill(); err == nil {
						proc = append(proc, &Service{Pid: d.Pid, Name: d.Name, TaskName: d.TaskName, Ppid: d.Ppid, Executor: d.Executor, ChaosTimeStamp: d.ChaosTimeStamp})
					}
				case onlyTask:
					if d.Executor == false{
						if err = d.Kill(); err == nil {
							proc = append(proc, &Service{Pid: d.Pid, Name: d.Name, TaskName: d.TaskName, Ppid: d.Ppid, Executor: d.Executor, ChaosTimeStamp: d.ChaosTimeStamp})
						}
					}
				case onlyExecutor:
					if d.Executor == true{
						if err = d.Kill(); err == nil {
							proc = append(proc, &Service{Pid: d.Pid, Name: d.Name, TaskName: d.TaskName, Ppid: d.Ppid, Executor: d.Executor, ChaosTimeStamp: d.ChaosTimeStamp})
						}
					}
			}
		}
	}
	return
}

func dockerFor(name string, docker []Docker) (proc []Process, err error) {
	for _, d := range docker {
		if strings.Compare(name, d.TaskName) == 0 {
			if err = d.Kill(); err == nil {
				proc = append(proc, &Docker{Id: d.Id, Name: d.Name, TaskName: d.TaskName, Image: d.Image, ChaosTimeStamp: d.ChaosTimeStamp})
			}
		}
	}
	return
}