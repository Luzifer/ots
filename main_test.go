package main

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func mustCloseListener(t *testing.T, l net.Listener, msg string) {
	t.Helper()

	if err := l.Close(); err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func TestGetListener(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ots-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("failed to remove temp dir: %v", err)
		}
	}()

	t.Run("tcp listener", func(t *testing.T) {
		l, err := getListener("127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to create TCP listener: %v", err)
		}
		if l.Addr().Network() != "tcp" {
			t.Errorf("expected tcp network, got %s", l.Addr().Network())
		}
		mustCloseListener(t, l, "failed to close TCP listener")
	})

	t.Run("unix listener", func(t *testing.T) {
		sockPath := filepath.Join(tmpDir, "test.sock")
		addr := "unix:" + sockPath

		l, err := getListener(addr)
		if err != nil {
			t.Fatalf("failed to create unix listener: %v", err)
		}
		if l.Addr().Network() != "unix" {
			t.Errorf("expected unix network, got %s", l.Addr().Network())
		}
		if l.Addr().String() != sockPath {
			t.Errorf("expected socket path %s, got %s", sockPath, l.Addr().String())
		}
		mustCloseListener(t, l, "failed to close unix listener")
	})

	t.Run("existing socket deletion", func(t *testing.T) {
		sockPath := filepath.Join(tmpDir, "test.sock")
		addr := "unix:" + sockPath

		lc := net.ListenConfig{}
		l, err := lc.Listen(context.Background(), "unix", sockPath)
		if err != nil {
			t.Fatalf("failed to create manual unix listener: %v", err)
		}
		mustCloseListener(t, l, "failed to close manual unix listener")

		l, err = getListener(addr)
		if err != nil {
			t.Fatalf("failed to create unix listener with existing socket: %v", err)
		}
		mustCloseListener(t, l, "failed to close unix listener with existing socket")
	})

	t.Run("nested directory creation", func(t *testing.T) {
		nestedSockPath := filepath.Join(tmpDir, "nested/dir/test.sock")
		addr := "unix:" + nestedSockPath

		l, err := getListener(addr)
		if err != nil {
			t.Fatalf("failed to create unix listener in nested dir: %v", err)
		}
		mustCloseListener(t, l, "failed to close nested unix listener")
	})
}
