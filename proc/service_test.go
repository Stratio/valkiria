package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"strings"
	"testing"
)

var (
	testKillService = func(t *testing.T) {
		rRead, eRead := ReadAllDaemons([]string{test.Unit})
		if eRead != nil {
			t.Fatalf("proc.testReadAllService - ERROR: %v", eRead)
		}
		rProc, eProc := ReadAllChildProcess(rRead, []string{"test.service"}, []string{})
		if eProc != nil {
			t.Fatalf("proc.testReadAllService - ERROR: %v", eProc)
		}
		if len(rProc) < 1 {
			t.Fatalf("proc.testReadAllService - ERROR: It should have almost 1 element.")
		}
		if !strings.EqualFold(test.MesosName, rProc[0].TaskName) {
			t.Fatalf("proc.testReadAllService - ERROR: Does not match the task recovered with the test.")
		}
		eKill := rProc[0].Kill()
		if eKill != nil {
			t.Fatalf("proc.testReadAllService - ERROR: %v", eKill)
		}
	}
	testReadAllService = func(t *testing.T) {
		rRead, eRead := ReadAllDaemons([]string{test.Unit})
		if eRead != nil {
			t.Fatalf("", eRead)
		}
		rProc, eProc := ReadAllChildProcess(rRead, []string{"test.service"}, []string{})
		if eProc != nil {
			t.Fatalf("proc.testReadAllService - ERROR: %v", eProc)
		}
		if len(rProc) < 1 {
			t.Fatalf("proc.testReadAllService - ERROR: It should have almost 1 element.")
		}
		if !strings.EqualFold(test.MesosName, rProc[0].TaskName) {
			t.Fatalf("proc.testReadAllService - ERROR: Does not match the task recovered with the test.")
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
	t.Run("testReadAllService", testReadAllService)
	t.Run("testKillService", testKillService)
	t.Run("testIsInBlackList", testIsInBlackList)

}
