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

	"github.com/ebitengine/purego"
)

type lib struct {
	lib uintptr
}

var (
	// Steam API function pointers - using different names to avoid conflicts with constants
	ptrAPI_RestartAppIfNecessary                    func(uint32) bool
	ptrAPI_InitFlat                                 func(uintptr) bool
	ptrAPI_RunCallbacks                             func()
	ptrAPI_SteamApps                                func() uintptr
	ptrAPI_ISteamApps_BGetDLCDataByIndex            func(uintptr, int32, uintptr, uintptr, uintptr, int32) bool
	ptrAPI_ISteamApps_BIsDlcInstalled               func(uintptr, uint32) bool
	ptrAPI_ISteamApps_GetAppInstallDir              func(uintptr, uint32, uintptr, int32) int32
	ptrAPI_ISteamApps_GetCurrentGameLanguage        func(uintptr) uintptr
	ptrAPI_ISteamApps_GetDLCCount                   func(uintptr) int32
	ptrAPI_SteamFriends                             func() uintptr
	ptrAPI_ISteamFriends_GetPersonaName             func(uintptr) uintptr
	ptrAPI_ISteamFriends_SetRichPresence            func(uintptr, uintptr, uintptr) bool
	ptrAPI_SteamInput                               func() uintptr
	ptrAPI_ISteamInput_GetConnectedControllers      func(uintptr, uintptr) int32
	ptrAPI_ISteamInput_GetInputTypeForHandle        func(uintptr, uint64) int32
	ptrAPI_ISteamInput_Init                         func(uintptr, bool) bool
	ptrAPI_ISteamInput_RunFrame                     func(uintptr, bool)
	ptrAPI_SteamRemoteStorage                       func() uintptr
	ptrAPI_ISteamRemoteStorage_FileWrite            func(uintptr, uintptr, uintptr, int32) bool
	ptrAPI_ISteamRemoteStorage_FileRead             func(uintptr, uintptr, uintptr, int32) int32
	ptrAPI_ISteamRemoteStorage_FileDelete           func(uintptr, uintptr) bool
	ptrAPI_ISteamRemoteStorage_GetFileSize          func(uintptr, uintptr) int32
	ptrAPI_SteamUser                                func() uintptr
	ptrAPI_ISteamUser_GetSteamID                    func(uintptr) uint64
	ptrAPI_SteamUserStats                           func() uintptr
	ptrAPI_ISteamUserStats_GetAchievement           func(uintptr, uintptr, uintptr) bool
	ptrAPI_ISteamUserStats_SetAchievement           func(uintptr, uintptr) bool
	ptrAPI_ISteamUserStats_ClearAchievement         func(uintptr, uintptr) bool
	ptrAPI_ISteamUserStats_StoreStats               func(uintptr) bool
	ptrAPI_SteamUtils                               func() uintptr
	ptrAPI_ISteamUtils_IsOverlayEnabled             func(uintptr) bool
	ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck    func(uintptr) bool
	ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput func(uintptr, int32, int32, int32, int32, int32) bool
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

	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		return 0, fmt.Errorf("steamworks: dlopen failed: %w", err)
	}

	// Register all Steam API function pointers
	purego.RegisterLibFunc(&ptrAPI_RestartAppIfNecessary, lib, flatAPI_RestartAppIfNecessary)
	purego.RegisterLibFunc(&ptrAPI_InitFlat, lib, flatAPI_InitFlat)
	purego.RegisterLibFunc(&ptrAPI_RunCallbacks, lib, flatAPI_RunCallbacks)
	purego.RegisterLibFunc(&ptrAPI_SteamApps, lib, flatAPI_SteamApps)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BGetDLCDataByIndex, lib, flatAPI_ISteamApps_BGetDLCDataByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsDlcInstalled, lib, flatAPI_ISteamApps_BIsDlcInstalled)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAppInstallDir, lib, flatAPI_ISteamApps_GetAppInstallDir)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetCurrentGameLanguage, lib, flatAPI_ISteamApps_GetCurrentGameLanguage)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetDLCCount, lib, flatAPI_ISteamApps_GetDLCCount)
	purego.RegisterLibFunc(&ptrAPI_SteamFriends, lib, flatAPI_SteamFriends)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetPersonaName, lib, flatAPI_ISteamFriends_GetPersonaName)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_SetRichPresence, lib, flatAPI_ISteamFriends_SetRichPresence)
	purego.RegisterLibFunc(&ptrAPI_SteamInput, lib, flatAPI_SteamInput)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetConnectedControllers, lib, flatAPI_ISteamInput_GetConnectedControllers)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetInputTypeForHandle, lib, flatAPI_ISteamInput_GetInputTypeForHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_Init, lib, flatAPI_ISteamInput_Init)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_RunFrame, lib, flatAPI_ISteamInput_RunFrame)
	purego.RegisterLibFunc(&ptrAPI_SteamRemoteStorage, lib, flatAPI_SteamRemoteStorage)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileWrite, lib, flatAPI_ISteamRemoteStorage_FileWrite)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileRead, lib, flatAPI_ISteamRemoteStorage_FileRead)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileDelete, lib, flatAPI_ISteamRemoteStorage_FileDelete)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_GetFileSize, lib, flatAPI_ISteamRemoteStorage_GetFileSize)
	purego.RegisterLibFunc(&ptrAPI_SteamUser, lib, flatAPI_SteamUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetSteamID, lib, flatAPI_ISteamUser_GetSteamID)
	purego.RegisterLibFunc(&ptrAPI_SteamUserStats, lib, flatAPI_SteamUserStats)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_GetAchievement, lib, flatAPI_ISteamUserStats_GetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_SetAchievement, lib, flatAPI_ISteamUserStats_SetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_ClearAchievement, lib, flatAPI_ISteamUserStats_ClearAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_StoreStats, lib, flatAPI_ISteamUserStats_StoreStats)
	purego.RegisterLibFunc(&ptrAPI_SteamUtils, lib, flatAPI_SteamUtils)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsOverlayEnabled, lib, flatAPI_ISteamUtils_IsOverlayEnabled)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck, lib, flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput, lib, flatAPI_ISteamUtils_ShowFloatingGamepadTextInput)

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

