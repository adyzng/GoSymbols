package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig("..\\config.ini")
	if err != nil {
		t.Error(err)
	}
}
