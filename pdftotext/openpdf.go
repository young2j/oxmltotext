// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pdftotext

import (
	"io"

	"github.com/young2j/oxmltotext/utils"

	"github.com/gen2brain/go-fitz"
)

// Open creates a new PdfParser and opens a PDF file at the specified path.
//
// Parameters:
//   - path: The path of the PDF file to open.
//
// Returns:
//   - *PdfParser: A pointer to the PdfParser object if the file was opened successfully.
//   - error: An error object if there was an error opening the file.
func Open(path string) (*PdfParser, error) {
	pp := newPdfParser()
	pdf, err := fitz.New(path)
	if err != nil {
		return nil, err
	}
	pp.pdf = pdf

	return pp, nil
}

// OpenReader creates a new PdfParser from an io.Reader.
//
// Parameters:
//   - r: The io.Reader from which to create the PdfParser.
//
// Returns:
//   - *PdfParser: The created PdfParser.
//   - error: Any error that occurred during the creation of the PdfParser.
func OpenReader(r io.Reader) (*PdfParser, error) {
	pp := newPdfParser()
	pdf, err := fitz.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	pp.pdf = pdf

	return pp, nil
}

// OpenURL creates a new PdfParser by reading the specified URL.
//
// Parameters:
//   - u: the URL to open as a string.
//
// Returns:
//   - pp: a pointer to a PdfParser object.
//   - statusCode: an integer representing the HTTP status code of the URL response.
//   - err: an error object, if any error occurred during the process.
func OpenURL(u string) (*PdfParser, int, error) {
	pp := newPdfParser()
	resp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return nil, statusCode, err
	}
	pdf, err := fitz.NewFromMemory(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}
	pp.pdf = pdf

	return pp, statusCode, nil
}
