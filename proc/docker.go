package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"regexp"
	"strings"
	"time"
)

const (
	running     = "running"
	taskMesosId = "MESOS_TASK_ID"
	equal       = "="
)

type docker struct {
	Id             string
	Name           string
	TaskName       string
	Image          string
	ChaosTimeStamp int
}

type function func(conatiner types.Container) (*docker, error)

func newDockerClientInitialization() (c *client.Client, err error) {
	if c, err = client.NewEnvClient(); err != nil {
		log.Infof("proc.docker.newDockerClientInitialization - ERROR: '%v'", err.Error())
	}
	return c, err
}

func (d *docker) Kill() (err error) {
	log.Debug("proc.docker.Kill")
	var clientDocker *client.Client
	if clientDocker, err = newDockerClientInitialization(); err == nil {
		d.ChaosTimeStamp = time.Now().UTC().Nanosecond()
		log.Infof("proc.docker.Kill - '%v' '%v' '%v' '%v'", d.Id, d.Name, d.Image, d.ChaosTimeStamp)
		err = clientDocker.ContainerKill(context.Background(), d.Id, "KILL")
		if err != nil {
			log.Infof("proc.docker.Kill - ERROR: '%v'", err.Error())
		}
	}
	return
}

func ReadAllDockers(patternContainerName string, functionToAdd function) (res []docker, err error) {
	log.Debug("proc.docker.ReadAllDockers")
	var validName = regexp.MustCompile(patternContainerName)
	var c *client.Client
	if c, err = newDockerClientInitialization(); err == nil {
		if containers, _ := c.ContainerList(context.Background(), types.ContainerListOptions{All: true}); err == nil {
			for _, container := range containers {
				if len(container.Names) > 0 && validName.Match([]byte(container.Names[0])) && strings.Contains(container.State, running) {
					if d, e := functionToAdd(container); e == nil {
						if d != nil {
							res = append(res, *d)
							log.Debugf("proc.docker.ReadAllDockers - append - '%v' '%v' '%v'", d.TaskName, d.Image, d.Name)
						}

					} else {
						log.Infof("proc.docker.ReadAllDockers - ERROR: %v", e.Error())
					}
				}
			}
		}
	}
	log.Debugf("proc.service.ReadAllDockers - lenService: '%v'", len(res))
	return
}

var FunctionToAddDockerContainerMesosCluster = func(container types.Container) (*docker, error) {
	//TODO: other new client is mandatory. bug in docker api
	c2, _ := client.NewEnvClient()
	insp, _ := c2.ContainerInspect(context.Background(), container.ID)
	var taskEnv string
	for _, e := range insp.Config.Env {
		if strings.Contains(e, taskMesosId) {
			taskEnv = strings.Split(e, equal)[1]
		}
	}
	return &docker{Id: container.ID, Name: container.Names[0], Image: container.Image, TaskName: taskEnv}, nil
}
