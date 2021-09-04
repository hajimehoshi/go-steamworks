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
// static uintptr_t callFunc_Bool_Ptr(uintptr_t f, uintptr_t arg0) {
//   return ((bool (*)(void*))(f))((void*)arg0);
// }
//
// static uintptr_t callFunc_Bool_Ptr_Ptr(uintptr_t f, uintptr_t arg0, uintptr_t arg1) {
//   return ((bool (*)(void*, void*))(f))((void*)arg0, (void*)arg1);
// }
//
// static uintptr_t callFunc_Bool_Ptr_Ptr_Ptr(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2) {
//   return ((bool (*)(void*, void*, void*))(f))((void*)arg0, (void*)arg1, (void*)arg2);
// }
//
// static uintptr_t callFunc_Bool_Ptr_Ptr_Ptr_Int32(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2, int32_t arg3) {
//   return ((bool (*)(void*, void*, void*, int32_t))(f))((void*)arg0, (void*)arg1, (void*)arg2, arg3);
// }
//
// static uintptr_t callFunc_Bool_Uint32(uintptr_t f, uint32_t arg0) {
//   return ((bool (*)(uint32_t))(f))(arg0);
// }
//
// static uintptr_t callFunc_Int32_Ptr_Ptr_Ptr_Int32(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2, int32_t arg3) {
//   return ((int32_t (*)(void*, void*, void*, int32_t))(f))((void*)arg0, (void*)arg1, (void*)arg2, arg3);
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

type funcType int

const (
	funcType_Bool funcType = iota
	funcType_Bool_Ptr
	funcType_Bool_Ptr_Ptr
	funcType_Bool_Ptr_Ptr_Ptr
	funcType_Bool_Ptr_Ptr_Ptr_Int32
	funcType_Bool_Uint32
	funcType_Int32_Ptr_Ptr_Ptr_Int32
	funcType_Ptr
	funcType_Ptr_Ptr
)

func (l *lib) call(ftype funcType, name string, args ...uintptr) (C.uintptr_t, error) {
	if l.procs == nil {
		l.procs = map[string]C.uintptr_t{}
	}

	if _, ok := l.procs[name]; !ok {
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		l.procs[name] = C.dlsym_(l.lib, cname)
	}

	f := l.procs[name]
	switch ftype {
	case funcType_Bool:
		return C.callFunc_Bool(f), nil
	case funcType_Bool_Ptr:
		return C.callFunc_Bool_Ptr(f, C.uintptr_t(args[0])), nil
	case funcType_Bool_Ptr_Ptr:
		return C.callFunc_Bool_Ptr_Ptr(f, C.uintptr_t(args[0]), C.uintptr_t(args[1])), nil
	case funcType_Bool_Ptr_Ptr_Ptr:
		return C.callFunc_Bool_Ptr_Ptr_Ptr(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2])), nil
	case funcType_Bool_Ptr_Ptr_Ptr_Int32:
		return C.callFunc_Bool_Ptr_Ptr_Ptr_Int32(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.int32_t(args[3])), nil
	case funcType_Bool_Uint32:
		return C.callFunc_Bool_Uint32(f, C.uint32_t(args[0])), nil
	case funcType_Int32_Ptr_Ptr_Ptr_Int32:
		return C.callFunc_Int32_Ptr_Ptr_Ptr_Int32(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.int32_t(args[3])), nil
	case funcType_Ptr:
		return C.callFunc_Ptr(f), nil
	case funcType_Ptr_Ptr:
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

func cBool(x bool) uintptr {
	if x {
		return 1
	}
	return 0
}

func RestartAppIfNecessary(appID int) bool {
	v, err := theLib.call(funcType_Bool_Uint32, flatAPI_RestartAppIfNecessary, uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func Init() bool {
	v, err := theLib.call(funcType_Bool, flatAPI_Init)
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func SteamApps() ISteamApps {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamApps)
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps C.uintptr_t

func (s steamApps) GetCurrentGameLanguage() string {
	v, err := theLib.call(funcType_Ptr_Ptr, flatAPI_ISteamApps_GetCurrentGameLanguage, uintptr(s))
	if err != nil {
		panic(err)
	}
	return C.GoString(C.uintptrToChar(v))
}

func SteamRemoteStorage() ISteamRemoteStorage {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamRemoteStorage)
	if err != nil {
		panic(err)
	}
	return steamRemoteStorage(v)
}

type steamRemoteStorage C.uintptr_t

func (s steamRemoteStorage) FileWrite(file string, data []byte) bool {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	defer runtime.KeepAlive(data)

	v, err := theLib.call(funcType_Bool_Ptr_Ptr_Ptr_Int32, flatAPI_ISteamRemoteStorage_FileWrite, uintptr(s), uintptr(unsafe.Pointer(cfile)), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	defer runtime.KeepAlive(data)

	v, err := theLib.call(funcType_Int32_Ptr_Ptr_Ptr_Int32, flatAPI_ISteamRemoteStorage_FileRead, uintptr(s), uintptr(unsafe.Pointer(cfile)), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	v, err := theLib.call(funcType_Bool_Ptr_Ptr, flatAPI_ISteamRemoteStorage_FileDelete, uintptr(s), uintptr(unsafe.Pointer(cfile)))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func SteamUserStats() ISteamUserStats {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamUserStats)
	if err != nil {
		panic(err)
	}
	return steamUserStats(v)
}

type steamUserStats C.uintptr_t

func (s steamUserStats) GetAchievement(name string) (achieved, success bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	v, err := theLib.call(funcType_Bool_Ptr_Ptr_Ptr, flatAPI_ISteamUserStats_GetAchievement, uintptr(s), uintptr(unsafe.Pointer(cname)), uintptr(unsafe.Pointer(&achieved)))
	if err != nil {
		panic(err)
	}
	success = byte(v) != 0

	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	v, err := theLib.call(funcType_Bool_Ptr_Ptr, flatAPI_ISteamUserStats_SetAchievement, uintptr(s), uintptr(unsafe.Pointer(cname)))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) ClearAchievement(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	v, err := theLib.call(funcType_Bool_Ptr_Ptr, flatAPI_ISteamUserStats_ClearAchievement, uintptr(s), uintptr(unsafe.Pointer(cname)))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) StoreStats() bool {
	v, err := theLib.call(funcType_Bool_Ptr, flatAPI_ISteamUserStats_StoreStats, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}
