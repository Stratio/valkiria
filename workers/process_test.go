package workers

import (
	log "github.com/Sirupsen/logrus"
	"testing"
	"github.com/Stratio/valkiria/test"
)

var (
	testProcessWorkerRead = func(t *testing.T) {
		var processWorker = NewWorkerProcess()
		processWorker.LoadProcesses()
		var msgRead = Msg{Res: make(chan interface{}, 1)}
		var msgKill = Msg{Res: make(chan interface{}, 1)}
		processWorker.ChRead <- &msgRead
		processWorker.ChRead <- &msgRead
		<- msgRead.Res
		if len(msgRead.Err) > 1 {
			t.Errorf("Should be 0 errors")
		}
		processWorker.ChKill <- &msgKill
		processWorker.ChKill <- &msgKill
		<- msgKill.Res
		if len(msgKill.Err) > 1 {
			t.Errorf("Should be 0 errors")
		}	}
	testProcessWorkerKill = func(t *testing.T) {
		var processWorker = NewWorkerProcess()
		processWorker.LoadProcesses()
		var msgRead = Msg{Res: make(chan interface{}, 1)}
		var msgKill = Msg{Res: make(chan interface{}, 1)}
		processWorker.ChKill <- &msgKill
		processWorker.ChKill <- &msgKill
		<- msgKill.Res
		if len(msgKill.Err) > 1 {
			t.Errorf("Should be 0 errors")
		}

		processWorker.ChRead <- &msgRead
		processWorker.ChRead <- &msgRead
		<- msgRead.Res
		if len(msgRead.Err) > 1 {
			t.Errorf("Should be 0 errors")
		}

	}
)

func TestWorkerReadLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	test.SetupDockerTest(t)
	defer test.TearDownDockerTest(t)
	defer test.TearDownDBusTest(t)
	t.Run("testProcessWorkerKill", testProcessWorkerKill)
	t.Run("testProcessWorkerRead", testProcessWorkerRead)
}