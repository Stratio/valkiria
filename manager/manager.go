package manager

import (
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/workers"
)

var (
	isStart         bool
	ManagerInstance *Manager
)

type Manager struct {
	WorkerProcessInstance *workers.WorkerProcess
}

func NewManager() *Manager {
	if !isStart {
		isStart = true
		ManagerInstance = new(Manager)
		ManagerInstance.WorkerProcessInstance = workers.NewWorkerProcess()
	}
	return ManagerInstance
}

func (m *Manager) Read() ([]plugin.Process, []error) {
	//.
	var msgRead = workers.Msg{Res: make(chan interface{}, 1)}
	defer close(msgRead.Res)
	//.
	ManagerInstance.WorkerProcessInstance.ChRead <- &msgRead
	<-msgRead.Res
	return msgRead.Proc, msgRead.Err
}

func (m *Manager) Shooter(name string, properties string) ([]plugin.Process, []error) {
	//.
	var msgKill = workers.Msg{Res: make(chan interface{}, 1), ProcessKill: workers.ProcKill{Name: name, Properties: properties}}
	defer close(msgKill.Res)
	//.
	ManagerInstance.WorkerProcessInstance.ChKill <- &msgKill
	<-msgKill.Res
	return msgKill.Proc, msgKill.Err
}
