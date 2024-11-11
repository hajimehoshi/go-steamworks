package steamworks

import (
	"sync"
	"time"
)

var interval time.Duration = 20 * time.Millisecond
var timeout time.Duration = 10 * time.Second

type iCallbackExpected int

const (
	iCallbackExpected_LeaderboardFindResult_t       iCallbackExpected = 1104
	iCallbackExpected_LeaderboardScoresDownloaded_t iCallbackExpected = 1105
	iCallbackExpected_LeaderboardScoreUploaded_t    iCallbackExpected = 1106
)

type callbackClient struct {
	timeout          time.Duration
	interval         time.Duration
	intervalTimer    *time.Timer
	steamUtils       ISteamUtils
	callbackArgsWait []*CallbackArgs
	callbackArgsChan chan *CallbackArgs
	closeSignal      chan bool
}
type CallbackArgs struct {
	CallbackAPI      SteamAPICall_t
	CallbackExpected iCallbackExpected
	CallbaseSize     int
	SuccessFunc      callbackSuccessFunc
	TimeoutFunc      callbackTimeoutFunc

	beginTime time.Time
}

type callbackSuccessFunc func(ret []byte)
type callbackTimeoutFunc func(callbackTime time.Time, callbackSpend time.Duration)

var callbackCli *callbackClient
var callbackOnce sync.Once

func defaultCallbackCli() *callbackClient {
	callbackOnce.Do(func() {
		callbackCli = &callbackClient{
			timeout:          timeout,
			interval:         interval,
			intervalTimer:    time.NewTimer(interval),
			steamUtils:       SteamUtils(),
			callbackArgsWait: make([]*CallbackArgs, 0),
			callbackArgsChan: make(chan *CallbackArgs, 10),
			closeSignal:      make(chan bool, 1),
		}
	})
	go callbackCli.run()
	return callbackCli
}

func (c *callbackClient) setCallback(callbackArgs *CallbackArgs) {
	callbackArgs.beginTime = time.Now()
	c.callbackArgsChan <- callbackArgs
}

func (c *callbackClient) run() {
	defer close(c.callbackArgsChan)
	defer close(c.closeSignal)
	for {
		<-c.intervalTimer.C
		select {
		case arg := <-c.callbackArgsChan:
			c.callbackArgsWait = append(c.callbackArgsWait, arg)
		default:
			if len(c.callbackArgsWait) == 0 {
				c.intervalTimer.Stop()
				if !c.waitArgs() {
					return
				}
			}
			RunCallbacks()
			for i := 0; i < len(c.callbackArgsWait); {
				a := c.callbackArgsWait[i]
				spend := time.Since(a.beginTime)
				if spend > c.timeout {
					a.TimeoutFunc(a.beginTime, spend)
					c.callbackArgsWait = append(c.callbackArgsWait[:i], c.callbackArgsWait[i+1:]...)
					continue
				}
				call, ok, fail := c.steamUtils.GetAPICallResult(a.CallbackAPI, a.CallbackExpected, a.CallbaseSize)
				if ok && !fail {
					a.SuccessFunc(call)
					c.callbackArgsWait = append(c.callbackArgsWait[:i], c.callbackArgsWait[i+1:]...)
					continue
				}
				i++
			}
		}
		c.intervalTimer.Reset(c.interval)
	}

}

func (c *callbackClient) waitArgs() bool {
	select {
	case arg := <-c.callbackArgsChan:
		c.callbackArgsWait = append(c.callbackArgsWait, arg)
		return true
	case <-c.closeSignal:
		return false
	}
}

func (c *callbackClient) close() {
	c.closeSignal <- true
}
