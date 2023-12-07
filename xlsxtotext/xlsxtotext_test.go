// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlsxtotext

import (
	"os"
	"testing"
)

var (
	xlsxPath = "../filesamples/file-sample_100kb.xlsx"
	xlsxURL  = "https://zzzx.snnu.edu.cn/__local/F/62/4E/896DC0778F426C757828CED677C_97EE9695_75E1.xlsx?e=.xlsx"
)

func TestOpen(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenSheet(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	texts, err := xp.ExtractSheetTexts(2)
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenReader(t *testing.T) {
	f, err := os.Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	xp, err := OpenReader(f, finfo.Size())
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenURL(t *testing.T) {
	xp, _, err := OpenURL(xlsxURL)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	t.Log(xp.NumSheets())

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseCharts(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	xp.SetParseCharts(true)

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseChartsNoFmt(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	xp.SetParseCharts(true)
	xp.SetDrawingsNoFmt(true)

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseDiagrams(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	xp.SetParseDiagrams(true)

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseImages(t *testing.T) {
	xp, err := Open(xlsxPath)
	if err != nil {
		t.Error(err)
	}
	defer xp.Close()

	xp.SetParseImages(true)

	texts, err := xp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}
