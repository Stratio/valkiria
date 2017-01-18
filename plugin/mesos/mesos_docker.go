package mesos

import (
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"regexp"
	"strings"
	"github.com/Stratio/valkiria/proc"
	"github.com/Stratio/valkiria/plugin"
)

const (
	running     = "running"
	taskMesosId = "MESOS_TASK_ID"
	frameWorkName = "FRAMEWORK_NAME"
	mesosFrameWorkName = "MESOS_FRAMEWORK_NAME"
	marathonApId = "MARATHON_APP_ID"
	equal       = "="
)

type function func(conatiner types.Container) (*proc.Docker, error)

var FunctionToAddDockerContainerMesosCluster = func(container types.Container) (*proc.Docker, error) {
	//TODO: other new client is mandatory. bug in docker api
	c2, _ := client.NewEnvClient()
	insp, _ := c2.ContainerInspect(context.Background(), container.ID)
	var taskEnv string
	var frameworkEnv string
	for _, e := range insp.Config.Env {
		if strings.Contains(e, taskMesosId) {
			taskEnv = strings.Split(e, equal)[1]
		}
		if strings.Contains(e, frameWorkName) {
			frameworkEnv = strings.Split(e, equal)[1]
		}
		if strings.Contains(e, mesosFrameWorkName) {
			frameworkEnv = strings.Split(e, equal)[1]
		}
		if strings.Contains(e, marathonApId) {
			frameworkEnv = strings.Split(e, equal)[1]
		}
	}
	return &proc.Docker{Id: container.ID, Name: container.Names[0], Image: container.Image, KillName: taskEnv, FrameWorkName: frameworkEnv}, nil
}

func (m *MesosConfig) GetDocker() (func ()([]plugin.Process, error)){
	return func ()([]plugin.Process, error){
		return ReadAllDockers(m.DockerConfigPattern, FunctionToAddDockerContainerMesosCluster)
	}
}

func newDockerClientInitialization() (c *client.Client, err error) {
	return client.NewEnvClient()
}

func ReadAllDockers(patternContainerName string, functionToAdd function) (res []plugin.Process, err error) {
	var validName = regexp.MustCompile(patternContainerName)
	var c *client.Client
	if c, err = newDockerClientInitialization(); err == nil {
		if containers, _ := c.ContainerList(context.Background(), types.ContainerListOptions{All: true}); err == nil {
			for _, container := range containers {
				if len(container.Names) > 0 && validName.Match([]byte(container.Names[0])) && strings.Contains(container.State, running) {
					if d, e := functionToAdd(container); e == nil {
						if d != nil {
							res = append(res, d)
						}

					}
				}
			}
		}
	}
	return
}
