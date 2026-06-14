# Design: `dev` CLI — Go replacement for PowerShell dev environment scripts

**Date:** 2026-06-14  
**Author:** Thomas Pöhlmann  
**Status:** Approved

---

## Overview

A cross-platform Go CLI named `dev` that replaces the PowerShell Docker-based dev environment scripts. Built with [cobra](https://github.com/spf13/cobra) for subcommand routing. A `ContainerEngine` interface abstracts all container and volume operations, with a Docker implementation (shelling out to the `docker` binary) as the first backend — making it straightforward to add Podman or other engines later.

Config files remain at `~/.environments.json` and `~/.identities.json` for zero migration friction. A new `~/.config/dev/config.json` holds all host-side settings (directories, image name, remote hosts).

**Target platforms:** Windows, Linux, macOS.

---

## Command Structure

```
dev [vol|dir]                          # root: launch container, or SSH to remotes

dev container create <name>
dev container list
dev container remove <name>
dev container stop
dev container attach <id>
dev container backup <name>
dev container restore <name> <path>
dev container copy <name> <path>

dev image pull
dev image tags

dev identity add --email --name --keyid
dev identity activate --email
```

---

## Project Structure

```
dev/
├── main.go
├── go.mod / go.sum
├── cmd/
│   ├── root.go               # dev [vol|dir|no-args]
│   ├── container/
│   │   ├── container.go      # parent command
│   │   ├── create.go
│   │   ├── list.go
│   │   ├── remove.go
│   │   ├── stop.go
│   │   ├── attach.go
│   │   ├── backup.go
│   │   ├── restore.go
│   │   └── copy.go
│   ├── image/
│   │   ├── image.go          # parent command
│   │   ├── pull.go
│   │   └── tags.go
│   └── identity/
│       ├── identity.go       # parent command
│       ├── add.go
│       └── activate.go
├── internal/
│   ├── config/
│   │   ├── config.go         # read/write ~/.config/dev/config.json
│   │   ├── environments.go   # read/write ~/.environments.json
│   │   └── identities.go     # read/write ~/.identities.json
│   ├── container/
│   │   ├── engine.go         # ContainerEngine interface + shared types
│   │   └── docker/
│   │       └── docker.go     # Docker implementation
│   └── ssh/
│       └── ssh.go            # remote SSH connection logic
└── docs/superpowers/specs/
    └── 2026-06-14-dev-cli-design.md
```

---

## Container Engine Interface

Defined in `internal/container/engine.go`. All commands receive a `ContainerEngine` via dependency injection from `main.go`. Adding Podman means adding `internal/container/podman/podman.go` implementing this interface — no changes to command code.

```go
type ContainerEngine interface {
    // Volume management
    CreateVolume(name string) error
    ListVolumes() ([]string, error)
    RemoveVolumes(names ...string) error

    // Container lifecycle
    RunInteractive(opts RunOptions) error
    Exec(containerID string, cmd []string) error
    StopContainersByImage(image string) error
    ListContainersByImage(image string) ([]Container, error)

    // Image management
    PullImage(image string) error
    ListImageTags(image string) ([]string, error)
}

type RunOptions struct {
    Image      string
    Name       string
    Mounts     []Mount
    Ports      []string
    Env        []string
    Privileged bool
    Memory     string
    Remove     bool
}

type Mount struct {
    Type   string // "bind" or "volume"
    Source string
    Target string
}

type Container struct {
    ID    string
    Image string
    Name  string
}
```

---

## Config Layer

### `~/.config/dev/config.json` — Host settings

Managed by `internal/config/config.go`. Config file values are primary; environment variables with the same names override them at runtime.

```json
{
  "image": "registry.tpoe.dev/dev",
  "ssh_directory": "/home/user/.ssh",
  "gpg_directory": "/home/user/.gnupg",
  "shared_directory": "/home/user/shared",
  "kube_directory": "/home/user/.kube",
  "timewarrior_directory": "/home/user/.local/share/timewarrior",
  "opencode_directory": "/home/user/.local/share/opencode",
  "env_file": "/home/user/.local/share/.env",
  "dict_file": "/home/user/dict.txt",
  "ngrok_directory": "/home/user/.config/ngrok",
  "npm_file": "/home/user/.npmrc",
  "github_copilot_directory": "/home/user/.config/github-copilot",
  "remote_dev_env": ["192.168.0.1", "192.168.0.2"]
}
```

Go struct:

```go
type Config struct {
    Image                  string   `json:"image"`
    SSHDirectory           string   `json:"ssh_directory"`
    GPGDirectory           string   `json:"gpg_directory"`
    SharedDirectory        string   `json:"shared_directory"`
    KubeDirectory          string   `json:"kube_directory"`
    TimewarriorDirectory   string   `json:"timewarrior_directory"`
    OpenCodeDirectory      string   `json:"opencode_directory"`
    EnvFile                string   `json:"env_file"`
    DictFile               string   `json:"dict_file"`
    NgrokDirectory         string   `json:"ngrok_directory"`
    NpmFile                string   `json:"npm_file"`
    GitHubCopilotDirectory string   `json:"github_copilot_directory"`
    RemoteDevEnv           []string `json:"remote_dev_env"`
}
```

**Override precedence:** config file value → overridden by matching env var.

| Env Var | Config field overridden |
|---|---|
| `DOCKER_DEV_ENV` | `image` |
| `SSH_DIRECTORY` | `ssh_directory` |
| `GPG_DIRECTORY` | `gpg_directory` |
| `SHARED_DIRECTORY` | `shared_directory` |
| `KUBE_DIRECTORY` | `kube_directory` |
| `TIMEWARRIOR_DIRECTORY` | `timewarrior_directory` |
| `OPENCODE_DIRECTORY` | `opencode_directory` |
| `ENV_FILE` | `env_file` |
| `DICT_FILE` | `dict_file` |
| `NGROK_DIRECTORY` | `ngrok_directory` |
| `NPM_FILE` | `npm_file` |
| `GITHUB_COPILOT_DIRECTORY` | `github_copilot_directory` |
| `REMOTE_DEV_ENV` | `remote_dev_env` (comma-separated → split to slice) |

### `~/.environments.json` — Per-environment runtime state

```go
type Environment struct {
    Name string   `json:"name"`
    Port []string `json:"port"`
    Tag  string   `json:"tag"`
}
```

### `~/.identities.json` — Git identities

```go
type Identity struct {
    Email  string `json:"email"`
    Name   string `json:"name"`
    KeyID  string `json:"keyid"`
    Active bool   `json:"active"`
}
```

All config files use atomic writes (write to temp file → rename) to prevent corruption. Missing files are treated as empty and created on first write.

---

## Command Behaviour

### `dev [vol|dir]` — root command

- **No args:** iterate `config.RemoteDevEnv` hosts, SSH into each sequentially.
- **Volume name arg:** mount as named Docker volume at `/root/workspace`; create companion volumes `<name>-history` and `<name>-resurrect` if absent; run container interactively.
- **Directory path arg:** bind-mount the path at `/root/workspace`; run container interactively.
- **Flags:** `--port <port>` (repeatable), `--tag <tag>` — both persisted to `~/.environments.json`.
- Injects `GIT_USER`, `GIT_EMAIL`, `GIT_SIGNINGKEY` from the active identity (if set).
- Mounts each configured directory/file from config; unset fields are silently skipped.
- Container flags: `--privileged --rm -it --memory 24gb`.

### `dev container create <name>`

Creates a Docker volume: `docker volume create --label=dev=yes <name>`.

### `dev container list`

Lists Docker volumes with the `dev=yes` label.

### `dev container remove <name>`

Removes `<name>`, `<name>-history`, and `<name>-resurrect` volumes.

### `dev container stop`

Stops all running containers whose image matches `config.Image`.

### `dev container attach <id>`

Runs `docker exec -it <id> /bin/zsh`.

### `dev container backup <name>`

Spins up a temporary Ubuntu container with the workspace, history, and zoxide volumes mounted plus a local `./backup-<name>/` bind mount. Creates:
- `backup.tar` — workspace, excluding `node_modules`, `.angular`, `.nx`, `dist`, `.pnpm-store`, `.next`
- `backup-history.tar` — shell history volume
- `backup-zoxide.tar` — zoxide database volume

### `dev container restore <name> <path>`

1. Calls create volume for `<name>`.
2. Spins up a temporary Ubuntu container with the same volume layout.
3. Extracts all three tar archives from `<path>` into the respective volumes.

### `dev container copy <name> <path>`

Reads from **stdin**, pipes into `docker exec -i <name> sh -c "cat > <path>"`. Cross-platform; composable with pipes (`pbpaste | dev container copy myenv /root/file`).

### `dev image pull`

Reads `~/.environments.json`, collects all unique non-empty `tag` values, then pulls `config.Image` (latest) and `config.Image:<tag>` for each tag.

### `dev image tags`

Lists local tags for `config.Image`.

### `dev identity add --email --name --keyid`

Upserts an identity into `~/.identities.json` and marks it active (all others become inactive).

### `dev identity activate --email`

Marks the specified identity as active.

---

## Error Handling

- All commands verify Docker is accessible before proceeding.
- Errors from `docker` subprocesses surface the original stderr output.
- Missing config files are treated as empty; created on first write.
- Unknown volume names or container IDs produce a clear error message.

---

## Out of Scope

- Shell utility commands (`watch`, `find_port`, `reload`, `remote_branches`, `refresh_nat`, `update`) — remain in the shell profile.
- Podman or other container engine backends — interface designed for it, only Docker implemented.
- A `dev config` subcommand — users edit `~/.config/dev/config.json` directly for now.
