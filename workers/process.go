package workers

import (
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/plugin/mesos"
)

var (
	isStart               bool
	WorkerProcessInstance *WorkerProcess
)

type ProcKill struct {
	Name       string
	Properties string
}

type Msg struct {
	ProcessKill ProcKill
	Proc        []plugin.Process
	Err         []error
	Res         chan interface{}
}

type WorkerProcess struct {
	PluginInstance 	plugin.Plugin
	GetIps     	func() ([]string, error)
	ChRead         	chan *Msg
	ChKill         	chan *Msg
}

func NewWorkerProcess() *WorkerProcess {
	if !isStart {
		isStart = true
		WorkerProcessInstance = new(WorkerProcess)
		WorkerProcessInstance.ChRead = make(chan *Msg, 100)
		WorkerProcessInstance.ChKill = make(chan *Msg, 100)
		//TODO: change mesos config, set by command line
		WorkerProcessInstance.PluginInstance = mesos.NewMesosConfig()
		WorkerProcessInstance.ConfigManager(WorkerProcessInstance.PluginInstance)
		go WorkerProcessInstance.WorkerProcess(WorkerProcessInstance.PluginInstance)
	}
	return WorkerProcessInstance
}

func (p *WorkerProcess) ConfigManager(plug plugin.Plugin) {

}

func (p *WorkerProcess) WorkerProcess(plug plugin.Plugin) {
	for {
		select {
		case msgRead := <-p.ChRead:
			msgRead.Res <- true
			//. high concurrency purpose
			for {
				select {
				case msgNew := <-p.ChRead:
					msgNew.Proc = msgRead.Proc
					msgNew.Err = msgRead.Err
					msgNew.Res <- true
				default:
					break
				}
				if len(p.ChRead) < 1 {
					break
				}
			}

		}
	}
}
