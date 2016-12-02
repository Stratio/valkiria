package proc

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"strconv"
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
)

const (
	daemonEnum = iota
	dockerEnum
	serviceEnum
	searchTypeEnum
)

type Processes struct {
	Daemons  []daemon
	Dockers  []docker
	Services []service
}

// Load processes from SO
func (p *Processes) LoadProcesses() (err error) {
	p.Daemons, err = ReadAllDaemons([]string{mesosMaster, mesosAgent, dcosMarathon, dcosZookeeper, mesosAgentPublic})
	p.Dockers, err = ReadAllDockers("", nil)
	p.Services, err = ReadAllChildProcess(p.Daemons)
	return
}

// Chaos time
// daemonInt is the number of daemons that can kill in  Chaos time
// serviceInt is the number of services that can kill in  Chaos time
// ddockerInt is the number of dockers that can kill in  Chaos time
// return error if apply
func (p *Processes) Chaos(daemonInt int, serviceInt int, dockerInt int) (err error) {
	log.Debug("proc.processes.Chaos")
	var res = Session{}
	res.Id = int(time.Now().Unix())
	res.Start = int(time.Now().Unix())
	res.SessionType = CHAOS
	for i := 0; i < daemonInt && i < len(p.Daemons); i++ {
		log.Debug("proc.processes.Chaos - daemons")
		rad := rand.Intn(len(p.Daemons))
		log.Debug("proc.processes.Chaos - daemons - random " + strconv.Itoa(int(rad)))
		err = p.Daemons[rad].Kill()
		res.Daemon = append(res.Daemon, p.Daemons[rad])
		p.Daemons = append(p.Daemons[:i], p.Daemons[i+1:]...)
		if err != nil {
			log.Debug("proc.processes.Chaos - daemons - ERROR: " + err.Error())
		}
	}

	for i := 0; i < serviceInt && i < len(p.Services); i++ {
		log.Debug("proc.processes.Chaos - services")
		rad := rand.Intn(len(p.Services))
		log.Debug("proc.processes.Chaos - services - random " + strconv.Itoa(int(rad)))
		err = p.Services[rad].Kill()
		res.Service = append(res.Service, p.Services[rad])
		p.Services = append(p.Services[:i], p.Services[i+1:]...)
		if err != nil {
			log.Debug("proc.processes.Chaos - services - ERROR: " + err.Error())
		}
	}

	for i := 0; i < dockerInt && i < len(p.Dockers); i++ {
		log.Debug("proc.processes.Chaos - dockers")
		rad := rand.Intn(len(p.Dockers))
		log.Debug("proc.processes.Chaos - dockers - random " + strconv.Itoa(int(rad)))
		err = p.Dockers[rad].Kill()
		res.Docker = append(res.Docker, p.Dockers[rad])
		p.Dockers = append(p.Dockers[:i], p.Dockers[i+1:]...)
		if err != nil {
			log.Debug("proc.processes.Chaos - dockers - ERROR: " + err.Error())
		}
	}

	res.Finish = int(time.Now().Unix())
	Sessions = append(Sessions, res)
	log.Debugf("proc.processes.Chaos - Sessions '%v'", len(Sessions))
	if err != nil {
		log.Debug("proc.processes.Chaos - ERROR: " + err.Error())
	}
	return
}

// Shooter is a method that kills tasks by order
// name task
// killExecutor kill executor task too else only service task
// serviceType 0 - daemon; 1 - docker; 2 - service; 3 -search in all; n default case
// true, nil -> ok
// false, nil -> empty slice for docker and/or service
// false, error -> error in kill call
func (p *Processes) Shooter(name string, serviceType int, killExecutor bool) (resBool bool, err error) {
	var timeStart = time.Now()
	log.Debugf("routes.processes.Shooter - Start '%v'", timeStart)
	var res = Session{}
	res.Id = int(time.Now().Unix())
	res.Start = int(time.Now().Unix())
	res.SessionType = SHOOTER
	log.Debugf("proc.processes.Shooter - Kill task '%v' type '%v' in session: '%v' '%v' '%v'", name, serviceType, res.Id, res.Start, res.SessionType)
	log.Debugf("proc.processes.Shooter - len(docker): '%v' len(service): '%v')", len(p.Dockers), len(p.Services))
	switch serviceType {
	case daemonEnum:
		resBool, err = daemonsFor(name, p.Daemons)

	case dockerEnum:
		resBool, err = dockerFor(name, p.Dockers)

	case serviceEnum:
		resBool, err = serviceFor(name, p.Services, killExecutor)

	case searchTypeEnum:
		resBool, err = daemonsFor(name, p.Daemons)
		resBool, err = dockerFor(name, p.Dockers)
		resBool, err = serviceFor(name, p.Services, killExecutor)

	default:
		err = errors.New("Type of service not supported.")
	}
	res.Finish = int(time.Now().Unix())
	Sessions = append(Sessions, res)
	log.Debugf("routes.processes.Shooter - Finish sesion '%v' in: '%v' with result '%v'", res.Id, time.Since(timeStart), resBool)
	return
}

func daemonsFor(name string, daemons []daemon) (resBool bool, err error) {
	for _, d := range daemons {
		if strings.Compare(name, d.Name) == 0 {
			log.Infof("proc.processes.Shooter.daemonsFor - Killing task '%v'", name)
			if err = d.Kill(); err != nil {
				log.Infof("proc.processes.Shooter.daemonsFor - Killing task ERROR: '%v'", err.Error())
			} else {
				resBool = true
			}
		}
	}
	return
}

func serviceFor(name string, services []service, killExecutor bool) (resBool bool, err error) {
	for _, d := range services {
		if strings.Compare(name, d.TaskName) == 0 {
			if killExecutor {
				log.Infof("proc.processes.Shooter.serviceFor - Killing task '%v' '%v'", name)
				if err = d.Kill(); err != nil {
					log.Infof("proc.processes.Shooter.serviceFor - Killing task ERROR: '%v'", err.Error())
				} else {
					resBool = true
				}
			} else {
				log.Infof("proc.processes.Shooter - Killing task '%v' '%v'", name)
				if !d.Executor {
					if err = d.Kill(); err != nil {
						log.Infof("proc.processes.Shooter - Killing task ERROR: '%v'", err.Error())
					} else {
						resBool = true
					}
				}
			}

		}
	}
	return
}

func dockerFor(name string, docker []docker) (resBool bool, err error) {
	for _, d := range docker {
		if strings.Compare(name, d.TaskName) == 0 {
			log.Infof("proc.processes.Shooter.dockerFor - Killing task '%v'", name)
			if err = d.Kill(); err != nil {
				log.Infof("proc.processes.Shooter.dockerFor - Killing task ERROR: '%v'", err.Error())
			} else {
				resBool = true
			}
		}
	}
	return
}
