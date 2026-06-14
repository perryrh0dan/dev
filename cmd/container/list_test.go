package container

import (
	"bytes"
	"strings"
	"testing"

	"github.com/perryrh0dan/dev/internal/container"
)

type fakeEngine struct {
	volumes      []string
	runningNames []string
	volErr       error
	runErr       error
}

func (f *fakeEngine) CreateVolume(name string) error                 { return nil }
func (f *fakeEngine) ListVolumes() ([]string, error)                 { return f.volumes, f.volErr }
func (f *fakeEngine) RemoveVolumes(names ...string) error            { return nil }
func (f *fakeEngine) RunInteractive(opts container.RunOptions) error { return nil }
func (f *fakeEngine) Exec(id string, cmd []string) error             { return nil }
func (f *fakeEngine) StopContainers(name string) error       { return nil }
func (f *fakeEngine) ListRunningDevContainerNames() ([]string, error) {
	return f.runningNames, f.runErr
}
func (f *fakeEngine) PullImage(image string) error               { return nil }
func (f *fakeEngine) ListImageTags(image string) ([]string, error) { return nil, nil }

func TestListCmd_ShowsStatusColumn(t *testing.T) {
	eng := &fakeEngine{
		volumes:      []string{"my-env", "other-env"},
		runningNames: []string{"my-env"},
	}

	cmd := newListCmd(eng)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "my-env") || !strings.Contains(out, "RUNNING") {
		t.Errorf("expected my-env RUNNING in output, got:\n%s", out)
	}
	if !strings.Contains(out, "other-env") || !strings.Contains(out, "STOPPED") {
		t.Errorf("expected other-env STOPPED in output, got:\n%s", out)
	}
}

func TestListCmd_RunningFlag_FiltersToRunning(t *testing.T) {
	eng := &fakeEngine{
		volumes:      []string{"my-env", "other-env"},
		runningNames: []string{"my-env"},
	}

	cmd := newListCmd(eng)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--running"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "my-env") {
		t.Errorf("expected my-env in output, got:\n%s", out)
	}
	if strings.Contains(out, "other-env") {
		t.Errorf("expected other-env to be filtered out, got:\n%s", out)
	}
}
