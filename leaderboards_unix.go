//go:build !windows

package steamworks

// non-windows platforms (linux, apple, freebsd) use VALVE_CALLBACK_PACK_SMALL as defined in steamclientpublic.h
// All of these structs follow strict padding, without using cgo, to match VALVE_CALLBACK_PACK_SMALL
// This means every struct member, if smaller than 4 bytes, is padded to 4 bytes

import "encoding/binary"

type leaderboardFindResult_t struct {
	// m_hSteamLeaderboard: 8 bytes
	// m_bLeaderboardFound: 1 byte (4 bytes padded)
	// total: 12 bytes
	data [12]byte
}

func (me leaderboardFindResult_t) Read() leaderboardFindResult {
	var result leaderboardFindResult
	result.steamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[0:8]))
	result.leaderboardFound = me.data[8] != 0
	return result
}

type leaderboardScoresDownloaded_t struct {
	// m_hSteamLeaderboard:        8 bytes
	// m_hSteamLeaderboardEntries: 8 bytes
	// m_cEntryCount:              4 bytes
	// total: 20 bytes
	data [20]byte
}

func (me leaderboardScoresDownloaded_t) Read() leaderboardScoresDownloaded {
	var result leaderboardScoresDownloaded
	result.steamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[0:8]))
	result.steamLeaderboardEntries = SteamLeaderboardEntries_t(binary.NativeEndian.Uint64(me.data[8:16]))
	result.entryCount = int32(binary.NativeEndian.Uint32(me.data[16:20]))
	return result
}

type leaderboardScoreUploaded_t struct {
	// m_bSuccess:            1 byte (4 bytes padded)
	// m_hSteamLeaderboard:   8 bytes
	// m_nScore:              4 bytes
	// m_bScoreChanged:       1 byte (4 bytes padded)
	// m_nGlobalRankNew:      4 bytes
	// m_nGlobalRankPrevious: 4 bytes
	// total: 28 bytes
	data [28]byte
}

func (me leaderboardScoreUploaded_t) Read() leaderboardScoreUploaded {
	var result leaderboardScoreUploaded
	result.success = me.data[0] != 0
	result.steamLeaderboard = SteamLeaderboard_t(binary.NativeEndian.Uint64(me.data[4:12]))
	result.score = int32(binary.NativeEndian.Uint32(me.data[12:16]))
	result.scoreChanged = me.data[16] != 0
	result.globalRankNew = int32(binary.NativeEndian.Uint32(me.data[20:24]))
	result.globalRankPrevious = int32(binary.NativeEndian.Uint32(me.data[24:28]))
	return result
}

type leaderboardEntry_t struct {
	// m_steamIDUser: 8 bytes
	// m_nGlobalRank: 4 bytes
	// m_nScore:      4 bytes
	// m_cDetails:    4 bytes
	// m_hUGC:        8 bytes
	// total:         28 bytes
	data [28]byte
}

func (me leaderboardEntry_t) Read() leaderboardEntry {
	var result leaderboardEntry
	result.steamIDUser = CSteamID(binary.NativeEndian.Uint64(me.data[0:8]))
	result.globalRank = int32(binary.NativeEndian.Uint32(me.data[8:12]))
	result.score = int32(binary.NativeEndian.Uint32(me.data[12:16]))
	result.details = int32(binary.NativeEndian.Uint32(me.data[16:20]))
	result.UGC = UGCHandle_t(binary.NativeEndian.Uint64(me.data[20:28]))
	return result
}