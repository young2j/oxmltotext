// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pptxtotext

import (
	"os"
	"testing"
)

var (
	pptxPath = "../filesamples/file-sample_500kb.pptx"
	pptxURL  = "https://zcc.czu.cn/_upload/article/files/06/b5/8a64cb854694bcd2265ad0b96c99/65ba8668-56f7-4cd6-ab51-b55349964a17.pptx"
)

func TestOpen(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenSlide(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	texts, err := pp.ExtractSlideTexts(1)
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenReader(t *testing.T) {
	f, err := os.Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	pp, err := OpenReader(f, finfo.Size())
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenURL(t *testing.T) {
	pp, _, err := OpenURL(pptxURL)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	t.Log(pp.NumSlides())

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseCharts(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	pp.SetParseCharts(true)

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseChartsNoFmt(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	pp.SetParseCharts(true)
	pp.SetDrawingsNoFmt(true)

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}
func TestParseDiagrams(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	pp.SetParseDiagrams(true)

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}
func TestParseImages(t *testing.T) {
	pp, err := Open(pptxPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	pp.SetParseImages(true)

	texts, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}
