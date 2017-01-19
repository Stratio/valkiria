package workers

import (
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/plugin/mesos"
)

var (
	isStart bool
	WorkerProcessInstance *WorkerProcess
)

type ProcKill struct {
	Name string
	Properties string
}

type Msg struct{
	ProcessKill ProcKill
	Proc []plugin.Process
	Err []error
	Res chan interface{}
}

type WorkerProcess struct{
	PluginInstance plugin.Plugin
	GetDaemons  func ()([]plugin.Process, error)
	GetServices func ()([]plugin.Process, error)
	GetDockers  func ()([]plugin.Process, error)
	GetFindAndKill func ([]plugin.Process, string, string)([]plugin.Process, []error)
	ChRead      chan *Msg
	ChKill      chan *Msg
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
	p.GetDaemons = plug.GetDaemons()
	p.GetServices = plug.GetServices()
	p.GetDockers = plug.GetDocker()
	p.GetFindAndKill = plug.FindAndKill()
}

func (p *WorkerProcess) LoadProcesses() (res []plugin.Process, err []error) {
	daemons, errDaemons := p.GetDaemons()
	if errDaemons != nil {
		err = append(err, errDaemons)
	}
	services, errServices := p.GetServices()
	if errServices != nil {
		err = append(err, errServices)
	}
	dockers, errDockers := p.GetDockers()
	if errDockers != nil {
		err = append(err, errDockers)
	}
	res = append(append(append(res, daemons...), services...), dockers...)
	return
}

func (p *WorkerProcess) WorkerProcess(plug plugin.Plugin){
	for{
		select {
		case msgRead := <-p.ChRead:
			msgRead.Proc, msgRead.Err = p.LoadProcesses()
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

		case msgKill := <- p.ChKill:
			process, _ := p.LoadProcesses()
			msgKill.Proc, msgKill.Err = p.GetFindAndKill(process, msgKill.ProcessKill.Name, msgKill.ProcessKill.Properties)
			msgKill.Res <- true
			for {
				select{
				case msgNew := <- p.ChKill:
					msgNew.Proc, msgNew.Err = p.GetFindAndKill(process, msgNew.ProcessKill.Name, msgNew.ProcessKill.Properties)
					msgNew.Res <- true
				default:
					break
				}
				if len(p.ChKill) < 1 {
					break
				}
			}
		}
	}
}