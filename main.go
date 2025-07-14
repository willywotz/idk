package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
)

var version = "v0"

func main() {
	if err := update(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("version: %s\n", version)
	_, _ = fmt.Scanln()
}

func update() error {
	if exe, _ := os.Executable(); strings.Contains(exe, "go-build") {
		return nil
	}

	fmt.Println("Checking for updates...")

	ctx := context.Background()
	repository := selfupdate.ParseSlug("willywotz/idk")
	release, err := selfupdate.UpdateSelf(ctx, version, repository)
	if err != nil {
		return fmt.Errorf("failed to update self: %w", err)
	}

	if release.GreaterThan(version) {
		fmt.Printf("Updated to version %s, restarting...\n", release.Version())

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		if _, err := os.StartProcess(exe, os.Args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		}); err != nil {
			return fmt.Errorf("failed to restart: %w", err)
		}

		os.Exit(0)
	}

	return nil
}
