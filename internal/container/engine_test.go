package container_test

import (
	"testing"

	"github.com/perryrh0dan/dev/internal/container"
)

// Compile-time check: any type claiming to implement ContainerEngine must satisfy the interface.
// This test just ensures the interface is importable and the types are defined.
func TestTypesExist(t *testing.T) {
	var _ container.ContainerEngine = (*mockEngine)(nil)
}

type mockEngine struct{}

func (m *mockEngine) CreateVolume(name string) error                        { return nil }
func (m *mockEngine) ListVolumes() ([]string, error)                        { return nil, nil }
func (m *mockEngine) RemoveVolumes(names ...string) error                   { return nil }
func (m *mockEngine) RunInteractive(opts container.RunOptions) error        { return nil }
func (m *mockEngine) Exec(containerID string, cmd []string) error           { return nil }
func (m *mockEngine) StopContainers(name string) error              { return nil }
func (m *mockEngine) ListRunningDevContainerNames() ([]string, error) {
	return nil, nil
}
func (m *mockEngine) PullImage(image string) error                          { return nil }
func (m *mockEngine) ListImageTags(image string) ([]string, error)          { return nil, nil }
