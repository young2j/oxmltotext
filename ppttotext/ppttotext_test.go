// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ppttotext

import (
	"os"
	"testing"
)

var (
	pptPath       = "../filesamples/file-sample_500kb.ppt"
	pptURL        = "https://zhangbaohui.snnu.edu.cn/__local/4/11/32/53981EBE55E513C149441E6174F_C64C30B4_C3A00.ppt?e=.ppt"
	tikaServerURL = "http://localhost:9998/tika"
)

func TestExtractFromPathByTika(t *testing.T) {
	res, _, err := ExtractFromPathByTika(pptPath, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURLByTika(t *testing.T) {
	res, _, err := ExtractFromURLByTika(pptURL, tikaServerURL)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReaderByTika(t *testing.T) {
	f, err := os.Open(pptPath)
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
