package config

import (
	"gopkg.in/ini.v1"
	"testing"
)

func TestIni(t *testing.T) {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg.Section("kafka").Key("topic"))
	t.Log(cfg.Section("kafka").Key("address"))
	t.Log(cfg.Section("tailLog").Key("path"))
}

func TestMapTo(t *testing.T) {
	cfg := new(AppConf)
	if err := ini.MapTo(&cfg, "./config.ini"); err != nil {
		t.Fatal(err)
	}
}
