//go:build windows

package steamworks

// 64-bit Windows uses VALVE_CALLBACK_PACK_LARGE as defined in steamclientpublic.h (8-byte alignment)

import (
	"encoding/binary"
)

type LeaderboardFindResult_t struct {
	// m_hSteamLeaderboard: 8 bytes
	// m_bLeaderboardFound: 1 byte (8 bytes padded)
	// total: 16 bytes
	data [16]byte
}

func (me LeaderboardFindResult_t) Read() LeaderboardFindResult {
	var result LeaderboardFindResult
	result.SteamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[0:8]))
	result.LeaderboardFound = me.data[8] != 0
	return result
}

type LeaderboardScoresDownloaded_t struct {
	// m_hSteamLeaderboard:        8 bytes
	// m_hSteamLeaderboardEntries: 8 bytes
	// m_cEntryCount:              4 bytes (8 bytes padded)
	// total: 24 bytes
	data [24]byte
}

func (me LeaderboardScoresDownloaded_t) Read() LeaderboardScoresDownloaded {
	var result LeaderboardScoresDownloaded
	result.SteamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[0:8]))
	result.SteamLeaderboardEntries = SteamLeaderboardEntries_t(binary.NativeEndian.Uint64(me.data[8:16]))
	result.EntryCount = int32(binary.NativeEndian.Uint32(me.data[16:20]))
	return result
}

type LeaderboardScoreUploaded_t struct {
	// m_bSuccess:            1 byte  (4 bytes padded)
	// m_hSteamLeaderboard:   8 bytes (12 bytes padded)
	// m_nScore:              4 bytes (8 bytes padded)
	// m_bScoreChanged:       1 byte  (4 bytes padded)
	// m_nGlobalRankNew:      4 bytes (4 bytes padded)
	// m_nGlobalRankPrevious: 4 bytes (4 bytes padded)
	// total: 48 bytes
	data [48]byte
}

func (me LeaderboardScoreUploaded_t) Read() LeaderboardScoreUploaded {
	var result LeaderboardScoreUploaded
	result.Success = me.data[0] != 0
	result.SteamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[4:12]))
	result.Score = int32(binary.NativeEndian.Uint32(me.data[16:20]))
	result.ScoreChanged = me.data[20] != 0
	result.GlobalRankNew = int32(binary.NativeEndian.Uint32(me.data[24:28]))
	result.GlobalRankPrevious = int32(binary.NativeEndian.Uint32(me.data[28:32]))
	return result
}

type LeaderboardEntry_t struct {
	// m_steamIDUser: 8 bytes
	// m_nGlobalRank: 4 bytes
	// m_nScore:      4 bytes (8 bytes padded)
	// m_cDetails:    4 bytes
	// m_hUGC:        8 bytes
	// total:         40 bytes
	data [40]byte
}

func (me LeaderboardEntry_t) Read() LeaderboardEntry {
	var result LeaderboardEntry
	result.SteamIDUser = CSteamID(binary.NativeEndian.Uint64(me.data[0:8]))
	result.GlobalRank = int32(binary.NativeEndian.Uint32(me.data[8:12]))
	result.Score = int32(binary.NativeEndian.Uint32(me.data[12:16]))
	result.Details = int32(binary.NativeEndian.Uint32(me.data[16:20]))
	result.UGC = UGCHandle_t(binary.NativeEndian.Uint64(me.data[24:32]))
	return result
}
