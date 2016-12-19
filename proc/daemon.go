package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/dbus"
	"time"
)

type Daemon struct {
	Pid            uint32
	Name           string
	Path           string
	ChaosTimeStamp int64
}

func (d *Daemon) Kill() (err error) {
	d.ChaosTimeStamp = time.Now().UTC().UnixNano()
	log.Debugf("proc.daemon.Kill - killing '%v' '%v' '%v' '%v'", d.Pid, d.Name, d.Path, d.ChaosTimeStamp)
	err = dbus.DbusInstance.KillUnit(d.Name)
	if err != nil {
		log.Errorf("proc.daemon.Kill - '%v' '%v' '%v' '%v' - '%v'", d.Pid, d.Name, d.Path, d.ChaosTimeStamp, err.Error())
	}
	return
}

func ReadAllDaemons(listDaemons []string) (res []Daemon, err error) {
	log.Debug("proc.daemon.ReadAllDaemons")
	for _, d := range listDaemons {
		if path, err := dbus.DbusInstance.GetUnit(d); err == nil {
			if pid, _ := dbus.DbusInstance.GetUnitPid(path); err == nil {
				res = append(res, Daemon{Pid: pid, Name: d, Path: path})
				log.Debugf("proc.daemon.ReadAllDaemons - append - '%v' '%v' '%v'", d, pid, path)
			}
		} else {
			log.Infof("proc.daemon.ReadAllDaemons - ERROR: '%v'", err.Error())
		}
	}
	log.Debugf("proc.daemon.ReadAllDaemons - lenDaemon: '%v'", len(res))
	return
}
