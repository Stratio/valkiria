package test

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/Stratio/valkiria/dbus"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"testing"
)

const (
	Level = log.DebugLevel
)

const (
	Unit            = "test.service"
	uid             = "0"
	unitServicePath = "/tmp/test.service"
	unitServiceLink = "/lib/systemd/system/test.service"
	MesosName       = "mesos-32156487435168416831"
	pathTest        = "/test/test/test/test/test/test/test/test/test/" + MesosName + "/test/test/test"
)

const (
	DockerImage         = "ubuntu"
	DockerContainerName = "testValkiria"
)

var (
	testFile = []byte("[Unit]\nDescription=test\n\n[Service]\nExecStart=/bin/bash -c 'while true; do echo hello; sleep 2; done'\nRestart=always\nWorkingDirectory=" + pathTest + "\n")
)

func SetupDBusTest(t *testing.T) {
	SkipTesAllDBusTest(t)
	AddTestServiceDBusTest(t)
}

func TearDownDBusTest(t *testing.T) {
	if err := os.Remove(unitServiceLink); err != nil {
		t.Errorf("Can not delete test link. %v", err)
	}
	if err := os.Remove(unitServicePath); err != nil {
		t.Errorf("Can not delete test path. %v", err)
	}
	err := exec.Command("rm", "-fr", "/test").Run()
	if err != nil {
		t.Errorf("Can not delete test path. %v", err)
	}
	exec.Command("systemctl", "stop", "test").Run()

}

func SkipTesAllDBusTest(t *testing.T) {
	if user, _ := user.Current(); !strings.EqualFold(uid, user.Uid) {
		t.Skipf("User must be root. Execute test with root privileges.")
	}
}

func AddTestServiceDBusTest(t *testing.T) {
	if err := ioutil.WriteFile(unitServicePath, testFile, 0644); err != nil {
		t.Skipf("Can not create test file. %v", err.Error())
	}
	exec.Command("rm", "-fr", "/lib/systemd/system/test.service").Run()
	if err := os.Link(unitServicePath, unitServiceLink); err != nil {
		os.Remove(unitServicePath)
		t.Skipf("Can not create test link. %v", err.Error())
	}
	err := exec.Command("mkdir", "-p", pathTest).Run()
	if err != nil {
		t.Skipf("Can not create directory for test. %v", err.Error())
	}
	exec.Command("systemctl", "stop", "test").Run()
}

func SetupDockerTest(t *testing.T) {
	exec.Command("docker", "rm", "-f", "testValkiria").Run()
	err := exec.Command("docker", "run", "-dit", "--env", "MESOS_TASK_ID=mesos_1232134", "--name", DockerContainerName, DockerImage, "/bin/bash").Run()
	if err != nil {
		t.Skipf("Can not create test container. Is it docker daemon running?? %v", err)
	}
}

func TearDownDockerTest(t *testing.T) {
	err := exec.Command("docker", "rm", "-f", "testValkiria").Run()
	if err != nil {
		t.Errorf("Can not remove test container. %v", err)
	}
}
