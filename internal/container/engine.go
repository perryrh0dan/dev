package container

// ContainerEngine abstracts container runtime operations.
// Implementations: docker/docker.go (Docker via CLI), future: podman/podman.go
type ContainerEngine interface {
	// Volume management
	CreateVolume(name string) error
	ListVolumes() ([]string, error)
	RemoveVolumes(names ...string) error

	// Container lifecycle
	RunInteractive(opts RunOptions) error
	Exec(containerID string, cmd []string) error
	StopContainers(name string) error
	ListRunningDevContainerNames() ([]string, error)

	// Image management
	PullImage(image string) error
	ListImageTags(image string) ([]string, error)
}

// RunOptions configures a container run invocation.
type RunOptions struct {
	Image      string
	Name       string
	Mounts     []Mount
	Ports      []string
	Env        []string
	Labels     map[string]string
	Privileged bool
	Memory     string
	Remove     bool
	Command    []string // args passed after the image name (overrides default entrypoint args)
}

// Mount describes a bind mount or named volume mount.
type Mount struct {
	Type   string // "bind" or "volume"
	Source string
	Target string
}
