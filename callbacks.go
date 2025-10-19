package steamworks

import (
	"slices"
)

var goCallbacks []func() bool
var goCallbacksToAppend []func() bool

func RunCallbacks() {
	runCallbacksSteam()

	goCallbacks = append(goCallbacks, goCallbacksToAppend...)
	goCallbacksToAppend = nil

	goCallbacks = slices.DeleteFunc(goCallbacks, func(f func() bool) bool {
		return f()
	})
}

// Return true when callback is finished and can be removed
func registerCallback(f func() bool) {
	if f == nil {
		panic("cannot register nil callback")
	}
	// If we were to add directly to goCallbacks, then it would cause bugs to register a callback from a callback itself.
	// So, append to a separate slice, and merge in at the start of the next RunCallbacks call.
	goCallbacksToAppend = append(goCallbacksToAppend, f)
}