// Helper function to convert C string pointer to Go string using reflect
func cStringFromPtr(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	// Find the length of the C string
	length := 0
	for {
		if *(*byte)(unsafe.Pointer(ptr + uintptr(length))) == 0 {
			break
		}
		length++
		// Safety limit to prevent infinite loops
		if length > 65536 {
			break
		}
	}

	if length == 0 {
		return ""
	}

	// Create a byte slice header pointing to the C string
	var result []byte
	for i := 0; i < length; i++ {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		result = append(result, b)
	}
	return string(result)
}

// Helper function to convert Go string to C string
func goStringToC(s string) (uintptr, func()) {
	bytes := append([]byte(s), 0)
	ptr := uintptr(unsafe.Pointer(&bytes[0]))
	return ptr, func() { runtime.KeepAlive(bytes) }
}

func RestartAppIfNecessary(appID uint32) bool {
	return ptrAPI_RestartAppIfNecessary(appID)
}

func Init() error {
	var msg steamErrMsg
	if !ptrAPI_InitFlat(uintptr(unsafe.Pointer(&msg))) {
		return fmt.Errorf("steamworks: InitFlat failed: %s", msg.String())
	}
	return nil
}

func RunCallbacks() {
	ptrAPI_RunCallbacks()
}

func SteamApps() ISteamApps {
	return steamApps(ptrAPI_SteamApps())
}

type steamApps uintptr

