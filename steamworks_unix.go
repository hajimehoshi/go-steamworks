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
// static uint8_t callFunc_Bool(uintptr_t f) {
//   return ((bool (*)())(f))();
// }
//
// static uint8_t callFunc_Bool_Ptr(uintptr_t f, uintptr_t arg0) {
//   return ((bool (*)(void*))(f))((void*)arg0);
// }
//
// static uint8_t callFunc_Bool_Ptr_Bool(uintptr_t f, uintptr_t arg0, uint8_t arg1) {
//   return ((bool (*)(void*, bool))(f))((void*)arg0, (bool)arg1);
// }
//
// static uint8_t callFunc_Bool_Ptr_Int32(uintptr_t f, uintptr_t arg0, int32_t arg1) {
//   return ((bool (*)(void*, int32_t))(f))((void*)arg0, arg1);
// }
//
// static uint8_t callFunc_Bool_Ptr_Int32_Int32_Int32_Int32_Int32(uintptr_t f, uintptr_t arg0, int32_t arg1, int32_t arg2, int32_t arg3, int32_t arg4, int32_t arg5) {
//   return ((bool (*)(void*, int32_t, int32_t, int32_t, int32_t, int32_t))(f))((void*)arg0, arg1, arg2, arg3, arg4, arg5);
// }
//
// static uint8_t callFunc_Bool_Ptr_Ptr(uintptr_t f, uintptr_t arg0, uintptr_t arg1) {
//   return ((bool (*)(void*, void*))(f))((void*)arg0, (void*)arg1);
// }
//
// static uint8_t callFunc_Bool_Ptr_Ptr_Ptr(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2) {
//   return ((bool (*)(void*, void*, void*))(f))((void*)arg0, (void*)arg1, (void*)arg2);
// }
//
// static uint8_t callFunc_Bool_Ptr_Ptr_Ptr_Int32(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2, int32_t arg3) {
//   return ((bool (*)(void*, void*, void*, int32_t))(f))((void*)arg0, (void*)arg1, (void*)arg2, arg3);
// }
//
// static uint8_t callFunc_Bool_Int32(uintptr_t f, uint32_t arg0) {
//   return ((bool (*)(uint32_t))(f))(arg0);
// }
//
// static int32_t callFunc_Int32_Ptr(uintptr_t f, uintptr_t arg0) {
//   return ((int32_t (*)(void*))(f))((void*)arg0);
// }
//
// static int32_t callFunc_Int32_Ptr_Int32_Ptr_Int32(uintptr_t f, uintptr_t arg0, int32_t arg1, uintptr_t arg2, int32_t arg3) {
//   return ((int32_t (*)(void*, int32_t, void*, int32_t))(f))((void*)arg0, arg1, (void*)arg2, arg3);
// }
//
// static int32_t callFunc_Int32_Ptr_Int32_Ptr_Ptr_Ptr_Int32(uintptr_t f, uintptr_t arg0, int32_t arg1, uintptr_t arg2, uintptr_t arg3, uintptr_t arg4, int32_t arg5) {
//   return ((int32_t (*)(void*, int32_t, void*, void*, void*, int32_t))(f))((void*)arg0, arg1, (void*)arg2, (void*)arg3, (void*)arg4, arg5);
// }
//
// static int32_t callFunc_Int32_Ptr_Int64(uintptr_t f, uintptr_t arg0, int64_t arg1) {
//   return ((int32_t (*)(void*, int64_t))(f))((void*)arg0, arg1);
// }
//
// static int32_t callFunc_Int32_Ptr_Ptr(uintptr_t f, uintptr_t arg0, uintptr_t arg1) {
//   return ((int32_t (*)(void*, void*))(f))((void*)arg0, (void*)arg1);
// }
//
// static int32_t callFunc_Int32_Ptr_Ptr_Ptr_Int32(uintptr_t f, uintptr_t arg0, uintptr_t arg1, uintptr_t arg2, int32_t arg3) {
//   return ((int32_t (*)(void*, void*, void*, int32_t))(f))((void*)arg0, (void*)arg1, (void*)arg2, arg3);
// }
//
// static int64_t callFunc_Int64_Ptr(uintptr_t f, uintptr_t arg0) {
//   return ((int64_t (*)(void*))(f))((void*)arg0);
// }
//
// static uintptr_t callFunc_Ptr(uintptr_t f) {
//   return (uintptr_t)((void* (*)())(f))();
// }
//
// static uintptr_t callFunc_Ptr_Ptr(uintptr_t f, uintptr_t arg0) {
//   return (uintptr_t)((void* (*)(void*))(f))((void*)arg0);
// }
//
// static void callFunc_Void(uintptr_t f) {
//   ((void (*)())(f))();
// }
//
// static void callFunc_Void_Ptr_Bool(uintptr_t f, uintptr_t arg0, uint8_t arg1) {
//   ((void (*)(void*, bool))(f))((void*)arg0, (bool)arg1);
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
	funcType_Bool_Ptr_Bool
	funcType_Bool_Ptr_Int32
	funcType_Bool_Ptr_Int32_Int32_Int32_Int32_Int32
	funcType_Bool_Ptr_Ptr
	funcType_Bool_Ptr_Ptr_Ptr
	funcType_Bool_Ptr_Ptr_Ptr_Int32
	funcType_Bool_Int32
	funcType_Int32_Int64
	funcType_Int32_Ptr
	funcType_Int32_Ptr_Int32_Ptr_Int32
	funcType_Int32_Ptr_Int32_Ptr_Ptr_Ptr_Int32
	funcType_Int32_Ptr_Int64
	funcType_Int32_Ptr_Ptr
	funcType_Int32_Ptr_Ptr_Ptr_Int32
	funcType_Int64_Ptr
	funcType_Ptr
	funcType_Ptr_Ptr
	funcType_Void
	funcType_Void_Ptr_Bool
)

