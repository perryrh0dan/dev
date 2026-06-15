package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/perryrh0dan/dev/internal/container"
)

// Engine is the Docker implementation of ContainerEngine.
// It shells out to the `docker` CLI binary.
type Engine struct{}

// New returns a new Docker Engine. Returns an error if docker is not accessible.
func New() (*Engine, error) {
	if err := exec.Command("docker", "info").Run(); err != nil {
		return nil, fmt.Errorf("docker is not running or not accessible: %w", err)
	}
	return &Engine{}, nil
}

func (e *Engine) CreateVolume(name string) error {
	return run("docker", "volume", "create", "--label=dev=yes", name)
}

func (e *Engine) ListVolumes() ([]string, error) {
	out, err := output("docker", "volume", "ls",
		"--filter", "label=dev=yes",
		"--format", "{{.Name}}")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func (e *Engine) RemoveVolumes(names ...string) error {
	args := append([]string{"volume", "rm"}, names...)
	return run("docker", args...)
}

func (e *Engine) RunInteractive(opts container.RunOptions) error {
	args := []string{"run"}
	if opts.Remove {
		args = append(args, "--rm")
	}
	args = append(args, "-it")
	if opts.Privileged {
		args = append(args, "--privileged")
	}
	if opts.Memory != "" {
		args = append(args, "--memory", opts.Memory)
	}
	if opts.Name != "" {
		args = append(args, "--name", opts.Name)
	}
	for _, m := range opts.Mounts {
		args = append(args, "--mount", buildMountFlag(m))
	}
	for _, p := range opts.Ports {
		args = append(args, "-p", fmt.Sprintf("%s:%s", p, p))
	}
	for _, env := range opts.Env {
		args = append(args, "-e", env)
	}
	for k, v := range opts.Labels {
		args = append(args, "--label", k+"="+v)
	}
	args = append(args, opts.Image)
	if len(opts.Command) > 0 {
		args = append(args, opts.Command...)
	}
	return runInteractive("docker", args...)
}

func (e *Engine) Exec(containerID string, cmd []string) error {
	args := append([]string{"exec", "-it", containerID}, cmd...)
	return runInteractive("docker", args...)
}

func (e *Engine) StopContainers(name string) error {
	out, err := output("docker", "ps", "-q", "--filter", "ancestor="+name)
	if err != nil {
		return err
	}
	ids := splitLines(out)
	if len(ids) == 0 {
		return nil
	}
	args := append([]string{"stop"}, ids...)
	return run("docker", args...)
}

func (e *Engine) ListRunningDevContainerNames() ([]string, error) {
	out, err := output("docker", "ps",
		"--filter", "label=dev=yes",
		"--format", "{{.Names}}")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func (e *Engine) PullImage(image string) error {
	return run("docker", "pull", image)
}

func (e *Engine) ListImageTags(image string) ([]string, error) {
	out, err := output("docker", "images", image, "--format", "{{.Tag}}")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

// buildMountFlag converts a Mount to a --mount flag value string.
func buildMountFlag(m container.Mount) string {
	return fmt.Sprintf("type=%s,source=%s,target=%s", m.Type, m.Source, m.Target)
}

// run executes a command, inheriting stdout/stderr.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runInteractive executes a command with stdin/stdout/stderr all attached (for TTY use).
func runInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// output executes a command and returns its stdout as a string.
func output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	return strings.TrimSpace(string(out)), err
}

// splitLines splits a newline-separated string into non-empty lines.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}
