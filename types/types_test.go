// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:build ocr

package types

import (
	"os"
	"testing"
)

func Test_defaultOcr_Run(t *testing.T) {
	ocr := NewDefaultOcr()
	defer ocr.Close()
	f, err := os.Open("../filesamples/img-sample_idcard.jpeg")
	if err != nil {
		t.Log(err)
	}
	defer f.Close()

	text, err := ocr.Run(f)
	if err != nil {
		t.Log(err)
	}

	t.Log(text)
}
