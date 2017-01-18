package mesos

import(
	"github.com/Stratio/valkiria/proc"
	"github.com/valkiria/dbus"
	"github.com/Stratio/valkiria/plugin"
)

func (m *MesosConfig) GetDaemons() (func ()([]plugin.Process, error)){
	return func ()([]plugin.Process, error){
		return ReadAllDaemons(m.DaemonConfigString)
	}
}

func ReadAllDaemons(listDaemons []string) (res []plugin.Process, err error) {
	for _, d := range listDaemons {
		if path, err := dbus.DbusInstance.GetUnit(d); err == nil {
			if pid, _ := dbus.DbusInstance.GetUnitPid(path); err == nil {
				res = append(res, &proc.Daemon{Pid: pid, KillName: d})
			}
		}
	}
	return
}