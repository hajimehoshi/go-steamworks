// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package steamworks

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

type lib struct {
	lib uintptr
}

var (
	// General
	ptrAPI_RestartAppIfNecessary func(uint32) bool
	ptrAPI_InitFlat              func(uintptr) ESteamAPIInitResult
	ptrAPI_RunCallbacks          func()

	// ISteamApps
	ptrAPI_SteamApps                         func() uintptr
	ptrAPI_ISteamApps_BGetDLCDataByIndex     func(uintptr, int32, uintptr, uintptr, uintptr, int32) bool
	ptrAPI_ISteamApps_BIsDlcInstalled        func(uintptr, AppId_t) bool
	ptrAPI_ISteamApps_GetAppInstallDir       func(uintptr, AppId_t, uintptr, int32) int32
	ptrAPI_ISteamApps_GetCurrentGameLanguage func(uintptr) string
	ptrAPI_ISteamApps_GetDLCCount            func(uintptr) int32

	// ISteamFriends
	ptrAPI_SteamFriends                  func() uintptr
	ptrAPI_ISteamFriends_GetPersonaName  func(uintptr) string
	ptrAPI_ISteamFriends_SetRichPresence func(uintptr, string, string) bool

	// ISteamInput
	ptrAPI_SteamInput                          func() uintptr
	ptrAPI_ISteamInput_GetConnectedControllers func(uintptr, uintptr) int32
	ptrAPI_ISteamInput_GetInputTypeForHandle   func(uintptr, InputHandle_t) int32
	ptrAPI_ISteamInput_Init                    func(uintptr, bool) bool
	ptrAPI_ISteamInput_RunFrame                func(uintptr, bool)

	// ISteamRemoteStorage
	ptrAPI_SteamRemoteStorage              func() uintptr
	ptrAPI_ISteamRemoteStorage_FileWrite   func(uintptr, string, uintptr, int32) bool
	ptrAPI_ISteamRemoteStorage_FileRead    func(uintptr, string, uintptr, int32) int32
	ptrAPI_ISteamRemoteStorage_FileDelete  func(uintptr, string) bool
	ptrAPI_ISteamRemoteStorage_GetFileSize func(uintptr, string) int32

	// ISteamUser
	ptrAPI_SteamUser             func() uintptr
	ptrAPI_ISteamUser_GetSteamID func(uintptr) CSteamID

	// ISteamUserStats
	ptrAPI_SteamUserStats                   func() uintptr
	ptrAPI_ISteamUserStats_GetAchievement   func(uintptr, string, uintptr) bool
	ptrAPI_ISteamUserStats_SetAchievement   func(uintptr, string) bool
	ptrAPI_ISteamUserStats_ClearAchievement func(uintptr, string) bool
	ptrAPI_ISteamUserStats_StoreStats       func(uintptr) bool

	// ISteamUtils
	ptrAPI_SteamUtils                               func() uintptr
	ptrAPI_ISteamUtils_IsOverlayEnabled             func(uintptr) bool
	ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck    func(uintptr) bool
	ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput func(uintptr, EFloatingGamepadTextInputMode, int32, int32, int32, int32) bool
)

