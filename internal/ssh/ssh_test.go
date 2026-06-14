package ssh

import (
	"testing"
)

func TestBuildSSHArgs(t *testing.T) {
	got := buildSSHArgs("192.168.0.1")
	want := []string{"192.168.0.1"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildSSHArgs_WithUser(t *testing.T) {
	got := buildSSHArgs("user@192.168.0.1")
	want := []string{"user@192.168.0.1"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestConnectAll_NoHosts_ReturnsError(t *testing.T) {
	err := ConnectAll([]string{})
	if err == nil {
		t.Error("expected error for empty hosts, got nil")
	}
}

