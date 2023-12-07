// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"unsafe"
)

// StringTobytes converts a string to a byte slice.
//
// It takes a string parameter `s` and returns a byte slice.
//
// This function is implemented using the `unsafe` package to achieve zero cost conversion.
func StringTobytes(s string) []byte {
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	return b
}

// BytesToString converts a byte slice to a string.
//
// It takes a parameter b, which is a byte slice. It returns a string.
//
// This function is implemented using the `unsafe` package to achieve zero cost conversion.
func BytesToString(b []byte) string {
	s := ""
	if len(b) > 0 {
		s = unsafe.String(&b[0], len(b))
	}
	return s
}
