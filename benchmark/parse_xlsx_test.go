// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "testing"

func Benchmark_ParseXlsxByGooxml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByGooxml()
	}
}

func Benchmark_ParseXlsxByTealeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByTealeg()
	}
}

func Benchmark_ParseXlsxByExcelize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByExcelize()
	}
}

func Benchmark_ParseXlsxByOxmlToText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByOxmlToText()
	}
}

func Benchmark_ParseXlsxByCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByCmd()
	}
}

func Benchmark_ParseXlsxByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseXlsxByTika()
	}
}
