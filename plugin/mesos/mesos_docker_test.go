package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/test"
	"testing"
)

const (
	patternTask          = "^/testValkiria"
	fakePatternMesosTask = "^fake"
	killNameDocker       = "mesos_1232134"
)

var (
	testReadKillDocker = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DockerConfigPattern = patternTask
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(aux, dockers...), killNameDocker, "killExecutor=0")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) < 1 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element.")
		}
	}
	testReadKillDockerFake = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DockerConfigPattern = patternTask
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(aux, dockers...), fakePatternMesosTask, "killExecutor=0")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) > 0 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element.")
		}
	}
)

func TestDockerLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDockerTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testReadDocker", testReadKillDocker)
	t.Run("testReadKillDockerFake", testReadKillDockerFake)
}
