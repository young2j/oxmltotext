// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:build !ocr

/*
Package ocr provides OCR interface implementation for not using Tesseract OCR.
*/
package ocr

import (
	"io"

	"github.com/young2j/oxmltotext/types"
)

type defaultOcr struct {
}

func NewDefaultOcr() types.OCR {
	return &defaultOcr{}
}

func (d *defaultOcr) Run(r io.Reader) (string, error) {
	return "", nil
}

func (d *defaultOcr) Close() error {
	return nil
}
