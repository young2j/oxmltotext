// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_ParsePdfByLpdf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePdfByLpdf()
	}
}

func Benchmark_ParsePdfByDocconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePdfByDocconv()
	}
}

func Benchmark_ParsePdfByFitz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePdfByFitz()
	}
}

func Benchmark_ParsePdfByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePdfByFitz()
	}
}
