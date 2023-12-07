// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pdftotext

import (
	"strings"
)

// SetPageSep sets the page text separator for the PdfParser. Default is "-"x100.
func (pp *PdfParser) SetPageSep(sep string) {
	pp.pageSep = sep
}

// NumPages returns the number of pages in the PDF.
func (pp *PdfParser) NumPages() int {
	return pp.pdf.NumPage()
}

// Close closes the opened pdf document of PdfParser.
func (pp *PdfParser) Close() error {
	return pp.pdf.Close()
}

// ExtractPageTexts extracts the text from the specified pages(start 0) of a PDF document.
//
// Parameters:
//   - pages: A variadic parameter representing the page numbers to extract the text from.
//
// Returns:
//   - A string containing the text content of the specified pages.
//   - An error if any error occurs during the extraction process.
func (pp *PdfParser) ExtractPageTexts(pages ...int) (string, error) {
	res := new(strings.Builder)
	for _, page := range pages {
		text, err := pp.pdf.Text(page)
		if err != nil {
			return res.String(), err
		}
		res.WriteString(text)
		res.WriteString(pp.pageSep)
	}

	return res.String(), nil
}

// ExtractTexts extracts the text from all pages of a PDF document.
//
// Parameters:
//   - None
//
// Returns:
//   - A string containing the text content of all pages seperated by the pageSep.
//   - An error if any error occurs during the extraction process.
func (pp *PdfParser) ExtractTexts() (string, error) {
	res := new(strings.Builder)
	for i := 0; i < pp.pdf.NumPage(); i++ {
		text, err := pp.pdf.Text(i)
		if err != nil {
			return res.String(), err
		}
		res.WriteString(text)
		res.WriteString(pp.pageSep)
	}

	return res.String(), nil
}
