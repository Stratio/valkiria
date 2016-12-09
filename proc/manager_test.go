package proc

import (
	"github.com/Stratio/valkiria/test"
	log "github.com/Sirupsen/logrus"
	"testing"
	"strings"
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
		if !strings.EqualFold(manager.Daemons[0].Name, test.Unit) || !strings.EqualFold(manager.Dockers[0].Name, "/"+test.DockerContainerName) || !strings.EqualFold(manager.Services[0].TaskName, test.MesosName){
			t.Fatalf("proc.testReadAllService - ERROR: Name of daemon, docker or service does not match")
		}

	}
	testChaos = func(t *testing.T) {
		var manager = new(Manager)
		manager.daemonConfigString = []string{test.Unit}
		manager.blackListServices = []string{}
		manager.dockerConfigPattern = "^/testValkiria"
		eLoad := manager.LoadProcesses()
		if eLoad != nil {
			t.Fatalf("proc.testChaos - ERROR: %v", eLoad)
		}
		if len(manager.Daemons) != 1 || len(manager.Dockers) != 1 || len(manager.Services) != 1 {
			t.Fatalf("proc.testReadAllService - ERROR: Should have 1 element by slice")
		}
		if !strings.EqualFold(manager.Daemons[0].Name, test.Unit) || !strings.EqualFold(manager.Dockers[0].Name, "/"+test.DockerContainerName) || !strings.EqualFold(manager.Services[0].TaskName, test.MesosName){
			t.Fatalf("proc.testReadAllService - ERROR: Name of daemon, docker or service does not match")
		}
		eChaosService := manager.Chaos(0,0,1)
		if eChaosService != nil{
			t.Fatalf("proc.testChaos - ERROR: %v", eChaosService)
		}
		if len(manager.Services) != 0 {
			t.Fatalf("proc.testChaos - ERROR: Slice dockers length can be 0")
		}
		eChaosDaemon := manager.Chaos(1,0,0)
		if eChaosDaemon != nil{
			t.Fatalf("proc.testChaos - ERROR: %v", eChaosDaemon)
		}
		if len(manager.Daemons) != 0 {
			t.Fatalf("proc.testChaos - ERROR: Slice daemons length can be 0")
		}
		eChaosDocker := manager.Chaos(0,1,0)
		if eChaosDocker != nil{
			t.Fatalf("proc.testChaos - ERROR: %v", eChaosDocker)
		}
		if len(manager.Dockers) != 0 {
			t.Fatalf("proc.testChaos - ERROR: Slice dockers length can be 0")
		}
	}
	testChaosFake = func(t *testing.T) {
		var manager = new(Manager)
		manager.Daemons = []daemon{daemon{Name:"fakeDaemon"}}
		manager.Dockers = []docker{docker{Id:"12321", TaskName:"fakeDocker"}}
		manager.Services = []service{service{Pid: 99999, TaskName:"fakeService"}}

		eChaosService := manager.Chaos(0,0,1)
		if eChaosService == nil{
			t.Fatalf("proc.testChaos - ERROR: Should be error")
		}
		eChaosDaemon := manager.Chaos(1,0,0)
		if eChaosDaemon == nil{
			t.Fatalf("proc.testChaos - ERROR: Should be error")
		}
		eChaosDocker := manager.Chaos(0,1,0)
		if eChaosDocker == nil{
			t.Fatalf("proc.testChaos - ERROR: Should be error")
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
		if !strings.EqualFold(manager.Daemons[0].Name, test.Unit) || !strings.EqualFold(manager.Dockers[0].Name, "/"+test.DockerContainerName) || !strings.EqualFold(manager.Services[0].TaskName, test.MesosName){
			t.Errorf("proc.testShooter - ERROR: Name of daemon, docker or service does not match")
		}
		rShoterService, eShoterService := manager.Shooter(manager.Services[0].TaskName, serviceEnum, true)
		if eShoterService != nil{
			t.Errorf("proc.testShooter - ERROR: %v", eShoterService)
		}
		if !rShoterService {
			t.Errorf("proc.testShooter - ERROR: Bool service can be true")
		}
		rShoterDocker, eShoterDocker := manager.Shooter(manager.Dockers[0].TaskName, dockerEnum, false)
		if eShoterDocker != nil{
			t.Errorf("proc.testShooter - ERROR: %v", eShoterDocker)
		}
		if !rShoterDocker {
			t.Errorf("proc.testShooter - ERROR: Bool docker can be true")
		}
		rShoterDaemon, eShoterDaemon := manager.Shooter(manager.Daemons[0].Name, daemonEnum, false)
		if eShoterDaemon != nil{
			t.Errorf("proc.testShooter - ERROR: %v", eShoterDaemon)
		}
		if !rShoterDaemon {
			t.Errorf("proc.testShooter - ERROR: Bool daemon can be true")
		}
		rShoterSearch, eShoterSearch := manager.Shooter(manager.Daemons[0].Name, searchTypeEnum, false)
		if eShoterSearch != nil{
			t.Errorf("proc.testShooter - ERROR: %v", eShoterSearch)
		}
		if rShoterSearch {
			t.Errorf("proc.testShooter - ERROR: Bool daemon default can be false")
		}
		_, eShoterDefault := manager.Shooter(manager.Daemons[0].Name, 10, false)
		if eShoterDefault == nil{
			t.Errorf("proc.testShooter - ERROR: Should be an error")
		}
	}
	testShooterFake = func(t *testing.T) {
		var manager = new(Manager)
		manager.Daemons = []daemon{daemon{Name:"fakeDaemon"}}
		manager.Dockers = []docker{docker{Id:"12321", TaskName:"fakeDocker"}}
		manager.Services = []service{service{Pid: 99999, TaskName:"fakeService"}}

		rShoterService, eShoterService := manager.Shooter(manager.Services[0].TaskName, serviceEnum, false)
		if eShoterService == nil{
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterService {
			t.Errorf("proc.testShooter - ERROR: Bool service can be false")
		}
		rShoterServiceTrue, eShoterServiceTrue := manager.Shooter(manager.Services[0].TaskName, serviceEnum, true)
		if eShoterServiceTrue == nil{
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterServiceTrue {
			t.Errorf("proc.testShooter - ERROR: Bool service can be false")
		}
		rShoterDocker, eShoterDocker := manager.Shooter(manager.Dockers[0].TaskName, dockerEnum, false)
		if eShoterDocker == nil{
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterDocker {
			t.Errorf("proc.testShooter - ERROR: Bool docker can be false")
		}
		rShoterDaemon, eShoterDaemon := manager.Shooter(manager.Daemons[0].Name, daemonEnum, false)
		if eShoterDaemon == nil{
			t.Errorf("proc.testShooter - ERROR: Shoould be error")
		}
		if rShoterDaemon {
			t.Errorf("proc.testShooter - ERROR: Bool daemon can be false")
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
	t.Run("testChaos", testChaos)
	t.Run("testChaosFake", testChaosFake)
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