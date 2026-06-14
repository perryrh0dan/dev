package docker

import (
	"testing"

	"github.com/perryrh0dan/dev/internal/container"
)

func TestBuildMountFlag_Bind(t *testing.T) {
	m := container.Mount{Type: "bind", Source: "/host/path", Target: "/container/path"}
	got := buildMountFlag(m)
	want := "type=bind,source=/host/path,target=/container/path"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildMountFlag_Volume(t *testing.T) {
	m := container.Mount{Type: "volume", Source: "myvolume", Target: "/root/workspace"}
	got := buildMountFlag(m)
	want := "type=volume,source=myvolume,target=/root/workspace"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEngine_ImplementsInterface(t *testing.T) {
	var _ container.ContainerEngine = (*Engine)(nil)
}
