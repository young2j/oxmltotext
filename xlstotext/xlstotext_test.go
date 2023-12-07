// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlstotext

import (
	"os"
	"testing"
)

var (
	xlsPath       = "../filesamples/file-sample_100kb.xls"
	xlsURL        = "https://zzrs.mastc.edu.cn/__local/D/84/4C/F211DEAE38FB979755B5F2DC38F_02B5B3E3_4200.xls?e=.xls"
	tikaServerURL = "http://localhost:9998/tika"
)

func TestExtractFromPath(t *testing.T) {
	res, err := ExtractFromPath(xlsPath)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURL(t *testing.T) {
	res, _, err := ExtractFromURL(xlsURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReader(t *testing.T) {
	f, err := os.Open(xlsPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	res, err := ExtractFromReader(f)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromPathByTika(t *testing.T) {
	res, _, err := ExtractFromPathByTika(xlsPath, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURLByTika(t *testing.T) {
	res, _, err := ExtractFromURLByTika(xlsURL, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReaderByTika(t *testing.T) {
	f, err := os.Open(xlsPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	res, _, err := ExtractFromReaderByTika(f, int(finfo.Size()), tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
