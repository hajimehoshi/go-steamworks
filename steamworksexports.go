package steamworks

import "C"
import "unsafe"

// This is exactly the go hook that was passed in by the app to receive steam warning messages
var userWarningMessageHook func(severity int, debugText string)

//export warningMessageGoHook
func warningMessageGoHook(severity C.int, debugText unsafe.Pointer) {
	snappedHook := userWarningMessageHook
	if snappedHook != nil {
		snappedHook(int(severity), cStringToGoString(uintptr(debugText), 0))
	}
}
