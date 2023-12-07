// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"fmt"
	"strings"
)

var strs = []string{
	"This", "is", "a", "long", "string", "that", "we", "want", "to", "split", "into", "a", "slice", "of", "strings", "based", "on", "a", "delimiter", "\n",
	"This", "is", "second", "long", "string", "that", "we", "want", "to", "split", "into", "a", "slice", "of", "strings", "based", "on", "a", "delimiter",
}

func JoinByAdd() string {
	res := ""
	for _, s := range strs {
		res += s + " "
	}

	return res
}

func JoinByStringsJoin() string {
	res := strings.Join(strs, " ")

	return res
}

func JoinBySprintf() string {
	res := "%s"
	for _, s := range strs {
		res = fmt.Sprintf(res, s+" %s")
	}

	return res
}

func JoinByBytesBuffer() string {
	buf := new(bytes.Buffer)
	for _, s := range strs {
		buf.WriteString(s)
		buf.WriteString(" ")
	}
	return buf.String()
}

func JoinByStringsBuilder() string {
	b := new(strings.Builder)

	for _, s := range strs {
		b.WriteString(s)
		b.WriteString(" ")
	}

	return b.String()
}
