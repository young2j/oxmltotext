// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pdftotext

import (
	"os"
	"testing"
)

var (
	pdfPath = "../filesamples/file-sample_500kb.pdf"
	// pdfURL  = "https://zzzx.snnu.edu.cn/__local/4/08/69/EAAC563EBFBD5EA895575AB4DB1_CEFEE2BD_10800.pdf?e=.pdf"
	// pdfURL = "https://zhangbaohui.snnu.edu.cn/info/1157/DevelopmentReciprocalTeach-ZHANGBH.pdf"
	pdfURL = "http://zhangbaohui.snnu.edu.cn/icse2012/files/2012_Nanjing_University_International_Education_Conference_Chinese_Proceedings.pdf"
)

func TestExtractFromPath(t *testing.T) {
	pp, err := Open(pdfPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	res, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractPage(t *testing.T) {
	pp, err := Open(pdfPath)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	res, err := pp.ExtractPageTexts(1, 2)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestExtractFromURL(t *testing.T) {
	pp, _, err := OpenURL(pdfURL)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	res, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
func TestExtractFromReader(t *testing.T) {
	f, err := os.Open(pdfPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	pp, err := OpenReader(f)
	if err != nil {
		t.Error(err)
	}
	defer pp.Close()

	res, err := pp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
