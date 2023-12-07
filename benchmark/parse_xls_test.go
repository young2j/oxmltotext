// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_ParseXlsToText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsToText()
	}
}

func Benchmark_ParseXlsByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsByTika()
	}
}
