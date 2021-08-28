// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

//go:build !windows
// +build !windows

package steamworks

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -ldl
//
// #include <stdbool.h>
// #include <stdint.h>
// #include <stdlib.h>
// #include <dlfcn.h>
//
// static uintptr_t dlsym_(uintptr_t handle, const char* name) {
//   return (uintptr_t)dlsym((void*)handle, name);
// }
//
// static const char* uintptrToChar(uintptr_t str) {
//   return (const char*)str;
// }
//
// static uintptr_t callFunc_Bool(uintptr_t f) {
//   return ((bool (*)())(f))();
// }
//
// static uintptr_t callFunc_Bool_Int(uintptr_t f, int arg0) {
//   return ((bool (*)(int))(f))(arg0);
// }
//
// static uintptr_t callFunc_Ptr(uintptr_t f) {
//   return (uintptr_t)((void* (*)())(f))();
// }
//
// static uintptr_t callFunc_Ptr_Ptr(uintptr_t f, uintptr_t arg0) {
//   return (uintptr_t)((void* (*)(void*))(f))((void*)arg0);
// }
import "C"

type lib struct {
	lib   C.uintptr_t
	procs map[string]C.uintptr_t
}

func (l *lib) call(name string, args ...uintptr) (C.uintptr_t, error) {
	if l.procs == nil {
		l.procs = map[string]C.uintptr_t{}
	}

	if _, ok := l.procs[name]; !ok {
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		l.procs[name] = C.dlsym_(l.lib, cname)
	}

	f := l.procs[name]
	switch name {
	case flatAPI_RestartAppIfNecessary:
		return C.callFunc_Bool_Int(f, C.int(args[0])), nil
	case flatAPI_Init:
		return C.callFunc_Bool(f), nil
	case flatAPI_SteamApps:
		return C.callFunc_Ptr(f), nil
	case flatAPI_ISteamApps_GetCurrentGameLanguage:
		return C.callFunc_Ptr_Ptr(f, C.uintptr_t(args[0])), nil
	}

	return 0, fmt.Errorf("steamworks: function %s not implemented", name)
}

func loadLib() (C.uintptr_t, error) {
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

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	lib := C.uintptr_t(uintptr(C.dlopen(cpath, C.RTLD_LAZY)))
	if lib == 0 {
		return 0, fmt.Errorf("steamworks: dlopen failed: %s", C.GoString(C.dlerror()))
	}

	return lib, nil
}

var theLib *lib

func init() {
	l, err := loadLib()
	if err != nil {
		panic(err)
	}
	theLib = &lib{
		lib: l,
	}
}

func RestartAppIfNecessary(appID int) bool {
	v, err := theLib.call(flatAPI_RestartAppIfNecessary, uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func Init() bool {
	v, err := theLib.call(flatAPI_Init)
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func SteamApps() ISteamApps {
	v, err := theLib.call(flatAPI_SteamApps)
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps C.uintptr_t

func (s steamApps) GetCurrentGameLanguage() string {
	v, err := theLib.call(flatAPI_ISteamApps_GetCurrentGameLanguage, uintptr(s))
	if err != nil {
		panic(err)
	}
	return C.GoString(C.uintptrToChar(v))
}
