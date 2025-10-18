package steamworks

import (
	"slices"
)

var goCallbacks []func() bool

func RunCallbacks() {
	runCallbacksSteam()

	goCallbacks = slices.DeleteFunc(goCallbacks, func(f func() bool) bool {
		return f()
	})
}

// Return true when callback is finished and can be removed
func registerCallback(f func() bool) {
	goCallbacks = append(goCallbacks, f)
}
