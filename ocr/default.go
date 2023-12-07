// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:build ocr

/*
Package ocr provides default OCR interface implementation for using Tesseract OCR.
*/
package ocr

import (
	"io"

	"github.com/young2j/oxmltotext/types"

	"github.com/otiai10/gosseract/v2"
)

type defaultOcr struct {
	client *gosseract.Client
}

// NewDefaultOcr initializes and returns a new instance of the default OCR implementation.
//
// It creates a new gosseract client and sets the language to "eng", "chi_sim", and "script/HanS".
//
// Returns:
//   - a pointer to defaultOcr struct, which implements the types.OCR interface.
func NewDefaultOcr() types.OCR {
	client := gosseract.NewClient()
	client.SetLanguage("eng", "chi_sim", "script/HanS")

	return &defaultOcr{
		client,
	}
}

// Run runs the OCR on the given input and returns the extracted text.
//
// Parameters:
//   - r: A reader containing the image data.
//
// Returns:
//   - string: The extracted text.
//   - error: Any error that occurred during the OCR process.
func (d *defaultOcr) Run(r io.Reader) (string, error) {
	imgData, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	err = d.client.SetImageFromBytes(imgData)
	if err != nil {
		return "", err
	}

	text, err := d.client.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}

// Close closes the ocr client of defaultOcr instance.
//
// It returns an error if there was a problem closing the client.
func (d *defaultOcr) Close() error {
	return d.client.Close()
}