func (l *lib) call(ftype funcType, name string, args ...uintptr) (C.uint64_t, error) {
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
		return C.uint64_t(C.callFunc_Bool(f)), nil
	case funcType_Bool_Ptr:
		return C.uint64_t(C.callFunc_Bool_Ptr(f, C.uintptr_t(args[0]))), nil
	case funcType_Bool_Ptr_Bool:
		return C.uint64_t(C.callFunc_Bool_Ptr_Bool(f, C.uintptr_t(args[0]), C.uint8_t(args[1]))), nil
	case funcType_Bool_Ptr_Int32:
		return C.uint64_t(C.callFunc_Bool_Ptr_Int32(f, C.uintptr_t(args[0]), C.int32_t(args[1]))), nil
	case funcType_Bool_Ptr_Int32_Int32_Int32_Int32_Int32:
		return C.uint64_t(C.callFunc_Bool_Ptr_Int32_Int32_Int32_Int32_Int32(f, C.uintptr_t(args[0]), C.int32_t(args[1]), C.int32_t(args[2]), C.int32_t(args[3]), C.int32_t(args[4]), C.int32_t(args[5]))), nil
	case funcType_Bool_Ptr_Ptr:
		return C.uint64_t(C.callFunc_Bool_Ptr_Ptr(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]))), nil
	case funcType_Bool_Ptr_Ptr_Ptr:
		return C.uint64_t(C.callFunc_Bool_Ptr_Ptr_Ptr(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]))), nil
	case funcType_Bool_Ptr_Ptr_Ptr_Int32:
		return C.uint64_t(C.callFunc_Bool_Ptr_Ptr_Ptr_Int32(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.int32_t(args[3]))), nil
	case funcType_Bool_Int32:
		return C.uint64_t(C.callFunc_Bool_Int32(f, C.uint32_t(args[0]))), nil
	case funcType_Int32_Ptr:
		return C.uint64_t(C.callFunc_Int32_Ptr(f, C.uintptr_t(args[0]))), nil
	case funcType_Int32_Ptr_Int32_Ptr_Int32:
		return C.uint64_t(C.callFunc_Int32_Ptr_Int32_Ptr_Int32(f, C.uintptr_t(args[0]), C.int32_t(args[1]), C.uintptr_t(args[2]), C.int32_t(args[3]))), nil
	case funcType_Int32_Ptr_Int32_Ptr_Ptr_Ptr_Int32:
		return C.uint64_t(C.callFunc_Int32_Ptr_Int32_Ptr_Ptr_Ptr_Int32(f, C.uintptr_t(args[0]), C.int32_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.int32_t(args[5]))), nil
	case funcType_Int32_Ptr_Int64:
		return C.uint64_t(C.callFunc_Int32_Ptr_Int64(f, C.uintptr_t(args[0]), C.int64_t(args[1]))), nil
	case funcType_Int32_Ptr_Ptr:
		return C.uint64_t(C.callFunc_Int32_Ptr_Ptr(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]))), nil
	case funcType_Int32_Ptr_Ptr_Ptr_Int32:
		return C.uint64_t(C.callFunc_Int32_Ptr_Ptr_Ptr_Int32(f, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.int32_t(args[3]))), nil
	case funcType_Int64_Ptr:
		return C.uint64_t(C.callFunc_Int64_Ptr(f, C.uintptr_t(args[0]))), nil
	case funcType_Ptr:
		return C.uint64_t(C.callFunc_Ptr(f)), nil
	case funcType_Ptr_Ptr:
		return C.uint64_t(C.callFunc_Ptr_Ptr(f, C.uintptr_t(args[0]))), nil
	case funcType_Void:
		C.callFunc_Void(f)
		return 0, nil
	case funcType_Void_Ptr_Bool:
		C.callFunc_Void_Ptr_Bool(f, C.uintptr_t(args[0]), C.uint8_t(args[1]))
		return 0, nil
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

