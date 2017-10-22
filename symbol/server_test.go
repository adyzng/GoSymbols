package symbol

import (
	"fmt"
	"testing"
)

func TestLoadBranch(t *testing.T) {
	ss := GetServer()
	if err := ss.LoadBuilders(); err != nil {
		t.Error(err)
	}
}

func TestAddBranch(t *testing.T) {
	bn, sn := "UDP_6_5_U2", "UDPv6.5U2"
	builder := GetServer().AddBuilder(bn, sn)

	if builder == nil {
		t.Fatal("Add branch failed.")
	} else {
		fmt.Printf("Add branch: %+v.\n", builder)
	}
}
