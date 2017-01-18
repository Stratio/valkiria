package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/dbus"
	"github.com/Stratio/valkiria/test"
	"os"
	"testing"
	"github.com/Stratio/valkiria/plugin"
)

const (
	unit = "test.service"
	fakeUnit        = "fakeUnit.service"
	unitServicePath = "/tmp/test.service"
	unitServiceLink = "/lib/systemd/system/test.service"
)

var (
	testReadKillDaemon = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DaemonConfigString = []string{unit}
		pluginInstance.DockerConfigPattern = "^/testValkiria"
		pluginInstance.DaemonListForChildServices = []string{unit}

		var daemons, _ = pluginInstance.GetDaemons()()
		var services, _ = pluginInstance.GetServices()()
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(append(append(aux, daemons...), services...), dockers...), unit, "killExecutor=0")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) < 1 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element")
		}
	}
	testReadKillDaemonFake = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		pluginInstance.DaemonConfigString = []string{fakeUnit}
		pluginInstance.DockerConfigPattern = "^/testValkiria"
		pluginInstance.DaemonListForChildServices = []string{fakeUnit}

		var daemons, _ = pluginInstance.GetDaemons()()
		var services, _ = pluginInstance.GetServices()()
		var dockers, _ = pluginInstance.GetDocker()()

		var functionRead = pluginInstance.FindAndKill()
		var aux []plugin.Process
		res, err := functionRead(append(append(append(aux, daemons...), services...), dockers...), fakeUnit, "killExecutor=1")
		if err != nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
		if len(res) > 0 {
			t.Errorf("plugin_mesos.testReadAllDaemons - Should be one element")
		}
	}
)

func TestDaemonLib(t *testing.T) {
	test.SetupDBusTest(t)
	defer test.TearDownDBusTest(t)
	log.SetLevel(test.Level)
	startDBusUnit(t)
	t.Run("testReadKillDaemon", testReadKillDaemon)
	t.Run("testReadKillDaemonFake", testReadKillDaemonFake)
}

//WARNING: this code can produce import cycle not allowed in test phase. Dont move to test_util
// If it is necessary to make changes in the tests, make changes in all the tests impacted
func startDBusUnit(t *testing.T) {
	if err := dbus.DbusInstance.NewDBus(); err != nil {
		t.Skipf("Error initializating D-Bus system. Stop the program. FATAL: %v", err)
	}
	if err := dbus.DbusInstance.StartUnit(test.Unit); err != nil {
		os.Remove(unitServiceLink)
		os.Remove(unitServicePath)
		t.Skipf("Can not start Unit. %v", err)
	}
}
