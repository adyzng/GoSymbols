package symbol

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/adyzng/goPDB/config"
)

func TestParseBuilds(t *testing.T) {
	builder := NewBranch("UDP_6_5_U2", "UDPv6.5U2")
	if err := builder.Init(); err != nil {
		t.Fatal(err)
	}

	lastBuild := ""
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
	builder := NewBranch("UDP_6_5_U2", "UDPV6.5U2")
	if err := builder.Init(); err != nil {
		t.Fatal(err)
	}

	if err := builder.Add(""); err != nil {
		time.Sleep(time.Second)
		t.Fatal(err)
	}

	idx, lastBuild := 0, builder.GetLatestID()
	total, err := builder.ParseSymbols(lastBuild, func(sym *Symbol) error {
		fmt.Printf(" %d: %+v\n", idx, sym)
		idx++
		return nil
	})
	fmt.Printf("Branch %s build %s has %d symbols.\n", builder.Name(), lastBuild, total)

	if err != nil {
		t.Error(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	b1 := NewBranch("UDP_6_5_U1", "UDPv6.5U1")
	if err := b1.Init(); err != nil {
		fmt.Printf("Init error: %v.\n", err)
	}
	fmt.Printf("%+v\n", b1)

	/*
		b2 := NewBranch("UDP_6_5_U1", "UDPV6.5U1")
		//json.NewEncoder(os.Stdout).Encode(&b1)
		arr := []Builder{b1, b2}

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(&arr)

		var b []Builder
		dec := gob.NewDecoder(&buffer)
		err := dec.Decode(&b)
		if err != nil {
			fmt.Println("Decode error:", err)
		} else {
			fmt.Printf("Decode: %+v.\n", b[0])
		}
	*/
}

func TestScanBuilders(t *testing.T) {
	root := config.Destination
	fmt.Println(root)

	if fi, err := ioutil.ReadDir(root); err == nil {
		for _, f := range fi {
			if f.IsDir() {
				fmt.Println(f.Name())
				b := NewBranch(f.Name(), f.Name())
				switch b.Init() {
				case nil:
					fmt.Printf("Exist branch %+v.\n", b)
				case ErrBranchNotInit:
					fmt.Printf("New branch %+v.\n", b)
					if err = b.Persist(); err != nil {
						t.Error(err)
					}
				default:
					fmt.Printf("Invalid branch %+v.\n", b)
				}
			}
		}
	}

	time.Sleep(time.Second)
}
