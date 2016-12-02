package dbus

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"os"
	"testing"
)

const (
	unit                    = "test.service"
	fakeUnit                = "fakeUnit.service"
	fakePath                = "/fake/path"
	fakePathNotRegexCompile = "fakepath"
	unitServicePath         = "/tmp/test.service"
	unitServiceLink         = "/lib/systemd/system/test.service"
)

var (
	testInitDefaultSharedObjectService = func(t *testing.T) {
		var err = DbusInstance.NewDBus()
		if err != nil {
			t.Errorf("dbus.TestinitDefaultSharedObjectService - ERROR: %v", err)
		}
	}
	testGetMachineId = func(t *testing.T) {
		_, err := DbusInstance.GetMachineId()
		if err != nil {
			t.Errorf("dbus.TestGetMachineId - ERROR: %v", err)
		}
	}
	testStartUnit = func(t *testing.T) {
		err := DbusInstance.StartUnit(unit)
		if err != nil {
			t.Errorf("dbus.TestStartUnit - ERROR: %v", err)
		}
		errFake := DbusInstance.StartUnit(fakeUnit)
		if errFake == nil {
			t.Errorf("dbus.TestStartUnit - ERROR: fakeUnit not exist but result is succes")
		}
	}
	testGetUnit = func(t *testing.T) {
		_, err := DbusInstance.GetUnit(unit)
		if err != nil {
			t.Errorf("dbus.TestGetUnit - ERROR: %v", err)
		}
		_, errFake := DbusInstance.GetUnit(fakeUnit)
		if errFake == nil {
			t.Errorf("dbus.TestGetUnit - ERROR: fakeUnit not exist but result is succes")
		}
	}
	testGetUnitPid = func(t *testing.T) {
		path, errGetUnit := DbusInstance.GetUnit(unit)
		if errGetUnit != nil {
			t.Errorf("dbus.TestGetUnitPid - ERROR: %v", errGetUnit)
		}
		_, err := DbusInstance.GetUnitPid(path)
		if err != nil {
			t.Errorf("dbus.TestGetUnitPid - ERROR: %v", err)
		}
		_, errFake := DbusInstance.GetUnitPid(fakePath)
		if errFake == nil {
			t.Errorf("dbus.TestGetUnitPid - ERROR: /fake/path not exist but result is succes")
		}
		_, errFakePathRegexCompile := DbusInstance.GetUnitPid(fakePathNotRegexCompile)
		if errFakePathRegexCompile == nil {
			t.Errorf("dbus.TestGetUnitPid - ERROR: fake/path not regex compile")
		}
	}
	testKillUnit = func(t *testing.T) {
		err := DbusInstance.KillUnit(unit)
		if err != nil {
			t.Errorf("dbus.TestKillUnit - ERROR: %v", err)
		}
		errFake := DbusInstance.KillUnit(fakeUnit)
		if errFake == nil {
			t.Errorf("dbus.TestKillUnit - ERROR: fakeUnit not exist but result is succes")
		}
	}
	testStopUnit = func(t *testing.T) {
		err := DbusInstance.StopUnit(unit)
		if err != nil {
			t.Errorf("dbus.TestStopUnit - ERROR: %v", err)
		}
		errFake := DbusInstance.StopUnit(fakeUnit)
		if errFake == nil {
			t.Errorf("dbus.TestStopUnit - ERROR: fakeUnit not exist but result is succes")
		}
	}
)

func TestDBusLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	defer test.TearDownDBusTest(t)
	startDBusUnit(t)
	t.Run("testInitDefaultSharedObjectService", testInitDefaultSharedObjectService)
	t.Run("testGetMachineId", testGetMachineId)
	t.Run("testStartUnit", testStartUnit)
	t.Run("testGetUnit", testGetUnit)
	t.Run("testGetUnitPid", testGetUnitPid)
	t.Run("testKillUnit", testKillUnit)
	t.Run("testStopUnit", testStopUnit)
}

//WARNING: this code can produce import cycle not allowed in test. dont move to test_util
// If it is necessary to make changes in the tests, make changes in all the tests impacted
func startDBusUnit(t *testing.T) {
	if err := DbusInstance.NewDBus(); err != nil {
		t.Skipf("Error initializating D-Bus system. Stop the program. FATAL: %v", err)
	}
	if err := DbusInstance.StartUnit(unit); err != nil {
		os.Remove(unitServiceLink)
		os.Remove(unitServicePath)
		t.Skipf("Can not start Unit. %v", err)
	}
}
