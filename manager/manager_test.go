package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"testing"
	"github.com/Stratio/valkiria/dbus"
	"os"
)

const(
	unit 		= "test.service"
	fakeUnit        = "fakeUnit.service"
	unitServicePath = "/tmp/test.service"
	unitServiceLink = "/lib/systemd/system/test.service"
)

var (
	testManagerRead = func(t *testing.T) {
		var manager = NewManager()
		res, err := manager.Read()
		if err != nil  {
			t.Fatal("Error reading process %v", err)
		}
		if len(res) < 1{
			t.Fatal("Error reading process")
		}
	}
	testManagerKillFake = func(t *testing.T) {
		var manager = NewManager()
		res, err := manager.Shooter(fakeUnit, "killExecutor=0")
		if len(err) > 0 {
			t.Fatal("Error reading process %v", err)
		}
		if len(res) > 0{
			t.Fatal("Error reading process")
		}
	}
	testManagerKillService = func(t *testing.T) {
		var manager = NewManager()
		res, err := manager.Shooter("mesos-32156487435168416831", "killExecutor=2")
		if len(err) > 0 {
			t.Fatal("Error reading process %v", err)
		}
		if len(res) < 1{
			t.Fatal("Error reading process")
		}
	}
)

func TestManagerLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	test.SetupDockerTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testManagerRead", testManagerRead)
	t.Run("testManagerKillFake", testManagerKillFake)
	t.Run("testManagerKillService", testManagerKillService)
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