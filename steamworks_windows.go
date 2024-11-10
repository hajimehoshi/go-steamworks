// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

package steamworks

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
typedef struct {
	unsigned long int m_steamIDUser;
	int m_nGlobalRank;
	int m_nScore;
	int m_cDetails;
	unsigned long int m_hUGC;
} LeaderboardEntry_t;
*/
import "C"

const is32Bit = unsafe.Sizeof(int(0)) == 4

func cStringToGoString(v uintptr, sizeHint int) string {
	bs := make([]byte, 0, sizeHint)
	for i := int32(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(v))
		v += unsafe.Sizeof(byte(0))
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs)
}

type dll struct {
	d     *windows.LazyDLL
	procs map[string]*windows.LazyProc
}

func (d *dll) call(name string, args ...uintptr) (uintptr, error) {
	if d.procs == nil {
		d.procs = map[string]*windows.LazyProc{}
	}
	if _, ok := d.procs[name]; !ok {
		d.procs[name] = d.d.NewProc(name)
	}
	r, _, err := d.procs[name].Call(args...)
	if err != nil {
		errno, ok := err.(windows.Errno)
		if !ok {
			return r, err
		}
		if errno != 0 {
			return r, err
		}
	}
	return r, nil
}

func loadDLL() (*dll, error) {
	dllName := "steam_api.dll"
	if !is32Bit {
		dllName = "steam_api64.dll"
	}

	return &dll{
		d: windows.NewLazyDLL(dllName),
	}, nil
}

var theDLL *dll

func init() {
	dll, err := loadDLL()
	if err != nil {
		panic(err)
	}
	theDLL = dll
}

