package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GetBinPath(name string) (string, error) {
	// 1. search $PATH (covers /opt/bushuray, /usr/bin, etc. if added to PATH)
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	// 2. same directory as the running executable (sibling binary)
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, nil
		}
	}

	// 3. hardcoded fallbacks
	paths := []string{
		"./" + name,
		filepath.Join(".", "bin", name),
		filepath.Join("/usr/bin", name),
		filepath.Join("/usr/local/bin", name),
	}

	for _, p := range paths {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return filepath.Abs(p)
		}
	}

	return "", fmt.Errorf("binary %q not found", name)
}
