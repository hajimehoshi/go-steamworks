// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

package steamworks

import "syscall"

func loadLib() (uintptr, error) {
	dllName := "steam_api64.dll"
	handle, err := syscall.LoadLibrary(dllName)
	return uintptr(handle), err
}
