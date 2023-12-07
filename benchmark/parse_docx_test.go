// Copyright (c) 2023 young2j
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT



package main

import (
	"testing"

	"baliance.com/gooxml"
)

func init() {
	gooxml.DisableLogging()
}

func Benchmark_ParseDocxByGooxml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocxByGooxml()
	}
}
func Benchmark_ParseDocxByGodocx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocxByGodocx()
	}
}
func Benchmark_ParseDocxByDocconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocxByDocconv()
	}
}

func Benchmark_ParseDocxByOxmlToText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocxByOxmlToText()
	}
}

func Benchmark_ParseDocxByTika(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseDocxByTika()
	}
}
