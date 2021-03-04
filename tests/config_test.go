package main

import (
	"github.com/deng00/go-base/config"
	"testing"
)

func TestConfigLocal(t *testing.T) {
	configIns, _ := config.GetConfigFromLocal()
	if configIns == nil || configIns.GetString("log.level") != "debug" {
		t.Errorf("get config from local failed")
	}
}

func TestConfigContent(t *testing.T) {
	content := []byte("level: debug")
	configIns := config.GetConfigFromContent(content)
	if configIns == nil || configIns.GetString("level") != "debug" {
		t.Errorf("get config from content failed")
	}
}
