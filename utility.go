// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The go-steamworks Authors

package steamworks

import (
	"runtime"
	"unsafe"
)

func cStringToGoString(v uintptr, sizeHint int) string {
	bs := make([]byte, 0, sizeHint)
	for i := int32(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(v))
		v += unsafe.Sizeof(byte(0))
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs)
}

// Helper function to convert Go string to C string
func goStringToCString(s string) (uintptr, func()) {
	bytes := append([]byte(s), 0)
	ptr := uintptr(unsafe.Pointer(&bytes[0]))
	return ptr, func() { runtime.KeepAlive(bytes) }
}

