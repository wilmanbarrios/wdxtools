package main

import (
	"os"
	"path/filepath"

	cmd "github.com/wilmanbarrios/wdxtools/internal/cmd"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	cmd.SetVersion(version)

	// Busybox pattern: if invoked as "diffh" or "numcrn" (via symlink),
	// behave as "wdxtools diffh ..." or "wdxtools numcrn ...".
	base := filepath.Base(os.Args[0])
	switch base {
	case "diffh", "numcrn":
		os.Args = append([]string{os.Args[0], base}, os.Args[1:]...)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
