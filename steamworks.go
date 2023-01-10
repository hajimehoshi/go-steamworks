// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

type AppId_t uint32
type CSteamID uint64
type InputHandle_t uint64

type ESteamInputType int32

const (
	ESteamInputType_Unknown              ESteamInputType = 0
	ESteamInputType_SteamController      ESteamInputType = 1
	ESteamInputType_XBox360Controller    ESteamInputType = 2
	ESteamInputType_XBoxOneController    ESteamInputType = 3
	ESteamInputType_GenericXInput        ESteamInputType = 4
	ESteamInputType_PS4Controller        ESteamInputType = 5
	ESteamInputType_AppleMFiController   ESteamInputType = 6 // Unused
	ESteamInputType_AndroidController    ESteamInputType = 7 // Unused
	ESteamInputType_SwitchJoyConPair     ESteamInputType = 8 // Unused
	ESteamInputType_SwitchJoyConSingle   ESteamInputType = 9 // Unused
	ESteamInputType_SwitchProController  ESteamInputType = 10
	ESteamInputType_MobileTouch          ESteamInputType = 11
	ESteamInputType_PS3Controller        ESteamInputType = 12
	ESteamInputType_PS5Controller        ESteamInputType = 13
	ESteamInputType_SteamDeckController  ESteamInputType = 14
	ESteamInputType_Count                ESteamInputType = 15
	ESteamInputType_MaximumPossibleValue ESteamInputType = 255
)

const (
	_STEAM_INPUT_MAX_COUNT = 16
)

type ISteamApps interface {
	GetAppInstallDir(appID AppId_t) string
	GetCurrentGameLanguage() string
}

type ISteamInput interface {
	GetConnectedControllers() []InputHandle_t
	GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType
	Init(bExplicitlyCallRunFrame bool) bool
	RunFrame()
}

type ISteamRemoteStorage interface {
	FileWrite(file string, data []byte) bool
	FileRead(file string, data []byte) int32
	FileDelete(file string) bool
	GetFileSize(file string) int32
}

type ISteamUser interface {
	GetSteamID() CSteamID
}

type ISteamUserStats interface {
	RequestCurrentStats() bool
	GetAchievement(name string) (achieved, success bool)
	SetAchievement(name string) bool
	ClearAchievement(name string) bool
	StoreStats() bool
}

const (
	flatAPI_RestartAppIfNecessary = "SteamAPI_RestartAppIfNecessary"
	flatAPI_Init                  = "SteamAPI_Init"
	flatAPI_RunCallbacks          = "SteamAPI_RunCallbacks"

	flatAPI_SteamApps                         = "SteamAPI_SteamApps_v008"
	flatAPI_ISteamApps_GetAppInstallDir       = "SteamAPI_ISteamApps_GetAppInstallDir"
	flatAPI_ISteamApps_GetCurrentGameLanguage = "SteamAPI_ISteamApps_GetCurrentGameLanguage"

	flatAPI_SteamInput                          = "SteamAPI_SteamInput_v006"
	flatAPI_ISteamInput_GetConnectedControllers = "SteamAPI_ISteamInput_GetConnectedControllers"
	flatAPI_ISteamInput_GetInputTypeForHandle   = "SteamAPI_ISteamInput_GetInputTypeForHandle"
	flatAPI_ISteamInput_Init                    = "SteamAPI_ISteamInput_Init"
	flatAPI_ISteamInput_RunFrame                = "SteamAPI_ISteamInput_RunFrame"

	flatAPI_SteamRemoteStorage              = "SteamAPI_SteamRemoteStorage_v016"
	flatAPI_ISteamRemoteStorage_FileWrite   = "SteamAPI_ISteamRemoteStorage_FileWrite"
	flatAPI_ISteamRemoteStorage_FileRead    = "SteamAPI_ISteamRemoteStorage_FileRead"
	flatAPI_ISteamRemoteStorage_FileDelete  = "SteamAPI_ISteamRemoteStorage_FileDelete"
	flatAPI_ISteamRemoteStorage_GetFileSize = "SteamAPI_ISteamRemoteStorage_GetFileSize"

	flatAPI_SteamUser             = "SteamAPI_SteamUser_v021"
	flatAPI_ISteamUser_GetSteamID = "SteamAPI_ISteamUser_GetSteamID"

	flatAPI_SteamUserStats                      = "SteamAPI_SteamUserStats_v012"
	flatAPI_ISteamUserStats_RequestCurrentStats = "SteamAPI_ISteamUserStats_RequestCurrentStats"
	flatAPI_ISteamUserStats_GetAchievement      = "SteamAPI_ISteamUserStats_GetAchievement"
	flatAPI_ISteamUserStats_SetAchievement      = "SteamAPI_ISteamUserStats_SetAchievement"
	flatAPI_ISteamUserStats_ClearAchievement    = "SteamAPI_ISteamUserStats_ClearAchievement"
	flatAPI_ISteamUserStats_StoreStats          = "SteamAPI_ISteamUserStats_StoreStats"
)