func (s steamApps) BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool) {
	var name [4096]byte
	v := ptrAPI_ISteamApps_BGetDLCDataByIndex(uintptr(s), int32(iDLC), uintptr(unsafe.Pointer(&appID)), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&name[0])), int32(len(name)))
	return appID, available, string(name[:cStringLen(name[:])]), v
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	return ptrAPI_ISteamApps_BIsDlcInstalled(uintptr(s), uint32(appID))
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v := ptrAPI_ISteamApps_GetAppInstallDir(uintptr(s), uint32(appID), uintptr(unsafe.Pointer(&path[0])), int32(len(path)))
	if v == 0 {
		return ""
	}
	return string(path[:v-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	v := ptrAPI_ISteamApps_GetCurrentGameLanguage(uintptr(s))
	if v == 0 {
		return ""
	}
	return cStringFromPtr(v)
}

func (s steamApps) GetDLCCount() int32 {
	return ptrAPI_ISteamApps_GetDLCCount(uintptr(s))
}

// Helper function to find length of C string in byte array
func cStringLen(b []byte) int {
	for i, v := range b {
		if v == 0 {
			return i
		}
	}
	return len(b)
}

func SteamFriends() ISteamFriends {
	return steamFriends(ptrAPI_SteamFriends())
}

type steamFriends uintptr

func (s steamFriends) GetPersonaName() string {
	v := ptrAPI_ISteamFriends_GetPersonaName(uintptr(s))
	if v == 0 {
		return ""
	}
	return cStringFromPtr(v)
}

func (s steamFriends) SetRichPresence(key, value string) bool {
	keyPtr, keyCleanup := goStringToC(key)
	defer keyCleanup()
	valuePtr, valueCleanup := goStringToC(value)
	defer valueCleanup()
	return ptrAPI_ISteamFriends_SetRichPresence(uintptr(s), keyPtr, valuePtr)
}

func SteamInput() ISteamInput {
	return steamInput(ptrAPI_SteamInput())
}

type steamInput uintptr

func (s steamInput) GetConnectedControllers() []InputHandle_t {
	var handles [_STEAM_INPUT_MAX_COUNT]InputHandle_t
	v := ptrAPI_ISteamInput_GetConnectedControllers(uintptr(s), uintptr(unsafe.Pointer(&handles[0])))
	return handles[:int(v)]
}

func (s steamInput) GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType {
	v := ptrAPI_ISteamInput_GetInputTypeForHandle(uintptr(s), uint64(inputHandle))
	return ESteamInputType(v)
}

func (s steamInput) Init(bExplicitlyCallRunFrame bool) bool {
	return ptrAPI_ISteamInput_Init(uintptr(s), bExplicitlyCallRunFrame)
}

func (s steamInput) RunFrame() {
	ptrAPI_ISteamInput_RunFrame(uintptr(s), false)
}

func SteamRemoteStorage() ISteamRemoteStorage {
	return steamRemoteStorage(ptrAPI_SteamRemoteStorage())
}

type steamRemoteStorage uintptr

func (s steamRemoteStorage) FileWrite(file string, data []byte) bool {
	filePtr, fileCleanup := goStringToC(file)
	defer fileCleanup()
	runtime.KeepAlive(data)
	return ptrAPI_ISteamRemoteStorage_FileWrite(uintptr(s), filePtr, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	filePtr, fileCleanup := goStringToC(file)
	defer fileCleanup()
	runtime.KeepAlive(data)
	return ptrAPI_ISteamRemoteStorage_FileRead(uintptr(s), filePtr, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	filePtr, fileCleanup := goStringToC(file)
	defer fileCleanup()
	return ptrAPI_ISteamRemoteStorage_FileDelete(uintptr(s), filePtr)
}

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	filePtr, fileCleanup := goStringToC(file)
	defer fileCleanup()
	return ptrAPI_ISteamRemoteStorage_GetFileSize(uintptr(s), filePtr)
}

func SteamUser() ISteamUser {
	return steamUser(ptrAPI_SteamUser())
}

type steamUser uintptr

func (s steamUser) GetSteamID() CSteamID {
	return CSteamID(ptrAPI_ISteamUser_GetSteamID(uintptr(s)))
}

func SteamUserStats() ISteamUserStats {
	return steamUserStats(ptrAPI_SteamUserStats())
}

type steamUserStats uintptr

func (s steamUserStats) GetAchievement(name string) (achieved, success bool) {
	namePtr, nameCleanup := goStringToC(name)
	defer nameCleanup()
	success = ptrAPI_ISteamUserStats_GetAchievement(uintptr(s), namePtr, uintptr(unsafe.Pointer(&achieved)))
	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	namePtr, nameCleanup := goStringToC(name)
	defer nameCleanup()
	return ptrAPI_ISteamUserStats_SetAchievement(uintptr(s), namePtr)
}

func (s steamUserStats) ClearAchievement(name string) bool {
	namePtr, nameCleanup := goStringToC(name)
	defer nameCleanup()
	return ptrAPI_ISteamUserStats_ClearAchievement(uintptr(s), namePtr)
}

func (s steamUserStats) StoreStats() bool {
	return ptrAPI_ISteamUserStats_StoreStats(uintptr(s))
}

func SteamUtils() ISteamUtils {
	return steamUtils(ptrAPI_SteamUtils())
}

type steamUtils uintptr

func (s steamUtils) IsOverlayEnabled() bool {
	return ptrAPI_ISteamUtils_IsOverlayEnabled(uintptr(s))
}

func (s steamUtils) IsSteamRunningOnSteamDeck() bool {
	return ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck(uintptr(s))
}

func (s steamUtils) ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool {
	return ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput(uintptr(s), int32(keyboardMode), textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight)
}
