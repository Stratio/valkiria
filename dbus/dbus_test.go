package dbus

import (
	"testing"
)

func TestinitDefaultSharedObjectService(t *testing.T) {
	var err = initDefaultSharedObjectService()
	if err != nil {
		t.Errorf("dbus.TestinitDefaultSharedObjectService - ERROR:", err)
	}
	t.Log("dbus.TestinitDefaultSharedObjectService")
}
