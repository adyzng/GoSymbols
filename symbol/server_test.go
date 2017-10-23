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
	ss.WalkBuilders(func(bd Builder) error {
		fmt.Printf("%+v\n", bd)
		return nil
	})
}

func TestAddBranch(t *testing.T) {
	bn, sn := "UDP_6_5_U2", "UDPv6.5U2"
	builder := GetServer().AddBuilder(bn, sn)

	if builder != nil {
		if err := builder.SetSubpath("", ""); err != nil {
			fmt.Printf("Set branch path failed: %v.\n", err)
			t.Error(err)
		}
		fmt.Printf("Add branch: %+v.\n", builder)
	} else {
		t.Fatal("Add branch failed.\n")
	}
}

func TestAddBranchForBrowseOnly(t *testing.T) {
	bn, sn := "UDP_6_5_U1", "UDPv6.5U1"
	builder := GetServer().AddBuilder(bn, sn)

	if builder != nil {
		if err := builder.SetSubpath("", ""); err != nil {
			fmt.Printf("Set branch path failed: %v\n", err)
			t.Error(err)
		}
		if builder.CanBrowse() {
			if err := builder.Persist(); err != nil {
				fmt.Printf("Failed to branch: %+v. %v.\n", builder, err)
				t.Error(err)
			}
		}
		fmt.Printf("Add branch: %+v.\n", builder)
	} else {
		t.Fatal("Add branch failed.\n")
	}
}
