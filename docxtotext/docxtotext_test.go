// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package docxtotext

import (
	"bytes"
	"image/jpeg"
	"os"
	"testing"
)

var (
	docxPath = "../filesamples/file-sample_100kb.docx"
	docxURL  = "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.docx"
)

func TestOpen(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenReader(t *testing.T) {
	f, err := os.Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	dp, err := OpenReader(f, finfo.Size())
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestOpenURL(t *testing.T) {
	dp, _, err := OpenURL(docxURL)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestExtractImages(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	imgs, err := dp.ExtractImages()
	if err != nil {
		t.Error(err)
	}

	for _, img := range imgs {
		//-- save to file --
		// fname := filepath.Base(img.Name)
		// f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	t.Log(err)
		// }
		// defer f.Close()
		// -- memory buffer --
		f := new(bytes.Buffer)

		err = jpeg.Encode(f, img.Raw, nil)
		if err != nil {
			t.Log(err)
		}
		t.Logf("img format: %v size:%v bytes", img.Format, len(f.Bytes()))
	}
}

func TestParseCharts(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	dp.SetParseCharts(true)

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseChartsNoFmt(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	dp.SetParseCharts(true)
	dp.SetDrawingsNoFmt(true)

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseDiagrams(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	dp.SetParseDiagrams(true)

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}

func TestParseImages(t *testing.T) {
	dp, err := Open(docxPath)
	if err != nil {
		t.Error(err)
	}
	defer dp.Close()

	dp.SetParseImages(true)

	texts, err := dp.ExtractTexts()
	if err != nil {
		t.Error(err)
	}

	t.Log(texts)
}
