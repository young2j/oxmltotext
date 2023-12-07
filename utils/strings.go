// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import "bytes"

// MaxLineLen returns the maximum line length in a given string.
//
// Parameters:
//   - s: the input string to check for maximum line length.
//
// Return type:
//   - int: the maximum line length in the given string.
func MaxLineLen(s string) int {
	maxLen := 0
	lineLen := 0
	for _, b := range StringTobytes(s) {
		if b != '\n' {
			lineLen++
		} else {
			if lineLen > maxLen {
				maxLen = lineLen
			}
			lineLen = 0
		}
	}
	if lineLen > maxLen {
		maxLen = lineLen
	}

	return maxLen
}

// MaxLineLenWithPrefix calculates the maximum line length in a string with a given prefix.
//
// Parameters:
//   - s: the input string
//   - prefix: the prefix to add to each line
//
// Returns:
//   - string: the modified string with the added prefix
//   - int: the maximum line length (including the prefix)
func MaxLineLenWithPrefix(s string, prefix []byte) (string, int) {
	maxLen := 0
	lineLen := 0
	newS := make([]byte, 0, len(s)+10)
	buf := bytes.NewBuffer(newS)
	for i, b := range StringTobytes(s) {
		if i == 0 {
			buf.Write(prefix)
		}

		if b != '\n' {
			lineLen++
			buf.WriteByte(b)
		} else {
			if lineLen > maxLen {
				maxLen = lineLen
			}
			lineLen = 0
			buf.WriteByte(b)
			buf.Write(prefix)
		}
	}

	if lineLen > maxLen {
		maxLen = lineLen
	}

	maxLen += len(prefix)

	return buf.String(), maxLen
}
