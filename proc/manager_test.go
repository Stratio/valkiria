package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"strings"
	"testing"
)

var (
	testLoadProcesses = func(t *testing.T) {
		var manager = new(Manager)
		manager.daemonConfigString = []string{test.Unit}
		manager.dockerConfigPattern = "^/testValkiria"
		eLoad := manager.LoadProcesses()
		if eLoad != nil {
			t.Fatalf("proc.testReadAllService - ERROR: %v", eLoad)
		}
		if len(manager.Daemons) != 1 || len(manager.Dockers) != 1 || len(manager.Services) != 1 {
			t.Fatalf("proc.testReadAllService - ERROR: Should have 1 element by slice")
		}
		if !strings.EqualFold(manager.Daemons[0].Name, test.Unit) || !strings.EqualFold(manager.Dockers[0].Name, "/"+test.DockerContainerName) || !strings.EqualFold(manager.Services[0].TaskName, test.MesosName) {
			t.Fatalf("proc.testReadAllService - ERROR: Name of daemon, docker or service does not match")
		}

	}
	testShooter = func(t *testing.T) {
		var manager = new(Manager)
		manager.daemonConfigString = []string{test.Unit}
		manager.blackListServices = []string{}
		manager.dockerConfigPattern = "^/testValkiria"
		eLoad := manager.LoadProcesses()
		if eLoad != nil {
			t.Errorf("proc.testShooter - ERROR: %v", eLoad)
		}
		if len(manager.Daemons) != 1 || len(manager.Dockers) != 1 || len(manager.Services) != 1 {
			t.Errorf("proc.testShooter - ERROR: Should have 1 element by slice")
		}
		if !strings.EqualFold(manager.Daemons[0].Name, test.Unit) || !strings.EqualFold(manager.Dockers[0].Name, "/"+test.DockerContainerName) || !strings.EqualFold(manager.Services[0].TaskName, test.MesosName) {
			t.Errorf("proc.testShooter - ERROR: Name of daemon, docker or service does not match")
		}
		rShoterService, eShoterService := manager.Shooter(manager.Services[0].TaskName, serviceEnum, true)
		if eShoterService != nil {
			t.Errorf("proc.testShooter - ERROR: %v", eShoterService)
		}
		if rShoterService == nil{
			t.Errorf("proc.testShooter - ERROR: service can be not nil")
		}
		rShoterDocker, eShoterDocker := manager.Shooter(manager.Dockers[0].TaskName, dockerEnum, false)
		if eShoterDocker != nil {
			t.Errorf("proc.testShooter - ERROR: %v", eShoterDocker)
		}
		if rShoterDocker == nil{
			t.Errorf("proc.testShooter - ERROR: docker can be not nil")
		}
		rShoterDaemon, eShoterDaemon := manager.Shooter(manager.Daemons[0].Name, daemonEnum, false)
		if eShoterDaemon != nil {
			t.Errorf("proc.testShooter - ERROR: %v", eShoterDaemon)
		}
		if rShoterDaemon == nil{
			t.Errorf("proc.testShooter - ERROR: daemon can be not nil")
		}
		rShoterSearch, eShoterSearch := manager.Shooter(manager.Daemons[0].Name, searchTypeEnum, false)
		if eShoterSearch != nil {
			t.Errorf("proc.testShooter - ERROR: %v", eShoterSearch)
		}
		if rShoterSearch != nil{
			t.Errorf("proc.testShooter - ERROR: daemon can be nil")
		}
		_, eShoterDefault := manager.Shooter(manager.Daemons[0].Name, 10, false)
		if eShoterDefault == nil {
			t.Errorf("proc.testShooter - ERROR: Should be an error")
		}
	}
	testShooterFake = func(t *testing.T) {
		var manager = new(Manager)
		manager.Daemons = []Daemon{Daemon{Name: "fakeDaemon"}}
		manager.Dockers = []Docker{Docker{Id: "12321", TaskName: "fakeDocker"}}
		manager.Services = []Service{Service{Pid: 99999, TaskName: "fakeService"}}

		rShoterService, eShoterService := manager.Shooter(manager.Services[0].TaskName, serviceEnum, false)
		if eShoterService == nil {
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterService != nil{
			t.Errorf("proc.testShooter - ERROR: service can be nil")
		}
		rShoterServiceTrue, eShoterServiceTrue := manager.Shooter(manager.Services[0].TaskName, serviceEnum, true)
		if eShoterServiceTrue == nil {
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterServiceTrue != nil{
			t.Errorf("proc.testShooter - ERROR: service can be nil")
		}
		rShoterDocker, eShoterDocker := manager.Shooter(manager.Dockers[0].TaskName, dockerEnum, false)
		if eShoterDocker == nil {
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterDocker != nil{
			t.Errorf("proc.testShooter - ERROR: docker can be nil")
		}
		rShoterDaemon, eShoterDaemon := manager.Shooter(manager.Daemons[0].Name, daemonEnum, false)
		if eShoterDaemon == nil {
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterDaemon != nil{
			t.Errorf("proc.testShooter - ERROR: daemon can be nil")
		}
	}
)

func TestManagerChaosLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	test.SetupDockerTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testLoadProcesses", testLoadProcesses)
}

func TestManagerShooterLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDBusTest(t)
	test.SetupDockerTest(t)
	startDBusUnit(t)
	defer test.TearDownDBusTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testShooter", testShooter)
	t.Run("testShooterFake", testShooterFake)
}
