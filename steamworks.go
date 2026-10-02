// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

type AppId_t uint32
type CSteamID uint64
type InputHandle_t uint64

type ESteamAPIInitResult int32

const (
	ESteamAPIInitResult_OK              ESteamAPIInitResult = 0
	ESteamAPIInitResult_FailedGeneric   ESteamAPIInitResult = 1
	ESteamAPIInitResult_NoSteamClient   ESteamAPIInitResult = 2
	ESteamAPIInitResult_VersionMismatch ESteamAPIInitResult = 3
)

type ESteamInputType int32

const (
	ESteamInputType_Unknown                  ESteamInputType = 0
	ESteamInputType_SteamController          ESteamInputType = 1
	ESteamInputType_XBox360Controller        ESteamInputType = 2
	ESteamInputType_XBoxOneController        ESteamInputType = 3
	ESteamInputType_GenericGamepad           ESteamInputType = 4
	ESteamInputType_PS4Controller            ESteamInputType = 5
	ESteamInputType_AppleMFiController       ESteamInputType = 6 // Unused
	ESteamInputType_AndroidController        ESteamInputType = 7 // Unused
	ESteamInputType_SwitchJoyConPair         ESteamInputType = 8 // Unused
	ESteamInputType_SwitchJoyConSingle       ESteamInputType = 9 // Unused
	ESteamInputType_SwitchProController      ESteamInputType = 10
	ESteamInputType_MobileTouch              ESteamInputType = 11
	ESteamInputType_PS3Controller            ESteamInputType = 12
	ESteamInputType_PS5Controller            ESteamInputType = 13
	ESteamInputType_SteamDeckController      ESteamInputType = 14
	ESteamInputType_SteamOSHandheld          ESteamInputType = 15
	ESteamInputType_Switch2ProController     ESteamInputType = 16
	ESteamInputType_SteamController2026      ESteamInputType = 17
	ESteamInputType_SteamFrameControllerPair ESteamInputType = 18
	ESteamInputType_Count                    ESteamInputType = 19
	ESteamInputType_MaximumPossibleValue     ESteamInputType = 255
)

// Deprecated: Use ESteamInputType_GenericGamepad.
const ESteamInputType_GenericXInput = ESteamInputType_GenericGamepad

const (
	_STEAM_INPUT_MAX_COUNT = 16
)

type EFloatingGamepadTextInputMode int32

const (
	EFloatingGamepadTextInputMode_ModeSingleLine    EFloatingGamepadTextInputMode = 0
	EFloatingGamepadTextInputMode_ModeMultipleLines EFloatingGamepadTextInputMode = 1
	EFloatingGamepadTextInputMode_ModeEmail         EFloatingGamepadTextInputMode = 2
	EFloatingGamepadTextInputMode_ModeNumeric       EFloatingGamepadTextInputMode = 3
)

type ISteamApps interface {
	BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool)
	BIsDlcInstalled(appID AppId_t) bool
	GetAppInstallDir(appID AppId_t) string
	GetCurrentGameLanguage() string
	GetDLCCount() int32
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
	GetAchievement(name string) (achieved, success bool)
	SetAchievement(name string) bool
	ClearAchievement(name string) bool
	StoreStats() bool
}

type ISteamUtils interface {
	IsOverlayEnabled() bool
	IsRunningOnSteamHardware() ESteamHardwareType
	GetSteamHardwareDefaultConfig() ESteamHardwareDefaultConfig
	IsRunningUnderProton() bool
	// Deprecated: Use IsRunningOnSteamHardware.
	IsSteamRunningOnSteamDeck() bool
	ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool
}

type ESteamHardwareType int32

const (
	ESteamHardwareType_None         ESteamHardwareType = 0
	ESteamHardwareType_SteamDeck    ESteamHardwareType = 1
	ESteamHardwareType_SteamMachine ESteamHardwareType = 2
	ESteamHardwareType_SteamFrame   ESteamHardwareType = 3
)

type ESteamHardwareDefaultConfig int32

const (
	ESteamHardwareDefaultConfig_None         ESteamHardwareDefaultConfig = 0
	ESteamHardwareDefaultConfig_Low          ESteamHardwareDefaultConfig = 1
	ESteamHardwareDefaultConfig_Medium       ESteamHardwareDefaultConfig = 2
	ESteamHardwareDefaultConfig_High         ESteamHardwareDefaultConfig = 3
	ESteamHardwareDefaultConfig_Max          ESteamHardwareDefaultConfig = 4
	ESteamHardwareDefaultConfig_SteamDeck    ESteamHardwareDefaultConfig = 5
	ESteamHardwareDefaultConfig_SteamMachine ESteamHardwareDefaultConfig = 6
	ESteamHardwareDefaultConfig_SteamFrame   ESteamHardwareDefaultConfig = 7
)