func registerFunctions(lib uintptr) {
	// General
	purego.RegisterLibFunc(&ptrAPI_RestartAppIfNecessary, lib, flatAPI_RestartAppIfNecessary)
	purego.RegisterLibFunc(&ptrAPI_InitFlat, lib, flatAPI_InitFlat)
	purego.RegisterLibFunc(&ptrAPI_RunCallbacks, lib, flatAPI_RunCallbacks)

	// ISteamApps
	purego.RegisterLibFunc(&ptrAPI_SteamApps, lib, flatAPI_SteamApps)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BGetDLCDataByIndex, lib, flatAPI_ISteamApps_BGetDLCDataByIndex)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_BIsDlcInstalled, lib, flatAPI_ISteamApps_BIsDlcInstalled)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetAppInstallDir, lib, flatAPI_ISteamApps_GetAppInstallDir)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetCurrentGameLanguage, lib, flatAPI_ISteamApps_GetCurrentGameLanguage)
	purego.RegisterLibFunc(&ptrAPI_ISteamApps_GetDLCCount, lib, flatAPI_ISteamApps_GetDLCCount)

	// ISteamFriends
	purego.RegisterLibFunc(&ptrAPI_SteamFriends, lib, flatAPI_SteamFriends)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_GetPersonaName, lib, flatAPI_ISteamFriends_GetPersonaName)
	purego.RegisterLibFunc(&ptrAPI_ISteamFriends_SetRichPresence, lib, flatAPI_ISteamFriends_SetRichPresence)

	// ISteamInput
	purego.RegisterLibFunc(&ptrAPI_SteamInput, lib, flatAPI_SteamInput)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetConnectedControllers, lib, flatAPI_ISteamInput_GetConnectedControllers)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_GetInputTypeForHandle, lib, flatAPI_ISteamInput_GetInputTypeForHandle)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_Init, lib, flatAPI_ISteamInput_Init)
	purego.RegisterLibFunc(&ptrAPI_ISteamInput_RunFrame, lib, flatAPI_ISteamInput_RunFrame)

	// ISteamRemoteStorage
	purego.RegisterLibFunc(&ptrAPI_SteamRemoteStorage, lib, flatAPI_SteamRemoteStorage)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileWrite, lib, flatAPI_ISteamRemoteStorage_FileWrite)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileRead, lib, flatAPI_ISteamRemoteStorage_FileRead)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_FileDelete, lib, flatAPI_ISteamRemoteStorage_FileDelete)
	purego.RegisterLibFunc(&ptrAPI_ISteamRemoteStorage_GetFileSize, lib, flatAPI_ISteamRemoteStorage_GetFileSize)

	// ISteamUser
	purego.RegisterLibFunc(&ptrAPI_SteamUser, lib, flatAPI_SteamUser)
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetSteamID, lib, flatAPI_ISteamUser_GetSteamID)

	// ISteamUserStats
	purego.RegisterLibFunc(&ptrAPI_SteamUserStats, lib, flatAPI_SteamUserStats)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_GetAchievement, lib, flatAPI_ISteamUserStats_GetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_SetAchievement, lib, flatAPI_ISteamUserStats_SetAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_ClearAchievement, lib, flatAPI_ISteamUserStats_ClearAchievement)
	purego.RegisterLibFunc(&ptrAPI_ISteamUserStats_StoreStats, lib, flatAPI_ISteamUserStats_StoreStats)

	// ISteamUtils
	purego.RegisterLibFunc(&ptrAPI_SteamUtils, lib, flatAPI_SteamUtils)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsOverlayEnabled, lib, flatAPI_ISteamUtils_IsOverlayEnabled)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_IsSteamRunningOnSteamDeck, lib, flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck)
	purego.RegisterLibFunc(&ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput, lib, flatAPI_ISteamUtils_ShowFloatingGamepadTextInput)
}

var theLib *lib

func init() {
	l, err := loadLib()
	if err != nil {
		panic(err)
	}
	registerFunctions(l)
	theLib = &lib{
		lib: l,
	}
}

func RestartAppIfNecessary(appID uint32) bool {
	return ptrAPI_RestartAppIfNecessary(appID)
}

func Init() error {
	var msg steamErrMsg
	if ptrAPI_InitFlat(uintptr(unsafe.Pointer(&msg))) != ESteamAPIInitResult_OK {
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
	return appID, available, CStringToGo(name[:]), v
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	return ptrAPI_ISteamApps_BIsDlcInstalled(uintptr(s), appID)
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v := ptrAPI_ISteamApps_GetAppInstallDir(uintptr(s), appID, uintptr(unsafe.Pointer(&path[0])), int32(len(path)))
	if v == 0 {
		return ""
	}
	return string(path[:v-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	return ptrAPI_ISteamApps_GetCurrentGameLanguage(uintptr(s))
}

func (s steamApps) GetDLCCount() int32 {
	return ptrAPI_ISteamApps_GetDLCCount(uintptr(s))
}

func SteamFriends() ISteamFriends {
	return steamFriends(ptrAPI_SteamFriends())
}

type steamFriends uintptr

func (s steamFriends) GetPersonaName() string {
	return ptrAPI_ISteamFriends_GetPersonaName(uintptr(s))
}

func (s steamFriends) SetRichPresence(key, value string) bool {
	return ptrAPI_ISteamFriends_SetRichPresence(uintptr(s), key, value)
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
	v := ptrAPI_ISteamInput_GetInputTypeForHandle(uintptr(s), inputHandle)
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
	return ptrAPI_ISteamRemoteStorage_FileWrite(uintptr(s), file, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	return ptrAPI_ISteamRemoteStorage_FileRead(uintptr(s), file, uintptr(unsafe.Pointer(&data[0])), int32(len(data)))
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	return ptrAPI_ISteamRemoteStorage_FileDelete(uintptr(s), file)
}

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	return ptrAPI_ISteamRemoteStorage_GetFileSize(uintptr(s), file)
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
	success = ptrAPI_ISteamUserStats_GetAchievement(uintptr(s), name, uintptr(unsafe.Pointer(&achieved)))
	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	return ptrAPI_ISteamUserStats_SetAchievement(uintptr(s), name)
}

func (s steamUserStats) ClearAchievement(name string) bool {
	return ptrAPI_ISteamUserStats_ClearAchievement(uintptr(s), name)
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
	return ptrAPI_ISteamUtils_ShowFloatingGamepadTextInput(uintptr(s), keyboardMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight)
}

func CStringToGo(name []byte) string {
	index := bytes.IndexByte(name, 0)
	var nameResult string
	if index < 0 {
		// No null terminator detected, so use the whole result
		nameResult = string(name)
	} else {
		// Null terminator detected, so use up to that point, excluding the null terminator
		nameResult = string(name[:index])
	}
	return nameResult
}
