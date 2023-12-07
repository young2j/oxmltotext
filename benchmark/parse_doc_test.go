// Copyright (c) 2023 young2j
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT


package main

import "testing"

func Benchmark_ParseDocByAntiword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocByAntiword()
	}
}

func Benchmark_ParseDocByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocByTika()
	}
}
