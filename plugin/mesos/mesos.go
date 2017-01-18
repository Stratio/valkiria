package mesos

import (
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/proc"
	"github.com/pkg/errors"
	"regexp"
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
	mesosHealthCheck    = "mesos-health-ch"
	dockerExecutor      = "docker"
)

const (
	all      string = "0"
	task     string = "1"
	executor string = "2"
)

type MesosConfig struct {
	DaemonConfigString         []string
	DaemonListForChildServices []string
	DockerConfigPattern        string
	BlackListServices          []string
}

func NewMesosConfig() *MesosConfig {
	return &MesosConfig{
		DaemonConfigString:         []string{mesosMaster, mesosAgentPublic, mesosAgent, dcosMarathon, dcosZookeeper},
		DaemonListForChildServices: []string{mesosAgent, mesosAgentPublic},
		DockerConfigPattern:        "^\\/mesos-.*",
		BlackListServices:          []string{mesosAgentLogrotate, mesosDockerExecutor, dockerExecutor, mesosHealthCheck},
	}
}

func (m *MesosConfig) FindAndKill() func([]plugin.Process, string, string) ([]plugin.Process, []error) {
	return func(processList []plugin.Process, name string, properties string) ([]plugin.Process, []error) {
		var process []plugin.Process
		var errorList []error
		//.
		validateProperties, _ := regexp.Compile("^killExecutor=[012]+")
		if !validateProperties.MatchString(properties) {
			errorList = append(errorList, errors.New("mesosPlugin - Bad properties field. Use killExecutor=[012]"))
			return nil, errorList
		}
		var err error
		for _, pro := range processList {
			switch pro.(type) {
			case *proc.Daemon:
				if pro.(*proc.Daemon).KillName == name {
					err = pro.Kill()
					if err == nil {
						pro.(*proc.Daemon).ChaosTimeStamp = time.Now().UnixNano()
						process = append(process, pro)
					}
				}
			case *proc.Docker:
				if pro.(*proc.Docker).KillName == name {
					err = pro.Kill()
					if err == nil {
						pro.(*proc.Docker).ChaosTimeStamp = time.Now().UnixNano()
						process = append(process, pro)
					}
				}
			case *proc.Service:
				if pro.(*proc.Service).KillName == name {
					switch strings.Split(properties, "=")[1] {
					case all:
						err = pro.Kill()
						if err == nil {
							pro.(*proc.Service).ChaosTimeStamp = time.Now().UnixNano()
							process = append(process, pro)
						}
					case task:
						if !pro.(*proc.Service).Executor {
							err = pro.Kill()
							if err == nil {
								pro.(*proc.Service).ChaosTimeStamp = time.Now().UnixNano()
								process = append(process, pro)
							}
						}
					case executor:
						if pro.(*proc.Service).Executor {
							err = pro.Kill()
							if err == nil {
								pro.(*proc.Service).ChaosTimeStamp = time.Now().UnixNano()
								process = append(process, pro)
							}
						}
					}
				}
			}
			if err != nil {
				errorList = append(errorList, err)
			}
		}
		return process, errorList
	}
}
