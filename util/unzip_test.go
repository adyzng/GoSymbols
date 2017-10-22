package util

import "testing"

func TestUnzip(t *testing.T) {
	zip := "wpt.zip"
	if err := Unzip(zip, "test"); err != nil {
		t.Error(err)
	}
}
