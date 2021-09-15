// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

type ISteamApps interface {
	GetAppInstallDir(appID AppId_t) string
	GetCurrentGameLanguage() string
}

type ISteamRemoteStorage interface {
	FileWrite(file string, data []byte) bool
	FileRead(file string, data []byte) int32
	FileDelete(file string) bool
	GetFileSize(file string) int32
}

type ISteamUserStats interface {
	GetAchievement(name string) (achieved, success bool)
	SetAchievement(name string) bool
	ClearAchievement(name string) bool
	StoreStats() bool
}

const (
	flatAPI_RestartAppIfNecessary = "SteamAPI_RestartAppIfNecessary"
	flatAPI_Init                  = "SteamAPI_Init"

	flatAPI_SteamApps                         = "SteamAPI_SteamApps_v008"
	flatAPI_ISteamApps_GetAppInstallDir       = "SteamAPI_ISteamApps_GetAppInstallDir"
	flatAPI_ISteamApps_GetCurrentGameLanguage = "SteamAPI_ISteamApps_GetCurrentGameLanguage"

	flatAPI_SteamRemoteStorage              = "SteamAPI_SteamRemoteStorage_v014"
	flatAPI_ISteamRemoteStorage_FileWrite   = "SteamAPI_ISteamRemoteStorage_FileWrite"
	flatAPI_ISteamRemoteStorage_FileRead    = "SteamAPI_ISteamRemoteStorage_FileRead"
	flatAPI_ISteamRemoteStorage_FileDelete  = "SteamAPI_ISteamRemoteStorage_FileDelete"
	flatAPI_ISteamRemoteStorage_GetFileSize = "SteamAPI_ISteamRemoteStorage_GetFileSize"

	flatAPI_SteamUserStats                   = "SteamAPI_SteamUserStats_v012"
	flatAPI_ISteamUserStats_GetAchievement   = "SteamAPI_ISteamUserStats_GetAchievement"
	flatAPI_ISteamUserStats_SetAchievement   = "SteamAPI_ISteamUserStats_SetAchievement"
	flatAPI_ISteamUserStats_ClearAchievement = "SteamAPI_ISteamUserStats_ClearAchievement"
	flatAPI_ISteamUserStats_StoreStats       = "SteamAPI_ISteamUserStats_StoreStats"
)
