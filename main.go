package main

import (
	"fmt"
	"os"

	"github.com/perryrh0dan/dev/cmd"
	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container/docker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	engine, err := docker.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rootCmd := cmd.NewRootCmd(cfg, engine)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
