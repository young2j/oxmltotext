// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package doctotext

import (
	"os"
	"testing"
)

var (
	docPath       = "../filesamples/file-sample_100kb.doc"
	docURL        = "https://zzrs.mastc.edu.cn/__local/6/3F/8E/562014FE1D4808E809ACA2A989F_920F1792_E800.doc?e=.doc"
	tikaServerURL = "http://localhost:9998/tika"
)

func TestExtractFromPath(t *testing.T) {
	res, err := ExtractFromPath(docPath)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURL(t *testing.T) {
	res, _, err := ExtractFromURL(docURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReader(t *testing.T) {
	f, err := os.Open(docPath)
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
	res, _, err := ExtractFromPathByTika(docPath, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURLByTika(t *testing.T) {
	res, _, err := ExtractFromURLByTika(docURL, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReaderByTika(t *testing.T) {
	f, err := os.Open(docPath)
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
