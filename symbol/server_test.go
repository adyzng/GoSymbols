package symbol

import (
	"fmt"
	"testing"

	"github.com/adyzng/GoSymbols/config"
)

func TestLoadBranch(t *testing.T) {
	ss := GetServer()
	if err := ss.LoadBranchs(); err != nil {
		t.Error(err)
	}
	total := 0
	ss.WalkBuilders(func(b Builder) error {
		fmt.Printf("Load %d: %+v\n", total, b)
		total++
		return nil
	})
}

func TestAddBranch(t *testing.T) {
	bn, sn := "UDP_6_5_U2", "UDPv6.5U2"
	builder := GetServer().Add(bn, sn)

	if builder != nil {
		fmt.Printf("Add branch: %+v.\n", builder)
	} else {
		t.Fatal("Add branch failed.\n")
	}
}

func TestScanBranchs(t *testing.T) {
	ss := GetServer()
	if err := ss.ScanStore(config.Destination); err != nil {
		t.Error(err)
	}

	total := 0
	ss.WalkBuilders(func(b Builder) error {
		fmt.Printf("%d: %+v\n", total, b)
		total++
		return nil
	})
}
