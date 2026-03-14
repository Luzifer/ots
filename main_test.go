package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestGetListener(t *testing.T) {
	// Test TCP listener
	l, err := getListener("127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create TCP listener: %v", err)
	}
	if l.Addr().Network() != "tcp" {
		t.Errorf("expected tcp network, got %s", l.Addr().Network())
	}
	l.Close()

	// Test Unix listener
	tmpDir, err := os.MkdirTemp("", "ots-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sockPath := filepath.Join(tmpDir, "test.sock")
	addr := "unix:" + sockPath

	l, err = getListener(addr)
	if err != nil {
		t.Fatalf("failed to create unix listener: %v", err)
	}
	if l.Addr().Network() != "unix" {
		t.Errorf("expected unix network, got %s", l.Addr().Network())
	}
	if l.Addr().String() != sockPath {
		t.Errorf("expected socket path %s, got %s", sockPath, l.Addr().String())
	}
	l.Close()

	// Test existing socket deletion
	l, err = net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to create manual unix listener: %v", err)
	}
	l.Close()

	l, err = getListener(addr)
	if err != nil {
		t.Fatalf("failed to create unix listener with existing socket: %v", err)
	}
	l.Close()

	// Test nested directory creation
	nestedSockPath := filepath.Join(tmpDir, "nested/dir/test.sock")
	addr = "unix:" + nestedSockPath
	l, err = getListener(addr)
	if err != nil {
		t.Fatalf("failed to create unix listener in nested dir: %v", err)
	}
	l.Close()
}
