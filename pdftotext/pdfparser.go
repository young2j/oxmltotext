// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pdftotext

import (
	"strings"

	"github.com/gen2brain/go-fitz"
)

// PdfParser is a wrapper around the go-fitz library.
type PdfParser struct {
	pdf     *fitz.Document
	pageSep string
}

func newPdfParser() *PdfParser {
	return &PdfParser{
		pageSep: strings.Repeat("-", 100) + "\n",
	}
}
