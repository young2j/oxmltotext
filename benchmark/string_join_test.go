// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_JoinByAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinByAdd()
	}
}
func Benchmark_JoinByStringsJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinByStringsJoin()
	}
}
func Benchmark_JoinBySprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinBySprintf()
	}
}
func Benchmark_JoinByBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinByBytesBuffer()
	}
}
func Benchmark_JoinByStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinByStringsBuilder()
	}
}
