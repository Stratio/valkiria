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
	dockerEnum
	serviceEnum
	daemonEnum
)

type Process interface{
	Kill() error
}

type Manager struct {
	daemonConfigString  []string
	dockerConfigPattern string
	blackListServices   []string
	Daemons             []Daemon
	Dockers             []Docker
	Services            []Service
}

func (m *Manager) ConfigManager() {
	//TODO: plugin with config - mesos config default
	m.daemonConfigString = []string{mesosMaster, mesosAgentPublic, mesosAgent, dcosMarathon, dcosZookeeper}
	m.dockerConfigPattern = "^\\/mesos-.*"
	m.blackListServices = []string{mesosAgentLogrotate}
}

// Load processes from SO
func (m *Manager) LoadProcesses() (err error) {
	m.Daemons, err = ReadAllDaemons(m.daemonConfigString)
	m.Dockers, err = ReadAllDockers(m.dockerConfigPattern, FunctionToAddDockerContainerMesosCluster)
	m.Services, err = ReadAllChildProcess(m.Daemons, m.blackListServices)
	return
}

// Shooter is a method that kills tasks by order
// name task
// killExecutor kill executor task too else only service task
// serviceType 0 - daemon; 1 - docker; 2 - service; 3 -search in all; n default case
// true, nil -> ok
// false, nil -> empty slice for docker and/or service
// false, error -> error in kill call
func (m *Manager) Shooter(name string, serviceType int, killExecutor bool) (proc Process, err error) {
	var timeStart = time.Now()
	log.Debugf("routes.processes.Shooter - Start '%v'", timeStart)
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
		if proc != nil {proc, err = dockerFor(name, m.Dockers)}
		if proc != nil {proc, err = serviceFor(name, m.Services, killExecutor)}

	default:
		err = errors.New("SERVICE TYPE NOT FOUND")
	}
	res.Finish = int(time.Now().Unix())
	Sessions = append(Sessions, res)
	log.Debugf("routes.processes.Shooter - Finish sesion '%v' in: '%v' with result '%v'", res.Id, time.Since(timeStart), proc)
	return
}

func daemonsFor(name string, daemons []Daemon) (proc Process, err error) {
	for _, d := range daemons {
		if strings.Compare(name, d.Name) == 0 {
			log.Infof("proc.processes.Shooter.daemonsFor - Killing task '%v'", name)
			if err = d.Kill(); err != nil {
				log.Infof("proc.processes.Shooter.daemonsFor - Killing task ERROR: '%v'", err.Error())
			} else {
				proc = &d
			}
		}
	}
	return
}

func serviceFor(name string, services []Service, killExecutor bool) (proc Process, err error) {
	for _, d := range services {
		if strings.Compare(name, d.TaskName) == 0 {
			if killExecutor {
				log.Infof("proc.processes.Shooter.serviceFor - Killing task '%v'", name)
				if err = d.Kill(); err != nil {
					log.Infof("proc.processes.Shooter.serviceFor - Killing task ERROR: '%v'", err.Error())
				} else {
					proc = &d
				}
			} else {
				log.Infof("proc.processes.Shooter - Killing task '%v'", name)
				if !d.Executor {
					if err = d.Kill(); err != nil {
						log.Infof("proc.processes.Shooter - Killing task ERROR: '%v'", err.Error())
					} else {
						proc = &d
					}
				}
			}

		}
	}
	return
}

func dockerFor(name string, docker []Docker) (proc Process, err error) {
	for _, d := range docker {
		if strings.Compare(name, d.TaskName) == 0 {
			log.Infof("proc.processes.Shooter.dockerFor - Killing task '%v'", name)
			if err = d.Kill(); err != nil {
				log.Infof("proc.processes.Shooter.dockerFor - Killing task ERROR: '%v'", err.Error())
			} else {
				proc = &d
			}
		}
	}
	return
}