func RestartAppIfNecessary(appID uint32) bool {
	v, err := theLib.call(funcType_Bool_Int32, flatAPI_RestartAppIfNecessary, uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func Init() error {
	var msg steamErrMsg
	v, err := theLib.call(funcType_Bool_Ptr, flatAPI_InitFlat, uintptr(unsafe.Pointer(&msg)))
	if err != nil {
		panic(err)
	}
	if ESteamAPIInitResult(v) != ESteamAPIInitResult_OK {
		return fmt.Errorf("steamworks: InitFlat failed: %d, %s", ESteamAPIInitResult(v), msg.String())
	}
	return nil
}

func RunCallbacks() {
	if _, err := theLib.call(funcType_Void, flatAPI_RunCallbacks); err != nil {
		panic(err)
	}
}

func SteamApps() ISteamApps {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamApps)
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps C.uintptr_t

func (s steamApps) BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool) {
	var name [4096]byte
	v, err := theLib.call(funcType_Int32_Ptr_Int32_Ptr_Ptr_Ptr_Int32, flatAPI_ISteamApps_BGetDLCDataByIndex, uintptr(s), uintptr(iDLC), uintptr(unsafe.Pointer(&appID)), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&name[0])), uintptr(len(name)))
	if err != nil {
		panic(err)
	}
	return appID, available, C.GoString((*C.char)(unsafe.Pointer(&name[0]))), byte(v) != 0
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	v, err := theLib.call(funcType_Bool_Ptr_Int32, flatAPI_ISteamApps_BIsDlcInstalled, uintptr(s), uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v, err := theLib.call(funcType_Int32_Ptr_Int32_Ptr_Int32, flatAPI_ISteamApps_GetAppInstallDir, uintptr(s), uintptr(appID), uintptr(unsafe.Pointer(&path[0])), uintptr(len(path)))
	if err != nil {
		panic(err)
	}
	return string(path[:uint32(v)-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	v, err := theLib.call(funcType_Ptr_Ptr, flatAPI_ISteamApps_GetCurrentGameLanguage, uintptr(s))
	if err != nil {
		panic(err)
	}
	return C.GoString(C.uintptrToChar(C.uintptr_t(v)))
}

func (s steamApps) GetDLCCount() int32 {
	v, err := theLib.call(funcType_Int32_Ptr, flatAPI_ISteamApps_GetDLCCount, uintptr(s))
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func SteamFriends() ISteamFriends {
	v, err := theLib.call(funcType_Ptr, flagAPI_SteamFriends)
	if err != nil {
		panic(err)
	}
	return steamFriends(v)
}

type steamFriends C.uintptr_t

func (s steamFriends) GetPersonaName() string {
	v, err := theLib.call(funcType_Ptr_Ptr, flatAPI_ISteamFriends_GetPersonaName, uintptr(s))
	if err != nil {
		panic(err)
	}
	return C.GoString(C.uintptrToChar(C.uintptr_t(v)))
}

func (s steamFriends) SetRichPresence(key, value string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	v, err := theLib.call(funcType_Bool_Ptr_Ptr_Ptr, flatAPI_ISteamFriends_SetRichPresence, uintptr(s), uintptr(unsafe.Pointer(ckey)), uintptr(unsafe.Pointer(cvalue)))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func SteamInput() ISteamInput {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamInput)
	if err != nil {
		panic(err)
	}
	return steamInput(v)
}

type steamInput C.uintptr_t

func (s steamInput) GetConnectedControllers() []InputHandle_t {
	var handles [_STEAM_INPUT_MAX_COUNT]InputHandle_t
	v, err := theLib.call(funcType_Int32_Ptr_Ptr, flatAPI_ISteamInput_GetConnectedControllers, uintptr(s), uintptr(unsafe.Pointer(&handles[0])))
	if err != nil {
		panic(err)
	}
	return handles[:int(v)]
}

func (s steamInput) GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType {
	v, err := theLib.call(funcType_Int32_Ptr_Int64, flatAPI_ISteamInput_GetInputTypeForHandle, uintptr(s), uintptr(inputHandle))
	if err != nil {
		panic(err)
	}
	return ESteamInputType(v)
}

func (s steamInput) Init(bExplicitlyCallRunFrame bool) bool {
	var callRunFrame uintptr
	if bExplicitlyCallRunFrame {
		callRunFrame = 1
	}
	v, err := theLib.call(funcType_Bool_Ptr_Bool, flatAPI_ISteamInput_Init, uintptr(s), callRunFrame)
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamInput) RunFrame() {
	if _, err := theLib.call(funcType_Void_Ptr_Bool, flatAPI_ISteamInput_RunFrame, uintptr(s), 0); err != nil {
		panic(err)
	}
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

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	v, err := theLib.call(funcType_Int32_Ptr, flatAPI_ISteamRemoteStorage_GetFileSize, uintptr(s), uintptr(unsafe.Pointer(cfile)))
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func SteamUser() ISteamUser {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamUser)
	if err != nil {
		panic(err)
	}
	return steamUser(v)
}

type steamUser C.uintptr_t

func (s steamUser) GetSteamID() CSteamID {
	v, err := theLib.call(funcType_Int64_Ptr, flatAPI_ISteamUser_GetSteamID, uintptr(s))
	if err != nil {
		panic(err)
	}
	return CSteamID(v)
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

func SteamUtils() ISteamUtils {
	v, err := theLib.call(funcType_Ptr, flatAPI_SteamUtils)
	if err != nil {
		panic(err)
	}
	return steamUtils(v)
}

type steamUtils C.uintptr_t

func (s steamUtils) IsSteamRunningOnSteamDeck() bool {
	v, err := theLib.call(funcType_Bool_Ptr, flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck, uintptr(s))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamUtils) ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool {
	v, err := theLib.call(funcType_Bool_Ptr_Int32_Int32_Int32_Int32_Int32, flatAPI_ISteamUtils_ShowFloatingGamepadTextInput, uintptr(s), uintptr(keyboardMode), uintptr(textFieldXPosition), uintptr(textFieldYPosition), uintptr(textFieldWidth), uintptr(textFieldHeight))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}
