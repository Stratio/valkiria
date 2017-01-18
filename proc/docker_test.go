package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"testing"
)

const (
	killNameDocker = "mesos_1232134"
)

var (
	testKillDockers = func(t *testing.T) {
		var d = Docker{KillName: killNameDocker}
		err := d.Kill()
		if err != nil {
			t.Errorf("proc.testKillDaemon - ERROR: %v", err)
		}
	}
)

func TestDockerLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDockerTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testKillDockers", testKillDockers)
}





