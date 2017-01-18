package proc

import (
	"github.com/docker/distribution/context"
	"github.com/docker/docker/client"
	"time"
)

type Docker struct {
	Id             string
	Name           string
	KillName       string
	FrameWorkName  string
	Image          string
	ChaosTimeStamp int64
}

func newDockerClientInitialization() (c *client.Client, err error) {
	return client.NewEnvClient()
}

func (d *Docker) Kill() (err error) {
	var clientDocker *client.Client
	if clientDocker, err = newDockerClientInitialization(); err == nil {
		err = clientDocker.ContainerKill(context.Background(), d.Id, "KILL")
		if err != nil {
			d.ChaosTimeStamp = time.Now().UTC().UnixNano()
		}
	}
	return
}
