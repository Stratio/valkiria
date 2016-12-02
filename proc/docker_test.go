package proc

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"testing"
)

const (
	patternTask          = "^/testValkiria"
	fakePatternMesosTask = "^fake"
)

var (
	testKillDockers = func(t *testing.T) {
		r, e := ReadAllDockers(patternTask, FunctionToAddDockerContainerMesosCluster)
		if e != nil {
			t.Errorf("proc.testKillDockers - ReadAllDockers ERROR: %v", e)
		}
		if len(r) != 1 {
			t.Fatalf("proc.testKillDockers - ReadAllDockers ERROR: Conatainer test not found")
		}
		eKill := r[0].Kill()
		if eKill != nil {
			t.Errorf("proc.testKillDockers - ERROR: %v", eKill)
		}
		eKillFake := r[0].Kill()
		if eKillFake == nil {
			t.Errorf("proc.testKillDockers - ERROR: %v", eKillFake)
		}
	}
	testReadAllDockers = func(t *testing.T) {
		r, e := ReadAllDockers(patternTask, FunctionToAddDockerContainerMesosCluster)
		if e != nil {
			t.Errorf("proc.testReadAllDockers - ERROR: %v", e)
		}
		if len(r) != 1 {
			t.Errorf("proc.testReadAllDockers - ERROR: Conatainer test not found")
		}
		rFake, eFake := ReadAllDockers(fakePatternMesosTask, FunctionToAddDockerContainerMesosCluster)
		if eFake != nil {
			t.Errorf("proc.testReadAllDockers - ERROR: %v", e)
		}
		if len(rFake) != 0 {
			t.Errorf("proc.testReadAllDockers - ERROR: Not pattern error and should be.")
		}
		rFakeFunc, eFakeFunc := ReadAllDockers(patternTask, fakeTestFunc)
		if eFakeFunc != nil {
			t.Errorf("proc.testReadAllDockers - ERROR: %v", eFakeFunc)
		}
		if len(rFakeFunc) != 0 {
			t.Errorf("proc.testReadAllDockers - ERROR: Slice ReadAllDockers should be empty.")
		}
	}
)

func TestDockerLib(t *testing.T) {
	log.SetLevel(test.Level)
	test.SetupDockerTest(t)
	defer test.TearDownDockerTest(t)
	t.Run("testReadAllDockers", testReadAllDockers)
	t.Run("testKill", testKillDockers)
}

var fakeTestFunc = func(container types.Container) (*docker, error) {
	//TODO: other new client is mandatory. bug in docker api
	return nil, errors.New("this is a test error, good error...")
}
