package dbus

import (
	"github.com/godbus/dbus"
	"fmt"
	"os"
	log "github.com/Sirupsen/logrus"
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
	killSignall int32 = 9
)

var (
	DbusInstance DSbusStruct
)

func init(){
	err := initDefaultSharedObjectService()
	if err != nil {
		fmt.Println("Error initializating D-Bus system. Stop the program.")
		os.Exit(-1)
	}

}

type DSbusStruct struct {
	c *dbus.Conn
	o dbus.BusObject
}

func initDefaultSharedObjectService () (err error){
	DbusInstance.c, err = dbus.SystemBus()
	DbusInstance.o = DbusInstance.c.Object(service, object)
	return
}

func (d *DSbusStruct) GetMachineId ()(r string, err error){
	log.Debug("dbus.dbus.GetMachineId")
	err = d.o.Call(getMachineId, 0).Store(&r)
	if err != nil {
		log.Debug("dbus.dbus.GetMachineId - ERROR: " + err.Error())
	}
	return
}

//"replace" "fail" "isolate" "ignore-dependencies" "ignore-requirements"
func (d *DSbusStruct) StopUnit (unitName string) (err error){
	log.Debug("dbus.dbus.StopUnit")
	var path dbus.ObjectPath
	err = d.o.Call(stopUnit, 0, unitName, stopMode).Store(&path)
	if err != nil {
		log.Debug("dbus.dbus.StopUnit - ERROR: " + err.Error())
	}
	return
}

func (d *DSbusStruct) KillUnit (unitName string) (err error){
	log.Debug("dbus.dbus.KillUnit")
	err = d.o.Call(killUnit, 0, unitName, killMode, killSignall).Store()
	if err != nil {
		log.Debug("dbus.dbus.KillUnit - ERROR: " + err.Error())
	}
	return
}

func (d *DSbusStruct) StartUnit (unitName string) (err error){
	var path dbus.ObjectPath
	err = d.o.Call(startUnit, 0, unitName, killMode).Store(&path)
	return
}

func (d *DSbusStruct) GetUnit (unitName string)(res string, err error){
	var unitPath dbus.ObjectPath
	conn := d.o.Call(getUnit, 0, unitName)
	err = conn.Store(&unitPath)
	res = string(unitPath)
	return
}

func (d *DSbusStruct) GetUnitPid (unitPath string)(res uint32, err error){
	obj := d.c.Object(service, dbus.ObjectPath(unitPath))
	variant, err := obj.GetProperty(getUnitproperties)
	res = variant.Value().(uint32)
	return
}
