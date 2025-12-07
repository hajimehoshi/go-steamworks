// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

type AppId_t uint32
type CSteamID uint64
type InputHandle_t uint64
type steamAPICall_t uint64
type SteamLeaderboard_t uint64
type SteamLeaderboardEntries_t uint64
type UGCHandle_t uint64

const k_iSteamUserStatsCallbacks = 1100
const LeaderboardFindResult_k_iCallback = k_iSteamUserStatsCallbacks + 4
const LeaderboardScoresDownloaded_k_iCallback = k_iSteamUserStatsCallbacks + 5
const LeaderboardScoreUploaded_k_iCallback = k_iSteamUserStatsCallbacks + 6

type ESteamAPIInitResult int32

const (
	ESteamAPIInitResult_OK              ESteamAPIInitResult = 0
	ESteamAPIInitResult_FailedGeneric   ESteamAPIInitResult = 1
	ESteamAPIInitResult_NoSteamClient   ESteamAPIInitResult = 2
	ESteamAPIInitResult_VersionMismatch ESteamAPIInitResult = 3
)

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

type EFloatingGamepadTextInputMode int32

const (
	EFloatingGamepadTextInputMode_ModeSingleLine    EFloatingGamepadTextInputMode = 0
	EFloatingGamepadTextInputMode_ModeMultipleLines EFloatingGamepadTextInputMode = 1
	EFloatingGamepadTextInputMode_ModeEmail         EFloatingGamepadTextInputMode = 2
	EFloatingGamepadTextInputMode_ModeNumeric       EFloatingGamepadTextInputMode = 3
)

type ELeaderboardDataRequest int32

const (
	ELeaderboardDataRequestGlobal           ELeaderboardDataRequest = 0
	ELeaderboardDataRequestGlobalAroundUser ELeaderboardDataRequest = 1
	ELeaderboardDataRequestFriends          ELeaderboardDataRequest = 2
	ELeaderboardDataRequestUsers            ELeaderboardDataRequest = 3
)

type ELeaderboardDisplayType int32

const (
	ELeaderboardDisplayTypeNone             ELeaderboardDisplayType = 0
	ELeaderboardDisplayTypeNumeric          ELeaderboardDisplayType = 1
	ELeaderboardDisplayTypeTimeSeconds      ELeaderboardDisplayType = 2
	ELeaderboardDisplayTypeTimeMilliSeconds ELeaderboardDisplayType = 3
)

type ELeaderboardSortMethod int32

const (
	ELeaderboardSortMethodNone       ELeaderboardSortMethod = 0
	ELeaderboardSortMethodAscending  ELeaderboardSortMethod = 1
	ELeaderboardSortMethodDescending ELeaderboardSortMethod = 2
)

type ELeaderboardUploadScoreMethod int32

const (
	ELeaderboardUploadScoreMethodNone        ELeaderboardUploadScoreMethod = 0
	ELeaderboardUploadScoreMethodKeepBest    ELeaderboardUploadScoreMethod = 1
	ELeaderboardUploadScoreMethodForceUpdate ELeaderboardUploadScoreMethod = 2
)

type leaderboardFindResult struct {
	steamLeaderboard SteamLeaderboard_t
	leaderboardFound bool
}

type leaderboardScoresDownloaded struct {
	steamLeaderboard        SteamLeaderboard_t
	steamLeaderboardEntries SteamLeaderboardEntries_t
	entryCount              int32
}

type leaderboardScoreUploaded struct {
	success            bool
	steamLeaderboard   SteamLeaderboard_t
	score              int32
	scoreChanged       bool
	globalRankNew      int32
	globalRankPrevious int32
}

type leaderboardEntry struct {
	steamIDUser CSteamID
	globalRank  int32
	score       int32
	details     int32
	UGC         UGCHandle_t
}

// What go users see
type LeaderboardEntry struct {
	steamIDUser CSteamID
	globalRank  int32
	score       int32
	details     []int32
	UGC         UGCHandle_t
}

