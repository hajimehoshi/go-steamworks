package steamworks

import "fmt"

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

type SteamLeaderboard_t uint64
type SteamLeaderboardEntries_t uint64

type LeaderboardFindResult_t struct {
	hSteamLeaderboard SteamLeaderboard_t
	bLeaderboardFound uint8
}

type LeaderboardScoresDownloaded_t struct {
	steamLeaderboard        SteamLeaderboard_t
	steamLeaderboardEntries SteamLeaderboardEntries_t
	entryCount              int32
}

type leaderboardScoreUploaded_t struct {
	bSuccess            uint8
	hSteamLeaderboard   SteamLeaderboard_t
	nScore              int32
	bScoreChanged       uint8
	nGlobalRankNew      int32
	nGlobalRankPrevious int32
}

type UGCHandle_t uint64

// Raw steam interface struct
type leaderboardEntry_t struct {
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
	bScoreChanged       uint8
	nGlobalRankNew      int32
	nGlobalRankPrevious int32
}

func (s steamUserStats) FindLeaderboard(name string, onComplete func(handle SteamLeaderboard_t, found bool, err error)) {
	handle := s.rawFindLeaderboard(name)
	registerCallback(func() bool {
		result, completed, success := steamUtilsGetAPICallResult[LeaderboardFindResult_t](SteamUtils().(steamUtils), handle, LeaderboardFindResult_k_iCallback)
		if !completed {
			return false
		}

		if success {
			found := result.bLeaderboardFound != 0
			onComplete(result.hSteamLeaderboard, found, nil)
		} else {
			onComplete(0, false, fmt.Errorf("failed to find leaderboard %s", name))
		}
		return true
	})
}

func (s steamUserStats) DownloadLeaderboardEntries(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardDataRequest ELeaderboardDataRequest, nRangeStart, nRangeEnd int32, onComplete func(entries []LeaderboardEntry, err error)) {
	v := s.rawDownloadLeaderboardEntries(hSteamLeaderboard, eLeaderboardDataRequest, nRangeStart, nRangeEnd)

	handle := SteamAPICall_t(v)
	registerCallback(func() bool {
		result, completed, success := steamUtilsGetAPICallResult[LeaderboardScoresDownloaded_t](SteamUtils().(steamUtils), handle, LeaderboardScoresDownloaded_k_iCallback)
		if !completed {
			return false
		}

		if success {
			if result.entryCount == 0 {
				onComplete(nil, nil)
				return true
			}
			entries := make([]LeaderboardEntry, result.entryCount)

			// Now grab all the entries with the detail count we learned
			for i := range result.entryCount {
				var ok bool
				ok, entries[i] = s.rawGetDownloadedLeaderboardEntry(result.steamLeaderboardEntries, i)
				if !ok {
					onComplete(nil, fmt.Errorf("failed to get leaderboard entry %d", i))
					return true
				}
			}

			onComplete(entries, nil)
		} else {
			onComplete(nil, fmt.Errorf("failed to download leaderboard entries"))
		}
		return true
	})
}

func (s steamUserStats) UploadLeaderboardScore(hSteamLeaderboard SteamLeaderboard_t, eLeaderboardUploadScoreMethod ELeaderboardUploadScoreMethod, score int32, details []int32, onComplete func(result LeaderboardScoreUploaded, err error)) {
	v := s.rawUploadLeaderboardScore(hSteamLeaderboard, eLeaderboardUploadScoreMethod, score, details)

	handle := SteamAPICall_t(v)
	registerCallback(func() bool {
		rawResult, completed, success := steamUtilsGetAPICallResult[leaderboardScoreUploaded_t](SteamUtils().(steamUtils), handle, LeaderboardScoreUploaded_k_iCallback)
		if !completed {
			return false
		}

		if !success {
			onComplete(LeaderboardScoreUploaded{}, fmt.Errorf("GetAPICallResult failed"))
		} else if rawResult.bSuccess == 0 {
			onComplete(LeaderboardScoreUploaded{}, fmt.Errorf("bSuccess is false"))
		} else {
			result := LeaderboardScoreUploaded{
				nScore:              rawResult.nScore,
				bScoreChanged:       rawResult.bScoreChanged,
				nGlobalRankNew:      rawResult.nGlobalRankNew,
				nGlobalRankPrevious: rawResult.nGlobalRankPrevious,
			}
			onComplete(result, nil)
		}
		return true
	})
}
