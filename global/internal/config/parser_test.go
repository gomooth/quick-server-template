package config

import (
	"strings"
	"testing"
)

func TestValidate_DuplicateConnectName(t *testing.T) {
	cfg := &ProjectConfig{}
	cfg.Data.Persistent.Enabled = true
	cfg.Data.Persistent.Connects = []dBConnectConfig{
		{Name: "platform", Driver: "sqlite", Dsn: "a.db"},
		{Name: "platform", Driver: "sqlite", Dsn: "b.db"},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected error for duplicate connect name, got nil")
	}
	if !strings.Contains(err.Error(), "duplicated") {
		t.Fatalf("error should mention duplicated, got: %v", err)
	}
}

func TestValidate_EmptyConnectName(t *testing.T) {
	cfg := &ProjectConfig{}
	cfg.Data.Persistent.Enabled = true
	cfg.Data.Persistent.Connects = []dBConnectConfig{
		{Name: "", Driver: "sqlite", Dsn: "a.db"},
	}

	err := cfg.validate()
	if err == nil {
		t.Fatal("expected error for empty connect name, got nil")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Fatalf("error should mention empty, got: %v", err)
	}
}

func TestValidate_UniqueConnectNames(t *testing.T) {
	cfg := &ProjectConfig{}
	cfg.Data.Persistent.Enabled = true
	cfg.Data.Persistent.Connects = []dBConnectConfig{
		{Name: "platform", Driver: "sqlite", Dsn: "a.db"},
		{Name: "analytics", Driver: "mysql", Dsn: "root@tcp(127.0.0.1)/analytics"},
	}

	if err := cfg.validate(); err != nil {
		t.Fatalf("expected nil for unique names, got: %v", err)
	}
}

func TestValidate_DisabledPersistentSkipsCheck(t *testing.T) {
	cfg := &ProjectConfig{}
	cfg.Data.Persistent.Enabled = false
	// 即使有重复名，disabled 时也不校验
	cfg.Data.Persistent.Connects = []dBConnectConfig{
		{Name: "platform", Driver: "sqlite", Dsn: "a.db"},
		{Name: "platform", Driver: "sqlite", Dsn: "b.db"},
	}

	if err := cfg.validate(); err != nil {
		t.Fatalf("expected nil when persistent disabled, got: %v", err)
	}
}
