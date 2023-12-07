// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_ParsePptxByGooxml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptxByGooxml()
	}
}

func Benchmark_ParsePptxByDocconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptxByDocconv()
	}
}

func Benchmark_ParsePptxByOxmlToText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptxByOxmlToText()
	}
}

func Benchmark_ParsePptxByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePptxByTika()
	}
}
