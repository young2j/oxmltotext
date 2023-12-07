// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlsxtotext

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/young2j/oxmltotext/utils"
)

var (
	re_SHARED       = regexp.MustCompile(`xl/sharedStrings\.xml`)
	re_SHEET        = regexp.MustCompile(`xl/worksheets/sheet(\d+)\.xml`)
	re_SHEET_RELS   = regexp.MustCompile(`xl/worksheets/_rels/sheet(\d+)\.xml\.rels`)
	re_CHARTS       = regexp.MustCompile(`xl/charts/chart\d+\.xml`)
	re_IMAGES       = regexp.MustCompile(`xl/media/image\d+\.(?:png|gif|jpg|jpeg)`)
	re_DIAGRAMS     = regexp.MustCompile(`xl/diagrams/data\d+\.xml`)
	re_DRAWINGS     = regexp.MustCompile(`xl/drawings/drawing\d+\.xml`)
	re_DRAWING_RELS = regexp.MustCompile(`xl/drawings/_rels/drawing(\d+)\.xml\.rels`)
)

// Open opens the specified xlsx file path and returns a new XlsxParser instance and an error, if any.
//
// Parameters:
//   - path: a string representing the path to the xlsx file.
//
// Returns:
//   - *XlsxParser: a pointer to the XlsxParser struct.
//   - error: an error, if any.
func Open(path string) (*XlsxParser, error) {
	xp := newXlsxParser()
	zipRc, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	xp.zipReadCloser = zipRc

	if err = matchZipFile(xp, &zipRc.Reader); err != nil {
		xp.Close()
		return nil, err
	}

	return xp, nil
}

// OpenReader opens a DocxParser for the given io.ReaderAt and file size.
//
// Parameters:
//   - r: The io.ReaderAt to read the docx file from.
//   - n: The size of the docx file.
//
// Returns:
//   - *DocxParser: The opened DocxParser object.
//   - error: Any error that occurred during the opening process.
func OpenReader(r io.ReaderAt, n int64) (*XlsxParser, error) {
	xp := newXlsxParser()
	zipReader, err := zip.NewReader(r, n)
	if err != nil {
		return nil, err
	}

	if err = matchZipFile(xp, zipReader); err != nil {
		xp.Close()
		return nil, err
	}

	return xp, nil
}

// OpenURL opens the specified xlsx file URL and returns a XlsxParser, status code, and error.
//
// Parameters:
//   - u (string): The URL to open.
//
// Returns:
//   - *XlsxParser: A pointer to a XlsxParser.
//   - int: The status code.
//   - error: An error object.
func OpenURL(u string) (*XlsxParser, int, error) {
	resp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return nil, statusCode, err
	}

	r := bytes.NewReader(resp.Body)
	xp, err := OpenReader(r, r.Size())

	return xp, statusCode, err
}

// matchZipFile iterates through the files in a zip.Reader
// and populates various maps in the XlsxParser struct
// based on the file names and their contents.
//
// Parameters:
// - xp: A pointer to an XlsxParser struct.
// - r: A pointer to a zip.Reader struct.
//
// Returns:
// - error: An error if any occurred during the iteration and population process.
func matchZipFile(xp *XlsxParser, r *zip.Reader) error {
	sheetsNum := max(len(r.File)-10, 1)
	xp.sheetFiles = make(map[int]*zip.File, sheetsNum)
	xp.chartsFiles = make(map[string]*zip.File, 4)
	xp.imagesFiles = make(map[string]*zip.File, 4)
	xp.diagramsFiles = make(map[string]*zip.File, 4)
	xp.drawingsFile = make(map[string]*zip.File, 4)
	xp.sheetRelsMap = make(map[int]map[string]string, sheetsNum)
	xp.drawingRelsMap = make(map[string]map[string]string, sheetsNum)

	for _, file := range r.File {
		switch {
		case re_SHARED.MatchString(file.Name):
			xp.sharedStringsFile = file
		case re_CHARTS.MatchString(file.Name):
			xp.chartsFiles[file.Name] = file
		case re_IMAGES.MatchString(file.Name):
			xp.imagesFiles[file.Name] = file
		case re_DIAGRAMS.MatchString(file.Name):
			xp.diagramsFiles[file.Name] = file
		case re_DRAWINGS.MatchString(file.Name):
			xp.drawingsFile[file.Name] = file

		default:
			matches := re_SHEET.FindStringSubmatch(file.Name)
			if len(matches) > 1 {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					xp.logWarn(err)
					continue
				}
				xp.sheetFiles[i] = file
				continue
			}

			matches = re_SHEET_RELS.FindStringSubmatch(file.Name)
			if len(matches) > 1 {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					xp.logWarn(err)
					continue
				}
				relsMap, err := utils.ParseRelsMap(file, "xl/")
				if err != nil {
					return err
				}
				xp.sheetRelsMap[i] = relsMap
				continue
			}

			matches = re_DRAWING_RELS.FindStringSubmatch(file.Name)
			if len(matches) > 1 {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					xp.logWarn(err)
					continue
				}
				relsMap, err := utils.ParseRelsMap(file, "xl/")
				if err != nil {
					return err
				}
				k := fmt.Sprintf("xl/drawings/drawing%d.xml", i)
				xp.drawingRelsMap[k] = relsMap
				continue
			}
		}
	}

	return nil
}
