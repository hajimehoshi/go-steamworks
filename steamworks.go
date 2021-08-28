// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

//go:generate go run gen.go

package steamworks

type ISteamApps interface {
	GetCurrentGameLanguage() string
}

const (
	flatAPI_RestartAppIfNecessary             = "SteamAPI_RestartAppIfNecessary"
	flatAPI_Init                              = "SteamAPI_Init"
	flatAPI_SteamApps                         = "SteamAPI_SteamApps_v008"
	flatAPI_ISteamApps_GetCurrentGameLanguage = "SteamAPI_ISteamApps_GetCurrentGameLanguage"
)
