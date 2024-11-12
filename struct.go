package steamworks

import (
	"reflect"
	"unsafe"
)

/*
typedef unsigned long long int SteamLeaderboard_t;
typedef unsigned long long int SteamLeaderboardEntries_t;
typedef unsigned char uint8;

typedef struct {
	unsigned long int m_steamIDUser;
	int m_nGlobalRank;
	int m_nScore;
	int m_cDetails;
	unsigned long int m_hUGC;
} LeaderboardEntry_t;

typedef struct {
	SteamLeaderboard_t m_hSteamLeaderboard;
	uint8 m_bLeaderboardFound;
}LeaderboardFindResult_t;

typedef struct{
	SteamLeaderboard_t m_hSteamLeaderboard;
	SteamLeaderboardEntries_t m_hSteamLeaderboardEntries;
	int m_cEntryCount;
}LeaderboardScoresDownloaded_t;

typedef struct{
	uint8 m_bSuccess;
	SteamLeaderboard_t m_hSteamLeaderboard;
	int m_nScore;
	uint8 m_bScoreChanged;
	int m_nGlobalRankNew;
	int m_nGlobalRankPrevious;
}LeaderboardScoreUploaded_t;
*/
import "C"

type IStruct interface {
	Size() uintptr
	CStructPtr() uintptr
}

type LeaderboardScoreUploaded_t struct {
	Success            bool
	SteamLeaderboard   SteamLeaderboard_t
	Score              int
	ScoreChanged       bool
	GlobalRankNew      int
	GlobalRankPrevious int
}

func (l LeaderboardScoreUploaded_t) FromCStruct(cstruct C.LeaderboardScoreUploaded_t) LeaderboardScoreUploaded_t {
	return LeaderboardScoreUploaded_t{
		Success:            cstruct.m_bSuccess != 0,
		SteamLeaderboard:   SteamLeaderboard_t(cstruct.m_hSteamLeaderboard),
		Score:              int(cstruct.m_nScore),
		ScoreChanged:       cstruct.m_bScoreChanged != 0,
		GlobalRankNew:      int(cstruct.m_nGlobalRankNew),
		GlobalRankPrevious: int(cstruct.m_nGlobalRankPrevious),
	}
}

func (l LeaderboardScoreUploaded_t) FromByte(b []byte) LeaderboardScoreUploaded_t {
	return l.FromCStruct(**(**C.LeaderboardScoreUploaded_t)(unsafe.Pointer(&b)))
}

func (l LeaderboardScoreUploaded_t) CStruct() C.LeaderboardScoreUploaded_t {
	return C.LeaderboardScoreUploaded_t{}
}

func (l LeaderboardScoreUploaded_t) Size() uintptr {
	return reflect.TypeOf(l.CStruct()).Size()
}

type LeaderboardScoresDownloaded_t struct {
	SteamLeaderboard        SteamLeaderboard_t
	SteamLeaderboardEntries SteamLeaderboardEntries_t
	EntryCount              int
}

func (l LeaderboardScoresDownloaded_t) FromCStruct(cstruct C.LeaderboardScoresDownloaded_t) LeaderboardScoresDownloaded_t {
	return LeaderboardScoresDownloaded_t{
		SteamLeaderboard:        SteamLeaderboard_t(cstruct.m_hSteamLeaderboard),
		SteamLeaderboardEntries: SteamLeaderboardEntries_t(cstruct.m_hSteamLeaderboardEntries),
		EntryCount:              int(cstruct.m_cEntryCount),
	}
}

func (l LeaderboardScoresDownloaded_t) FromByte(b []byte) LeaderboardScoresDownloaded_t {
	return l.FromCStruct(**(**C.LeaderboardScoresDownloaded_t)(unsafe.Pointer(&b)))
}

func (l LeaderboardScoresDownloaded_t) CStruct() C.LeaderboardScoresDownloaded_t {
	return C.LeaderboardScoresDownloaded_t{}
}

func (l LeaderboardScoresDownloaded_t) Size() uintptr {
	return reflect.TypeOf(l.CStruct()).Size()
}

type LeaderboardFindResult_t struct {
	SteamLeaderboard SteamLeaderboard_t
	LeaderboardFound bool
}

func (l LeaderboardFindResult_t) FromByte(b []byte) LeaderboardFindResult_t {
	return l.FromCStruct(**(**C.LeaderboardFindResult_t)(unsafe.Pointer(&b)))
}

func (l LeaderboardFindResult_t) FromCStruct(cstruct C.LeaderboardFindResult_t) LeaderboardFindResult_t {
	return LeaderboardFindResult_t{
		SteamLeaderboard: SteamLeaderboard_t(cstruct.m_hSteamLeaderboard),
		LeaderboardFound: cstruct.m_bLeaderboardFound != 0,
	}
}

func (l LeaderboardFindResult_t) CStruct() C.LeaderboardFindResult_t {
	return C.LeaderboardFindResult_t{}
}

func (l LeaderboardFindResult_t) Size() uintptr {
	return reflect.TypeOf(l.CStruct()).Size()
}

type LeaderboardEntry_t struct {
	SteamIDUser CSteamID
	GlobalRank  int
	Score       int
	Details     int
	UGC         UGCHandle_t
}

func (l LeaderboardEntry_t) FromCStruct(cstruct C.LeaderboardEntry_t) LeaderboardEntry_t {
	return LeaderboardEntry_t{
		SteamIDUser: CSteamID(cstruct.m_steamIDUser),
		GlobalRank:  int(cstruct.m_nGlobalRank),
		Score:       int(cstruct.m_nScore),
		Details:     int(cstruct.m_cDetails),
		UGC:         UGCHandle_t(cstruct.m_hUGC),
	}
}

func (l LeaderboardEntry_t) CStruct() C.LeaderboardEntry_t {
	return C.LeaderboardEntry_t{}
}

func (l LeaderboardEntry_t) Size() uintptr {
	return reflect.TypeOf(l.CStruct()).Size()
}
