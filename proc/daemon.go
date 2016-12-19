package proc

import (
	"github.com/Stratio/valkiria/dbus"
	"time"
)

type Daemon struct {
	Pid            uint32
	KillName       string
	Path           string
	ChaosTimeStamp int64
}

func (d *Daemon) Kill() (err error) {
	err = dbus.DbusInstance.KillUnit(d.KillName)
	if err != nil {
		d.ChaosTimeStamp = time.Now().UTC().UnixNano()
	}
	return
}


