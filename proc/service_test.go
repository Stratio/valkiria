package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"testing"

)

const(
	killNameFake = "mesos-32156487435168416831"
)

var (
	testKillService = func(t *testing.T) {
		var d = Service{KillName: killNameFake, Pid:123456}
		err := d.Kill()
		if err == nil {
			t.Errorf("proc.testKillDaemon - Should be an error")
		}
	}
)

func TestServiceLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	t.Run("testKillService", testKillService)
}


