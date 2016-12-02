package test

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/docker/docker/client"
	//"github.com/docker/distribution/context"
	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/api/types"
	"io/ioutil"
	"strings"
	"testing"
	"os/user"
	"os"
	"os/exec"
)

const (
	Level           = log.InfoLevel
)

const (
	uid             = "0"
	unitServicePath = "/tmp/test.service"
	unitServiceLink = "/lib/systemd/system/test.service"
)

const (
	DockerImage             = "ubuntu"
	DockerContainerName 	= "testValkiria"

)

var (
	testFile = []byte("[Unit]\nDescription=test\n\n[Service]\nExecStart=/bin/bash -c 'while true; do echo hello; sleep 2; done'\nRestart=always\n")
)

func SetupDBusTest(t *testing.T) {
	SkipTesAllDBusTest(t)
	AddTestServiceDBusTest(t)
}

func TearDownDBusTest(t *testing.T) {
	if err := os.Remove(unitServiceLink); err != nil {
		t.Fatalf("Can not delete test link. %v", err)
	}
	if err := os.Remove(unitServicePath); err != nil {
		t.Fatalf("Can not delete test path. %v", err)
	}
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
	if err := os.Link(unitServicePath, unitServiceLink); err != nil {
		os.Remove(unitServicePath)
		t.Skipf("Can not create test link. %v", err.Error())
	}
}

func SetupDockerTest(t *testing.T) {
	err := exec.Command("docker", "run", "-dit", "--env", "MESOS_TASK_ID=mesos_1232134", "--name", DockerContainerName, DockerImage, "/bin/bash").Run()
	if err != nil {
		t.Skipf("Can not create test container. Is it docker daemon running?? %v", err)
	}
}

func TearDownDockerTest(t *testing.T) {
	err := exec.Command("docker", "rm", "testValkiria").Run()
	if err != nil {
		t.Fatalf("Can not remove test container. %v", err)
	}
}