package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/dbus"
	"github.com/Stratio/valkiria/test"
	"os"
	"testing"
)

const (
	fakeUnit        = "fakeUnit.service"
	unitServicePath = "/tmp/test.service"
	unitServiceLink = "/lib/systemd/system/test.service"
)

var (
	testKillDaemon = func(t *testing.T) {
		var d = Daemon{KillName: test.Unit}
		err := d.Kill()
		if err != nil {
			t.Errorf("proc.testKillDaemon - ERROR: %v", err)
		}
		var dFake = Daemon{KillName: fakeUnit}
		errFake := dFake.Kill()
		if errFake == nil {
			t.Error("proc.testKillDaemon - ERROR: fakeUnit not exist but result is succes")
		}
	}
)

func TestDaemonLib(t *testing.T) {
	test.SetupDBusTest(t)
	defer test.TearDownDBusTest(t)
	log.SetLevel(test.Level)
	startDBusUnit(t)
	t.Run("testKill", testKillDaemon)
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