type ISteamFriends interface {
	GetPersonaName() string
	SetRichPresence(string, string) bool
}

const (
	flatAPI_RestartAppIfNecessary = "SteamAPI_RestartAppIfNecessary"
	flatAPI_InitFlat              = "SteamAPI_InitFlat"
	flatAPI_RunCallbacks          = "SteamAPI_RunCallbacks"

	flatAPI_SteamApps                         = "SteamAPI_SteamApps_v009"
	flatAPI_ISteamApps_BGetDLCDataByIndex     = "SteamAPI_ISteamApps_BGetDLCDataByIndex"
	flatAPI_ISteamApps_BIsDlcInstalled        = "SteamAPI_ISteamApps_BIsDlcInstalled"
	flatAPI_ISteamApps_GetAppInstallDir       = "SteamAPI_ISteamApps_GetAppInstallDir"
	flatAPI_ISteamApps_GetCurrentGameLanguage = "SteamAPI_ISteamApps_GetCurrentGameLanguage"
	flatAPI_ISteamApps_GetDLCCount            = "SteamAPI_ISteamApps_GetDLCCount"

	flatAPI_SteamFriends                  = "SteamAPI_SteamFriends_v018"
	flatAPI_ISteamFriends_GetPersonaName  = "SteamAPI_ISteamFriends_GetPersonaName"
	flatAPI_ISteamFriends_SetRichPresence = "SteamAPI_ISteamFriends_SetRichPresence"

	flatAPI_SteamInput                          = "SteamAPI_SteamInput_v007"
	flatAPI_ISteamInput_GetConnectedControllers = "SteamAPI_ISteamInput_GetConnectedControllers"
	flatAPI_ISteamInput_GetInputTypeForHandle   = "SteamAPI_ISteamInput_GetInputTypeForHandle"
	flatAPI_ISteamInput_Init                    = "SteamAPI_ISteamInput_Init"
	flatAPI_ISteamInput_RunFrame                = "SteamAPI_ISteamInput_RunFrame"

	flatAPI_SteamRemoteStorage              = "SteamAPI_SteamRemoteStorage_v016"
	flatAPI_ISteamRemoteStorage_FileWrite   = "SteamAPI_ISteamRemoteStorage_FileWrite"
	flatAPI_ISteamRemoteStorage_FileRead    = "SteamAPI_ISteamRemoteStorage_FileRead"
	flatAPI_ISteamRemoteStorage_FileDelete  = "SteamAPI_ISteamRemoteStorage_FileDelete"
	flatAPI_ISteamRemoteStorage_GetFileSize = "SteamAPI_ISteamRemoteStorage_GetFileSize"

	flatAPI_SteamUser             = "SteamAPI_SteamUser_v023"
	flatAPI_ISteamUser_GetSteamID = "SteamAPI_ISteamUser_GetSteamID"

	flatAPI_SteamUserStats                   = "SteamAPI_SteamUserStats_v013"
	flatAPI_ISteamUserStats_GetAchievement   = "SteamAPI_ISteamUserStats_GetAchievement"
	flatAPI_ISteamUserStats_SetAchievement   = "SteamAPI_ISteamUserStats_SetAchievement"
	flatAPI_ISteamUserStats_ClearAchievement = "SteamAPI_ISteamUserStats_ClearAchievement"
	flatAPI_ISteamUserStats_StoreStats       = "SteamAPI_ISteamUserStats_StoreStats"

	flatAPI_SteamUtils                                = "SteamAPI_SteamUtils_v011"
	flatAPI_ISteamUtils_IsOverlayEnabled              = "SteamAPI_ISteamUtils_IsOverlayEnabled"
	flatAPI_ISteamUtils_IsRunningOnSteamHardware      = "SteamAPI_ISteamUtils_IsRunningOnSteamHardware"
	flatAPI_ISteamUtils_GetSteamHardwareDefaultConfig = "SteamAPI_ISteamUtils_GetSteamHardwareDefaultConfig"
	flatAPI_ISteamUtils_IsRunningUnderProton          = "SteamAPI_ISteamUtils_IsRunningUnderProton"
	flatAPI_ISteamUtils_ShowFloatingGamepadTextInput  = "SteamAPI_ISteamUtils_ShowFloatingGamepadTextInput"
)

type steamErrMsg [1024]byte

func (s *steamErrMsg) String() string {
	for i, b := range s {
		if b == 0 {
			return string(s[:i])
		}
	}
	return ""
}
