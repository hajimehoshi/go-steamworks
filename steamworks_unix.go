// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

//go:build !windows

package steamworks

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ebitengine/purego"
)

func loadLib() (uintptr, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return 0, err
	}

	ext := ".so"
	if runtime.GOOS == "darwin" {
		ext = ".dylib"
	}
	path := filepath.Join(dir, "libsteam_api"+ext)
	if err := os.WriteFile(path, libSteamAPI, 0644); err != nil {
		return 0, err
	}

	lib, err := purego.Dlopen(path, purego.RTLD_LAZY|purego.RTLD_LOCAL)
	if err != nil {
		return 0, fmt.Errorf("steamworks: dlopen failed: %w", err)
	}

	return lib, nil
}