func RestartAppIfNecessary(appID uint32) bool {
	v, err := theDLL.call(flatAPI_RestartAppIfNecessary, uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func Init() error {
	var msg steamErrMsg
	v, err := theDLL.call(flatAPI_InitFlat, uintptr(unsafe.Pointer(&msg[0])))
	if err != nil {
		panic(err)
	}
	if ESteamAPIInitResult(v) != ESteamAPIInitResult_OK {
		return fmt.Errorf("steamworks: InitFlat failed: %d, %s", ESteamAPIInitResult(v), msg.String())
	}
	return nil
}

func RunCallbacks() {
	if _, err := theDLL.call(flatAPI_RunCallbacks); err != nil {
		panic(err)
	}
}

func SteamApps() ISteamApps {
	v, err := theDLL.call(flatAPI_SteamApps)
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps uintptr

func (s steamApps) BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool) {
	var name [4096]byte
	v, err := theDLL.call(flatAPI_ISteamApps_BGetDLCDataByIndex, uintptr(s), uintptr(iDLC), uintptr(unsafe.Pointer(&appID)), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&name[0])), uintptr(len(name)))
	if err != nil {
		panic(err)
	}
	return appID, available, cStringToGoString(v, len(name)), byte(v) != 0
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	v, err := theDLL.call(flatAPI_ISteamApps_BIsDlcInstalled, uintptr(s), uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v, err := theDLL.call(flatAPI_ISteamApps_GetAppInstallDir, uintptr(s), uintptr(appID), uintptr(unsafe.Pointer(&path[0])), uintptr(len(path)))
	if err != nil {
		panic(err)
	}
	return string(path[:uint32(v)-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	v, err := theDLL.call(flatAPI_ISteamApps_GetCurrentGameLanguage, uintptr(s))
	if err != nil {
		panic(err)
	}
	return cStringToGoString(v, 256)
}

func (s steamApps) GetDLCCount() int32 {
	v, err := theDLL.call(flatAPI_ISteamApps_GetDLCCount, uintptr(s))
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func SteamFriends() ISteamFriends {
	v, err := theDLL.call(flagAPI_SteamFriends)
	if err != nil {
		panic(err)
	}
	return steamFriends(v)
}

type steamFriends uintptr

func (s steamFriends) GetPersonaName() string {
	v, err := theDLL.call(flatAPI_ISteamFriends_GetPersonaName, uintptr(s))
	if err != nil {
		panic(err)
	}
	return cStringToGoString(v, 64)
}

func (s steamFriends) SetRichPresence(key, value string) bool {
	ckey := append([]byte(key), 0)
	defer runtime.KeepAlive(ckey)
	cvalue := append([]byte(value), 0)
	defer runtime.KeepAlive(cvalue)

	v, err := theDLL.call(flatAPI_ISteamFriends_SetRichPresence, uintptr(s), uintptr(unsafe.Pointer(&ckey[0])), uintptr(unsafe.Pointer(&cvalue[0])))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func SteamInput() ISteamInput {
	v, err := theDLL.call(flatAPI_SteamInput)
	if err != nil {
		panic(err)
	}
	return steamInput(v)
}

type steamInput uintptr

func (s steamInput) GetConnectedControllers() []InputHandle_t {
	var handles [_STEAM_INPUT_MAX_COUNT]InputHandle_t
	v, err := theDLL.call(flatAPI_ISteamInput_GetConnectedControllers, uintptr(s), uintptr(unsafe.Pointer(&handles[0])))
	if err != nil {
		panic(err)
	}
	return handles[:int(v)]
}

func (s steamInput) GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType {
	v, err := theDLL.call(flatAPI_ISteamInput_GetInputTypeForHandle, uintptr(s), uintptr(inputHandle))
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
	// The error value seems unreliable.
	v, _ := theDLL.call(flatAPI_ISteamInput_Init, uintptr(s), callRunFrame)
	return byte(v) != 0
}

func (s steamInput) RunFrame() {
	if _, err := theDLL.call(flatAPI_ISteamInput_RunFrame, uintptr(s), 0); err != nil {
		panic(err)
	}
}

func SteamRemoteStorage() ISteamRemoteStorage {
	v, err := theDLL.call(flatAPI_SteamRemoteStorage)
	if err != nil {
		panic(err)
	}
	return steamRemoteStorage(v)
}

type steamRemoteStorage uintptr

func (s steamRemoteStorage) FileWrite(file string, data []byte) bool {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	defer runtime.KeepAlive(data)

	v, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileWrite, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	defer runtime.KeepAlive(data)

	v, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileRead, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}

	return int32(v)
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	v, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileDelete, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	v, err := theDLL.call(flatAPI_ISteamRemoteStorage_GetFileSize, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])))
	if err != nil {
		panic(err)
	}

	return int32(v)
}

func SteamUser() ISteamUser {
	v, err := theDLL.call(flatAPI_SteamUser)
	if err != nil {
		panic(err)
	}
	return steamUser(v)
}

type steamUser uintptr

func (s steamUser) GetSteamID() CSteamID {
	if is32Bit {
		// On 32bit machines, syscall cannot treat a returned value as 64bit.
		panic("GetSteamID is not implemented on 32bit Windows")
	}
	v, err := theDLL.call(flatAPI_ISteamUser_GetSteamID, uintptr(s))
	if err != nil {
		panic(err)
	}
	return CSteamID(v)
}

func SteamUserStats() ISteamUserStats {
	v, err := theDLL.call(flatAPI_SteamUserStats)
	if err != nil {
		panic(err)
	}
	return steamUserStats(v)
}

type steamUserStats uintptr

func (s steamUserStats) RequestCurrentStats() bool {
	v, err := theDLL.call(flatAPI_ISteamUserStats_RequestCurrentStats, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) GetAchievement(name string) (achieved, success bool) {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, err := theDLL.call(flatAPI_ISteamUserStats_GetAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])), uintptr(unsafe.Pointer(&achieved)))
	if err != nil {
		panic(err)
	}

	success = byte(v) != 0
	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, err := theDLL.call(flatAPI_ISteamUserStats_SetAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) ClearAchievement(name string) bool {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, err := theDLL.call(flatAPI_ISteamUserStats_ClearAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) StoreStats() bool {
	v, err := theDLL.call(flatAPI_ISteamUserStats_StoreStats, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) FindLeaderboard(leaderboardName string) SteamAPICall_t {
	cname := append([]byte(leaderboardName), 0)
	defer runtime.KeepAlive(cname)
	v, err := theDLL.call(flatAPI_ISteamUserStats_FindLeaderboard, uintptr(s), uintptr(unsafe.Pointer(&cname[0])))
	if err != nil {
		panic(err)
	}

	return SteamAPICall_t(v)
}

func (s steamUserStats) FindOrCreateLeaderboard(leaderboardName string, sortMethod ELeaderboardSortMethod, displayType ELeaderboardDisplayType) SteamAPICall_t {
	cname := append([]byte(leaderboardName), 0)
	defer runtime.KeepAlive(cname)
	v, err := theDLL.call(flatAPI_ISteamUserStats_FindOrCreateLeaderboard, uintptr(s), uintptr(unsafe.Pointer(&cname[0])), uintptr(sortMethod), uintptr(displayType))
	if err != nil {
		panic(err)
	}

	return SteamAPICall_t(v)
}
func (s steamUserStats) GetLeaderboardName(leaderboard SteamLeaderboard_t) string {
	v, err := theDLL.call(flatAPI_ISteamUserStats_GetLeaderboardName, uintptr(s), uintptr(leaderboard))
	if err != nil {
		panic(err)
	}
	return cStringToGoString(v, 64)
}

func (s steamUserStats) DownloadLeaderboardEntries(leaderboard SteamLeaderboard_t, dataRequest ELeaderboardDataRequest, rangeStart, rangeEnd int) SteamAPICall_t {
	v, err := theDLL.call(flatAPI_ISteamUserStats_DownloadLeaderboardEntries, uintptr(s), uintptr(leaderboard), uintptr(dataRequest), uintptr(rangeStart), uintptr(rangeEnd))
	if err != nil {
		panic(err)
	}

	return SteamAPICall_t(v)
}

func (s steamUserStats) DownloadLeaderboardEntriesForUsers(leaderboard SteamLeaderboard_t, prgUsers []CSteamID) SteamAPICall_t {
	prgUsers = append(prgUsers, 0)
	defer runtime.KeepAlive(prgUsers)
	v, err := theDLL.call(flatAPI_ISteamUserStats_DownloadLeaderboardEntriesForUsers, uintptr(s), uintptr(leaderboard), uintptr(unsafe.Pointer(&prgUsers[0])), uintptr(len(prgUsers)-1))
	if err != nil {
		panic(err)
	}
	return SteamAPICall_t(v)
}

func (s steamUserStats) GetDownloadedLeaderboardEntry(entries SteamLeaderboardEntries_t, index int, detailsMax int) (success bool, entry *LeaderboardEntry_t, details []int32) {
	details = append(details, 0)
	defer runtime.KeepAlive(details)
	var entryC C.LeaderboardEntry_t
	v, err := theDLL.call(flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry, uintptr(s), uintptr(entries), uintptr(index), uintptr(unsafe.Pointer(&entryC)), uintptr(unsafe.Pointer(&details[0])), uintptr(detailsMax))
	if err != nil {
		panic(err)
	}
	entry = &LeaderboardEntry_t{
		SteamIDUser: CSteamID(entryC.m_steamIDUser),
		GlobalRank:  int(entryC.m_nGlobalRank),
		Score:       int(entryC.m_nScore),
		Details:     int(entryC.m_cDetails),
		UGC:         UGCHandle_t(entryC.m_hUGC),
	}

	success = byte(v) != 0
	return
}

func (s steamUserStats) UploadLeaderboardScore(leaderboard SteamLeaderboard_t, uploadScoreMethod ELeaderboardUploadScoreMethod, score int32, scoreDetails ...int32) SteamAPICall_t {
	scoreDetails = append(scoreDetails, 0)
	defer runtime.KeepAlive(scoreDetails)
	v, err := theDLL.call(flatAPI_ISteamUserStats_UploadLeaderboardScore, uintptr(s), uintptr(leaderboard), uintptr(uploadScoreMethod), uintptr(score), uintptr(unsafe.Pointer(&scoreDetails[0])), uintptr(len(scoreDetails)-1))
	if err != nil {
		panic(err)
	}

	return SteamAPICall_t(v)
}

func SteamUtils() ISteamUtils {
	v, err := theDLL.call(flatAPI_SteamUtils)
	if err != nil {
		panic(err)
	}
	return steamUtils(v)
}

type steamUtils uintptr

func (s steamUtils) IsSteamRunningOnSteamDeck() bool {
	v, err := theDLL.call(flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUtils) ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool {
	v, err := theDLL.call(flatAPI_ISteamUtils_ShowFloatingGamepadTextInput, uintptr(s), uintptr(keyboardMode), uintptr(textFieldXPosition), uintptr(textFieldYPosition), uintptr(textFieldWidth), uintptr(textFieldHeight))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamUtils) GetAPICallResult(apiCall SteamAPICall_t, callbackExpected int, callbaseSize int) (callback []byte, success bool, pbFailed bool) {
	callback = make([]byte, callbaseSize)
	v, err := theDLL.call(flatAPI_ISteamUtils_GetAPICallResult, uintptr(s), uintptr(apiCall), uintptr(unsafe.Pointer(&callback[0])), uintptr(callbaseSize), uintptr(callbackExpected), uintptr(unsafe.Pointer(&pbFailed)))
	if err != nil {
		panic(err)
	}
	success = byte(v) != 0
	return
}
