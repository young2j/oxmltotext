// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_ParsePptByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptByTika()
	}
}

func Benchmark_ParsePptByTikaCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptByTikaCmd()
	}
}

func Benchmark_ParsePptByUnoconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptByUnoconv()
	}
}