type LeaderboardScoreUploaded struct {
	nScore              int32
	bScoreChanged       bool
	nGlobalRankNew      int32
	nGlobalRankPrevious int32
}

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
	FindLeaderboard(name string, onComplete func(handle SteamLeaderboard_t, found bool, err error))
	DownloadLeaderboardEntries(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardDataRequest ELeaderboardDataRequest, nRangeStart, nRangeEnd int32, onComplete func(entries []LeaderboardEntry, err error))
	UploadLeaderboardScore(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardUploadScoreMethod ELeaderboardUploadScoreMethod, nScore int32, pScoreDetails []int32, onComplete func(result LeaderboardScoreUploaded, err error))
	GetLeaderboardEntryCount(hSteamLeaderboard SteamLeaderboard_t) int32
}

type ISteamUtils interface {
	IsOverlayEnabled() bool
	IsSteamRunningOnSteamDeck() bool
	ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool
}

type ISteamFriends interface {
	GetPersonaName() string
	SetRichPresence(string, string) bool
}

const (
	flatAPI_RestartAppIfNecessary = "SteamAPI_RestartAppIfNecessary"
	flatAPI_InitFlat              = "SteamAPI_InitFlat"
	flatAPI_RunCallbacks          = "SteamAPI_RunCallbacks"

	flatAPI_SteamApps                         = "SteamAPI_SteamApps_v008"
	flatAPI_ISteamApps_BGetDLCDataByIndex     = "SteamAPI_ISteamApps_BGetDLCDataByIndex"
	flatAPI_ISteamApps_BIsDlcInstalled        = "SteamAPI_ISteamApps_BIsDlcInstalled"
	flatAPI_ISteamApps_GetAppInstallDir       = "SteamAPI_ISteamApps_GetAppInstallDir"
	flatAPI_ISteamApps_GetCurrentGameLanguage = "SteamAPI_ISteamApps_GetCurrentGameLanguage"
	flatAPI_ISteamApps_GetDLCCount            = "SteamAPI_ISteamApps_GetDLCCount"

	flatAPI_SteamFriends                  = "SteamAPI_SteamFriends_v017"
	flatAPI_ISteamFriends_GetPersonaName  = "SteamAPI_ISteamFriends_GetPersonaName"
	flatAPI_ISteamFriends_SetRichPresence = "SteamAPI_ISteamFriends_SetRichPresence"

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

	flatAPI_SteamUser             = "SteamAPI_SteamUser_v023"
	flatAPI_ISteamUser_GetSteamID = "SteamAPI_ISteamUser_GetSteamID"

	flatAPI_SteamUserStats                                = "SteamAPI_SteamUserStats_v013"
	flatAPI_ISteamUserStats_GetAchievement                = "SteamAPI_ISteamUserStats_GetAchievement"
	flatAPI_ISteamUserStats_SetAchievement                = "SteamAPI_ISteamUserStats_SetAchievement"
	flatAPI_ISteamUserStats_ClearAchievement              = "SteamAPI_ISteamUserStats_ClearAchievement"
	flatAPI_ISteamUserStats_StoreStats                    = "SteamAPI_ISteamUserStats_StoreStats"
	flatAPI_ISteamUserStats_FindLeaderboard               = "SteamAPI_ISteamUserStats_FindLeaderboard"
	flatAPI_ISteamUserStats_DownloadLeaderboardEntries    = "SteamAPI_ISteamUserStats_DownloadLeaderboardEntries"
	flatAPI_ISteamUserStats_UploadLeaderboardScore        = "SteamAPI_ISteamUserStats_UploadLeaderboardScore"
	flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry = "SteamAPI_ISteamUserStats_GetDownloadedLeaderboardEntry"
	flatAPI_ISteamUserStats_GetLeaderboardEntryCount      = "SteamAPI_ISteamUserStats_GetLeaderboardEntryCount"

	flatAPI_SteamUtils                               = "SteamAPI_SteamUtils_v010"
	flatAPI_ISteamUtils_IsOverlayEnabled             = "SteamAPI_ISteamUtils_IsOverlayEnabled"
	flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck    = "SteamAPI_ISteamUtils_IsSteamRunningOnSteamDeck"
	flatAPI_ISteamUtils_ShowFloatingGamepadTextInput = "SteamAPI_ISteamUtils_ShowFloatingGamepadTextInput"
	flatAPI_ISteamUtils_GetAPICallResult             = "SteamAPI_ISteamUtils_GetAPICallResult"
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
