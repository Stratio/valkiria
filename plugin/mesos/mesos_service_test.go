package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"testing"
	//"github.com/Stratio/valkiria/proc"
	"github.com/Stratio/valkiria/plugin"
)

const (
	killName = "mesos-32156487435168416831"
)

var (
	testReadKillService = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DaemonConfigString = []string{unit}
		pluginInstance.DockerConfigPattern = "^/testValkiria"
		pluginInstance.DaemonListForChildServices = []string{unit}

		var daemons, _ = pluginInstance.GetDaemons()()
		var services, _ = pluginInstance.GetServices()()
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(append(append(aux, daemons...), services...), dockers...), killName, "killExecutor=0")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) < 1 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element")
		}
	}
	testReadKillServiceTask = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DaemonConfigString = []string{unit}
		pluginInstance.DockerConfigPattern = "^/testValkiria"
		pluginInstance.DaemonListForChildServices = []string{unit}

		var daemons, _ = pluginInstance.GetDaemons()()
		var services, _ = pluginInstance.GetServices()()
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		_, err := functionRead(append(append(append(aux, daemons...), services...), dockers...), killName, "killExecutor=1")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
	}
	testReadKillServiceExecutor = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DaemonConfigString = []string{unit}
		pluginInstance.DockerConfigPattern = "^/testValkiria"
		pluginInstance.DaemonListForChildServices = []string{unit}

		var daemons, _ = pluginInstance.GetDaemons()()
		var services, _ = pluginInstance.GetServices()()
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(append(append(aux, daemons...), services...), dockers...), killName, "killExecutor=2")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) < 1 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element")
		}
	}
	testIsInBlackList = func(t *testing.T) {
		b := isInBlackList(test.Unit, []string{test.Unit})
		if !b {
			t.Fatalf("proc.testChaos - ERROR: Is in black list and do not match")
		}
		f := isInBlackList(test.Unit, []string{"fakeUnit"})
		if f {
			t.Fatalf("proc.testChaos - ERROR: Is not in black list and match")
		}
	}
)

func TestServiceLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	t.Run("testReadKillService", testReadKillService)
	t.Run("testIsInBlackList", testIsInBlackList)
}

func TestServiceTaskLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	t.Run("testReadKillServiceTask", testReadKillServiceTask)
}

func TestServiceExecutorLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	t.Run("testReadKillServiceExecutor", testReadKillServiceExecutor)
}
