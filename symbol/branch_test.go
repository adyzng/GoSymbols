package symbol

import (
	"fmt"
	"testing"
	"time"
)

func TestParseBuilds(t *testing.T) {
	lastBuild := ""
	builder := NewBranch("UDP_6_5_U2", "UDPv6.5U2")

	total, err := builder.ParseBuilds(func(b *Build) error {
		fmt.Printf("%s: %+v \n", b.ID, b)
		lastBuild = b.ID
		return nil
	})
	fmt.Printf("Branch %s has %d builds.\n", builder.Name(), total)
	if err != nil {
		t.Error(err)
	}

	idx := 0
	total, err = builder.ParseSymbols(lastBuild, func(sym *Symbol) error {
		fmt.Printf(" %d: %+v\n", idx, sym)
		idx++
		return nil
	})
	fmt.Printf("Branch %s build %s has %d symbols.\n", builder.Name(), lastBuild, total)
	if err != nil {
		t.Error(err)
	}
}

func TestAddBuild(t *testing.T) {
	builder := NewBranch("UDP_6_5_U2", "UDPv6.5U2")
	if err := builder.AddBuild(""); err != nil {
		time.Sleep(time.Second)
		t.Fatal(err)
	}

	idx, lastBuild := 0, builder.GetLatestID()
	total, err := builder.ParseSymbols(lastBuild, func(sym *Symbol) error {
		fmt.Printf(" %d: %+v\n", idx, sym)
		idx++
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Branch %s build %s has %d symbols.\n", builder.Name(), lastBuild, total)
}
