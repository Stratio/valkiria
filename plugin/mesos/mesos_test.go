package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/test"
	"testing"
)

const(
	fakeKillExecutor = "killExecutor=a"
)

var (
	testKillFunctionFake = func(t *testing.T) {
		var pluginInstance = NewMesosConfig()
		var functionRead = pluginInstance.FindAndKill()
		_, err := functionRead(nil, "fake", fakeKillExecutor)
		if err == nil {
			t.Errorf("plugin_mesos.testReadAllDaemons - ERROR: %v", err)
		}
	}
)

func TestMesosPluginLib(t *testing.T) {
	log.SetLevel(test.Level)
	t.Run("testKillFunctionFake", testKillFunctionFake)
}
