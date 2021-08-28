// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

const version = "151"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "go-steamworks")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := downloadSDKZip(dir); err != nil {
		return err
	}

	return nil
}

func downloadSDKZip(dir string) error {
	f, err := os.Open(fmt.Sprintf("steamworks_sdk_%s.zip", version))
	if err != nil {
		if os.IsNotExist(err) {
			const sdkURL = "https://partner.steamgames.com/downloads/steamworks_sdk_" + version + ".zip"
			return fmt.Errorf("steamworks_sdk.zip must exist; download it from %s with your Steamworks account", sdkURL)
		}
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	r, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return err
	}

	libs := map[string]string {
		"sdk/redistributable_bin/linux32/libsteam_api.so": "libsteam_api.so",
		"sdk/redistributable_bin/linux64/libsteam_api.so": "libsteam_api64.so",
		"sdk/redistributable_bin/osx/libsteam_api.dylib":  "libsteam_api.dylib",
		"sdk/redistributable_bin/steam_api.dll":           "steam_api.dll",
		"sdk/redistributable_bin/win64/steam_api64.dll":   "steam_api64.dll",
	}
	for _, f := range r.File {
		name, ok := libs[f.Name]
		if !ok {
			continue
		}

		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		out, err := os.Create(name)
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, r); err != nil {
			return err
		}
	}

	return nil
}
