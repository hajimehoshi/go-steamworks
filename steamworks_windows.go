// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021 The go-steamworks Authors

package steamworks

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

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

func (d *dll) call(name string, args ...uintptr) (uintptr, uintptr, error) {
	if d.procs == nil {
		d.procs = map[string]*windows.LazyProc{}
	}
	if _, ok := d.procs[name]; !ok {
		d.procs[name] = d.d.NewProc(name)
	}
	r1, r2, err := d.procs[name].Call(args...)
	if err != nil {
		errno, ok := err.(windows.Errno)
		if !ok {
			return r1, r2, err
		}
		if errno != 0 {
			return r1, r2, err
		}
	}
	return r1, r2, nil
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
	v, _, err := theDLL.call(flatAPI_RestartAppIfNecessary, uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func Init() error {
	var msg steamErrMsg
	v, _, err := theDLL.call(flatAPI_InitFlat, uintptr(unsafe.Pointer(&msg[0])))
	if err != nil {
		panic(err)
	}
	if ESteamAPIInitResult(v) != ESteamAPIInitResult_OK {
		return fmt.Errorf("steamworks: InitFlat failed: %d, %s", ESteamAPIInitResult(v), msg.String())
	}
	return nil
}

func runCallbacksSteam() {
	if _, _, err := theDLL.call(flatAPI_RunCallbacks); err != nil {
		panic(err)
	}
}

func SteamApps() ISteamApps {
	v, _, err := theDLL.call(flatAPI_SteamApps)
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps uintptr

func (s steamApps) BGetDLCDataByIndex(iDLC int) (appID AppId_t, available bool, pchName string, success bool) {
	var name [4096]byte
	v, _, err := theDLL.call(flatAPI_ISteamApps_BGetDLCDataByIndex, uintptr(s), uintptr(iDLC), uintptr(unsafe.Pointer(&appID)), uintptr(unsafe.Pointer(&available)), uintptr(unsafe.Pointer(&name[0])), uintptr(len(name)))
	if err != nil {
		panic(err)
	}
	return appID, available, cStringToGoString(uintptr(unsafe.Pointer(&name[0])), len(name)), byte(v) != 0
}

func (s steamApps) BIsDlcInstalled(appID AppId_t) bool {
	v, _, err := theDLL.call(flatAPI_ISteamApps_BIsDlcInstalled, uintptr(s), uintptr(appID))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func (s steamApps) GetAppInstallDir(appID AppId_t) string {
	var path [4096]byte
	v, _, err := theDLL.call(flatAPI_ISteamApps_GetAppInstallDir, uintptr(s), uintptr(appID), uintptr(unsafe.Pointer(&path[0])), uintptr(len(path)))
	if err != nil {
		panic(err)
	}
	if uint32(v) == 0 {
		return ""
	}
	return string(path[:uint32(v)-1])
}

func (s steamApps) GetCurrentGameLanguage() string {
	v, _, err := theDLL.call(flatAPI_ISteamApps_GetCurrentGameLanguage, uintptr(s))
	if err != nil {
		panic(err)
	}
	return cStringToGoString(v, 256)
}

func (s steamApps) GetDLCCount() int32 {
	v, _, err := theDLL.call(flatAPI_ISteamApps_GetDLCCount, uintptr(s))
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func SteamFriends() ISteamFriends {
	v, _, err := theDLL.call(flagAPI_SteamFriends)
	if err != nil {
		panic(err)
	}
	return steamFriends(v)
}

type steamFriends uintptr

func (s steamFriends) GetPersonaName() string {
	v, _, err := theDLL.call(flatAPI_ISteamFriends_GetPersonaName, uintptr(s))
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

	v, _, err := theDLL.call(flatAPI_ISteamFriends_SetRichPresence, uintptr(s), uintptr(unsafe.Pointer(&ckey[0])), uintptr(unsafe.Pointer(&cvalue[0])))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

func SteamInput() ISteamInput {
	v, _, err := theDLL.call(flatAPI_SteamInput)
	if err != nil {
		panic(err)
	}
	return steamInput(v)
}

type steamInput uintptr

func (s steamInput) GetConnectedControllers() []InputHandle_t {
	var handles [_STEAM_INPUT_MAX_COUNT]InputHandle_t
	v, _, err := theDLL.call(flatAPI_ISteamInput_GetConnectedControllers, uintptr(s), uintptr(unsafe.Pointer(&handles[0])))
	if err != nil {
		panic(err)
	}
	return handles[:int(v)]
}

func (s steamInput) GetInputTypeForHandle(inputHandle InputHandle_t) ESteamInputType {
	v, _, err := theDLL.call(flatAPI_ISteamInput_GetInputTypeForHandle, uintptr(s), uintptr(inputHandle))
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
	v, _, _ := theDLL.call(flatAPI_ISteamInput_Init, uintptr(s), callRunFrame)
	return byte(v) != 0
}

func (s steamInput) RunFrame() {
	if _, _, err := theDLL.call(flatAPI_ISteamInput_RunFrame, uintptr(s), 0); err != nil {
		panic(err)
	}
}

func SteamRemoteStorage() ISteamRemoteStorage {
	v, _, err := theDLL.call(flatAPI_SteamRemoteStorage)
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

	v, _, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileWrite, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamRemoteStorage) FileRead(file string, data []byte) int32 {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	defer runtime.KeepAlive(data)

	v, _, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileRead, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if err != nil {
		panic(err)
	}

	return int32(v)
}

func (s steamRemoteStorage) FileDelete(file string) bool {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	v, _, err := theDLL.call(flatAPI_ISteamRemoteStorage_FileDelete, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamRemoteStorage) GetFileSize(file string) int32 {
	cfile := append([]byte(file), 0)
	defer runtime.KeepAlive(cfile)

	v, _, err := theDLL.call(flatAPI_ISteamRemoteStorage_GetFileSize, uintptr(s), uintptr(unsafe.Pointer(&cfile[0])))
	if err != nil {
		panic(err)
	}

	return int32(v)
}

func SteamUser() ISteamUser {
	v, _, err := theDLL.call(flatAPI_SteamUser)
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
	v, _, err := theDLL.call(flatAPI_ISteamUser_GetSteamID, uintptr(s))
	if err != nil {
		panic(err)
	}
	return CSteamID(v)
}

func SteamUserStats() ISteamUserStats {
	v, _, err := theDLL.call(flatAPI_SteamUserStats)
	if err != nil {
		panic(err)
	}
	return steamUserStats(v)
}

type steamUserStats uintptr

func (s steamUserStats) GetAchievement(name string) (achieved, success bool) {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, _, err := theDLL.call(flatAPI_ISteamUserStats_GetAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])), uintptr(unsafe.Pointer(&achieved)))
	if err != nil {
		panic(err)
	}

	success = byte(v) != 0
	return
}

func (s steamUserStats) SetAchievement(name string) bool {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, _, err := theDLL.call(flatAPI_ISteamUserStats_SetAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) ClearAchievement(name string) bool {
	cname := append([]byte(name), 0)
	defer runtime.KeepAlive(cname)

	v, _, err := theDLL.call(flatAPI_ISteamUserStats_ClearAchievement, uintptr(s), uintptr(unsafe.Pointer(&cname[0])))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) StoreStats() bool {
	v, _, err := theDLL.call(flatAPI_ISteamUserStats_StoreStats, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUserStats) GetLeaderboardEntryCount(hSteamLeaderboard SteamLeaderboard_t) int32 {
	var v uintptr
	var err error
	
	if is32Bit {
		// On 32-bit, split the 64-bit hSteamLeaderboard parameter
		leaderboardLow := uintptr(uint32(hSteamLeaderboard))
		leaderboardHigh := uintptr(uint32(hSteamLeaderboard >> 32))
		v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetLeaderboardEntryCount, uintptr(s), leaderboardLow, leaderboardHigh)
	} else {
		v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetLeaderboardEntryCount, uintptr(s), uintptr(hSteamLeaderboard))
	}
	
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func (s steamUserStats) rawFindLeaderboard(name string) SteamAPICall_t {
	cName := append([]byte(name), 0)
	defer runtime.KeepAlive(cName)
	v1, v2, err := theDLL.call(flatAPI_ISteamUserStats_FindLeaderboard, uintptr(s), uintptr(unsafe.Pointer(&cName[0])))
	if err != nil {
		panic(err)
	}
	var result SteamAPICall_t
	if is32Bit {
		// On 32-bit: v1=EAX (low 32 bits), v2=EDX (high 32 bits)
		// Must convert to uint64 BEFORE shifting to avoid truncation
		result = SteamAPICall_t(uint64(v1) | (uint64(v2) << 32))
	} else {
		result = SteamAPICall_t(v1)
	}
	return result
}

func (s steamUserStats) rawGetDownloadedLeaderboardEntry(hSteamLeaderboardEntries SteamLeaderboardEntries_t, index int32) (success bool, entry LeaderboardEntry) {
	var rawEntry leaderboardEntry_t
	var v uintptr
	var err error
	
	if is32Bit {
		// On 32-bit, split the 64-bit hSteamLeaderboardEntries parameter
		entriesLow := uintptr(uint32(hSteamLeaderboardEntries))
		entriesHigh := uintptr(uint32(hSteamLeaderboardEntries >> 32))
		v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry, uintptr(s), entriesLow, entriesHigh, uintptr(index), uintptr(unsafe.Pointer(&rawEntry)), uintptr(0), uintptr(0))
	} else {
		v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry, uintptr(s), uintptr(hSteamLeaderboardEntries), uintptr(index), uintptr(unsafe.Pointer(&rawEntry)), uintptr(0), uintptr(0))
	}
	
	if err != nil {
		panic(err)
	}
	if byte(v) == 0 {
		return false, LeaderboardEntry{}
	}

	readEntry := rawEntry.Read()
	if readEntry.details > 0 {
		entry.details = make([]int32, readEntry.details)
		if is32Bit {
			entriesLow := uintptr(uint32(hSteamLeaderboardEntries))
			entriesHigh := uintptr(uint32(hSteamLeaderboardEntries >> 32))
			v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry, uintptr(s), entriesLow, entriesHigh, uintptr(index), uintptr(unsafe.Pointer(&rawEntry)), uintptr(unsafe.Pointer(&entry.details[0])), uintptr(readEntry.details))
		} else {
			v, _, err = theDLL.call(flatAPI_ISteamUserStats_GetDownloadedLeaderboardEntry, uintptr(s), uintptr(hSteamLeaderboardEntries), uintptr(index), uintptr(unsafe.Pointer(&rawEntry)), uintptr(unsafe.Pointer(&entry.details[0])), uintptr(readEntry.details))
		}
		if err != nil {
			panic(err)
		}
		if byte(v) == 0 {
			return false, LeaderboardEntry{}
		}
	}

	entry.globalRank = readEntry.globalRank
	entry.score = readEntry.score
	entry.steamIDUser = readEntry.steamIDUser
	entry.UGC = readEntry.UGC

	success = byte(v) != 0
	return
}

func (s steamUserStats) rawDownloadLeaderboardEntries(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardDataRequest ELeaderboardDataRequest, nRangeStart, nRangeEnd int32) SteamAPICall_t {
	var v1, v2 uintptr
	var err error
	
	if is32Bit {
		// On 32-bit, split the 64-bit hSteamLeaderboard parameter
		leaderboardLow := uintptr(uint32(hSteamLeaderboard))
		leaderboardHigh := uintptr(uint32(hSteamLeaderboard >> 32))
		v1, v2, err = theDLL.call(flatAPI_ISteamUserStats_DownloadLeaderboardEntries, uintptr(s), leaderboardLow, leaderboardHigh, uintptr(eLeaderboardDataRequest), uintptr(nRangeStart), uintptr(nRangeEnd))
	} else {
		v1, v2, err = theDLL.call(flatAPI_ISteamUserStats_DownloadLeaderboardEntries, uintptr(s), uintptr(hSteamLeaderboard), uintptr(eLeaderboardDataRequest), uintptr(nRangeStart), uintptr(nRangeEnd))
	}
	
	if err != nil {
		panic(err)
	}
	if is32Bit {
		return SteamAPICall_t(uint64(v1) | (uint64(v2) << 32))
	}
	return SteamAPICall_t(v1)
}

func (s steamUserStats) rawUploadLeaderboardScore(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardUploadScoreMethod ELeaderboardUploadScoreMethod, score int32, details []int32) SteamAPICall_t {
	var detailsPtr uintptr
	if len(details) > 0 {
		detailsPtr = uintptr(unsafe.Pointer(&details[0]))
	}
	
	var v1, v2 uintptr
	var err error
	
	if is32Bit {
		// On 32-bit, split the 64-bit hSteamLeaderboard parameter
		leaderboardLow := uintptr(uint32(hSteamLeaderboard))
		leaderboardHigh := uintptr(uint32(hSteamLeaderboard >> 32))
		v1, v2, err = theDLL.call(flatAPI_ISteamUserStats_UploadLeaderboardScore, uintptr(s), leaderboardLow, leaderboardHigh, uintptr(eLeaderboardUploadScoreMethod), uintptr(score), detailsPtr, uintptr(len(details)))
	} else {
		v1, v2, err = theDLL.call(flatAPI_ISteamUserStats_UploadLeaderboardScore, uintptr(s), uintptr(hSteamLeaderboard), uintptr(eLeaderboardUploadScoreMethod), uintptr(score), detailsPtr, uintptr(len(details)))
	}
	
	if err != nil {
		panic(err)
	}
	if is32Bit {
		return SteamAPICall_t(uint64(v1) | (uint64(v2) << 32))
	}
	return SteamAPICall_t(v1)
}

func SteamUtils() ISteamUtils {
	v, _, err := theDLL.call(flatAPI_SteamUtils)
	if err != nil {
		panic(err)
	}
	return steamUtils(v)
}

type steamUtils uintptr

func (s steamUtils) IsSteamRunningOnSteamDeck() bool {
	v, _, err := theDLL.call(flatAPI_ISteamUtils_IsSteamRunningOnSteamDeck, uintptr(s))
	if err != nil {
		panic(err)
	}

	return byte(v) != 0
}

func (s steamUtils) ShowFloatingGamepadTextInput(keyboardMode EFloatingGamepadTextInputMode, textFieldXPosition, textFieldYPosition, textFieldWidth, textFieldHeight int32) bool {
	v, _, err := theDLL.call(flatAPI_ISteamUtils_ShowFloatingGamepadTextInput, uintptr(s), uintptr(keyboardMode), uintptr(textFieldXPosition), uintptr(textFieldYPosition), uintptr(textFieldWidth), uintptr(textFieldHeight))
	if err != nil {
		panic(err)
	}
	return byte(v) != 0
}

// Basically a member function, but implemented as a standalone function because of generics limitations.
func steamUtilsGetAPICallResult[T any](s steamUtils, apiCall SteamAPICall_t, callbackType int) (result T, completed, success bool) {
	var failed uintptr
	var v uintptr
	var err error

	if is32Bit {
		// On 32-bit Windows, we need to manually split the 64-bit value
		// and pass it as two consecutive stack arguments
		apiCallLow := uintptr(uint32(apiCall))
		apiCallHigh := uintptr(uint32(apiCall >> 32))
		v, _, err = theDLL.call(flatAPI_ISteamUtils_GetAPICallResult,
			uintptr(s),
			apiCallLow,
			apiCallHigh,
			uintptr(unsafe.Pointer(&result)),
			uintptr(unsafe.Sizeof(result)),
			uintptr(callbackType),
			uintptr(unsafe.Pointer(&failed)))
	} else {
		// On 64-bit platforms, pass the 64-bit value directly
		v, _, err = theDLL.call(flatAPI_ISteamUtils_GetAPICallResult,
			uintptr(s),
			uintptr(apiCall),
			uintptr(unsafe.Pointer(&result)),
			uintptr(unsafe.Sizeof(result)),
			uintptr(callbackType),
			uintptr(unsafe.Pointer(&failed)))
	}

	if err != nil {
		panic(err)
	}
	completed = byte(v) != 0
	success = failed == 0
	return
}
