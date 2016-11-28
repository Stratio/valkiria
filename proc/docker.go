package proc

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/distribution/context"
	"strings"
	"time"
	log "github.com/Sirupsen/logrus"
)

const (
	mesos = "mesos"
	running = "running"
	taskMesosId = "MESOS_TASK_ID"
	equal = "="
)

type docker struct{
	Id string
	Name string
	TaskName string
	Image string
	ChaosTimeStamp int
}

func (d *docker) Kill () (err error){
	log.Debug("proc.docker.Kill")
	if client, err := client.NewEnvClient(); err == nil{
		d.ChaosTimeStamp = time.Now().UTC().Nanosecond()
		log.Infof("proc.docker.Kill - '%v' '%v' '%v' '%v'", d.Id, d.Name, d.Image, d.ChaosTimeStamp)
		err = client.ContainerKill(context.Background(), d.Id, "KILL")
	}
	if err != nil {
		log.Infof("proc.docker.Kill - ERROR: '%v'", err.Error())
	}
	return
}

func ReadAllDockers ()(res []docker, err error) {
	log.Debug("proc.docker.ReadAllDockers")
	if c, err := client.NewEnvClient(); err == nil{
		if containers, err := c.ContainerList(context.Background(), types.ContainerListOptions{All: true}); err == nil {
			for _, container := range containers {
				if len(container.Names) > 0 && strings.Contains(container.Names[0], mesos) && strings.Contains(container.State, running) {
					//TODO: other new client is mandatory. bug in docker api
					c2, _ := client.NewEnvClient()
					insp, err := c2.ContainerInspect(context.Background(), container.ID)
					if err == nil {
						var taskEnv string
						for _, e := range insp.Config.Env {
							if strings.Contains(e, taskMesosId) {
								taskEnv = strings.Split(e, equal)[1]
							}
						}
						res = append(res, docker{Id: container.ID, Name: container.Names[0], Image: container.Image, TaskName: taskEnv})
						log.Debugf("proc.docker.ReadAllDockers - append - '%v' '%v' '%v'", taskEnv, container.Image, container.Names[0])
					}
				}
			}
		}
	}
	if err != nil {
		log.Infof("proc.docker.ReadAllDockers ERROR: '%v'", err.Error())
	}
	log.Debugf("proc.service.ReadAllDockers - lenService: '%v'", len(res))
	return
}