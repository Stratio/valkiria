package dbus

import (
	"github.com/godbus/dbus"
	"os"
	log "github.com/Sirupsen/logrus"
	"regexp"
	"github.com/pkg/errors"
)

const (
	service = "org.freedesktop.systemd1"
	object = "/org/freedesktop/systemd1"
	getMachineId = "org.freedesktop.DBus.Peer.GetMachineId"
	getUnit = "org.freedesktop.systemd1.Manager.GetUnit"
	startUnit = "org.freedesktop.systemd1.Manager.StartUnit"
	stopUnit = "org.freedesktop.systemd1.Manager.StopUnit"
	killUnit = "org.freedesktop.systemd1.Manager.KillUnit"
	getUnitproperties = "org.freedesktop.systemd1.Service.MainPID"
	stopMode = "replace"
	killMode = "main"
	startMode = "replace"
	killSignall int32 = 9
	systemdPathRegex = `^/[a-zA-Z0-9/_\-.]*$`
)

var (
	DbusInstance DSbusStruct
)

func init(){
	err := initDefaultSharedObjectService()
	if err != nil {
		log.Fatalf("dbus.init - Error initializating D-Bus system. Stop the program. FATAL: %v", err)
		os.Exit(-1)
	}

}

type DSbusStruct struct {
	c *dbus.Conn
	o dbus.BusObject
}

func initDefaultSharedObjectService () (err error){
	log.Debug("dbus.initDefaultSharedObjectService")
	if DbusInstance.c, err = dbus.SystemBus(); err != nil{
		log.Infof("dbus.initDefaultSharedObjectService - ERROR: %v", err)
	} else {
		DbusInstance.o = DbusInstance.c.Object(service, object)
	}
	log.Debugf("dbus.initDefaultSharedObjectService - RES: %v", DbusInstance.c.Names())
	return
}

func (d *DSbusStruct) GetMachineId ()(r string, err error){
	log.Debug("dbus.GetMachineId")
	err = d.o.Call(getMachineId, 0).Store(&r)
	if err != nil {
		log.Infof("dbus.GetMachineId - ERROR: %v", err.Error())
	}
	log.Debugf("dbus.GetMachineId - RES: %v", r)
	return
}

// replace, fail, isolate, ignore-dependencies, ignore-requirements
func (d *DSbusStruct) StopUnit (unitName string) (err error){
	log.Debug("dbus.StopUnit")
	var path dbus.ObjectPath
	err = d.o.Call(stopUnit, 0, unitName, stopMode).Store(&path)
	if err != nil {
		log.Infof("dbus.StopUnit - ERROR: %v", err.Error())
	}
	log.Debugf("dbus.StopUnit - RES: %v", path)
	return
}

func (d *DSbusStruct) KillUnit (unitName string) (err error){
	log.Debug("dbus.KillUnit")
	err = d.o.Call(killUnit, 0, unitName, killMode, killSignall).Store()
	if err != nil {
		log.Infof("dbus.KillUnit - ERROR: %v", err.Error())
	}
	log.Debug("dbus.dbus.KillUnit - RES: EMPTY RES")
	return
}

// replace, fail, isolate, ignore-dependencies, ignore-requirements
func (d *DSbusStruct) StartUnit (unitName string) (err error){
	log.Debug("dbus.StartUnit")
	var path dbus.ObjectPath
	if err = d.o.Call(startUnit, 0, unitName, startMode).Store(&path); err != nil {
		log.Infof("dbus.StartUnit - ERROR: %v", err)
	}
	log.Debugf("dbus.StartUnit - RES: %v", path)
	return
}

func (d *DSbusStruct) GetUnit (unitName string)(res string, err error){
	log.Debug("dbus.GetUnit")
	var unitPath dbus.ObjectPath
	conn := d.o.Call(getUnit, 0, unitName)
	err = conn.Store(&unitPath)
	if err != nil {
		log.Infof("dbus.GetUnit - ERROR: %v", err.Error())
	} else {
		res = string(unitPath)
	}
	log.Debugf("dbus.GetUnit - RES: %v", res)
	return
}

func (d *DSbusStruct) GetUnitPid (unitPath string)(res uint32, err error){
	log.Debug("dbus.GetUnitPid")
	var validPath = regexp.MustCompile(systemdPathRegex)
	if b := validPath.Match([]byte(unitPath)); b {
		obj := d.c.Object(service, dbus.ObjectPath(unitPath))
		variant, err := obj.GetProperty(getUnitproperties)
		if err != nil {
			log.Infof("dbus.GetUnitPid - ERROR: %v", err.Error())
			return res, err
		} else {
			res = variant.Value().(uint32)
		}
	} else {
		log.Infof("dbus.GetUnitPid - ERROR: Unit path can not compile with regex.")
		err = errors.New("dbus.GetUnitPid - ERROR: Invalid unit path: " + unitPath + ". Can not compile with validation expression.")
	}
	log.Debugf("dbus.GetUnitPid - RES: %v", res)
	return
}
