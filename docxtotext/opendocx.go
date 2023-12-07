// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package docxtotext

import (
	"archive/zip"
	"bytes"
	"io"
	"regexp"

	"github.com/young2j/oxmltotext/utils"
)

var (
	re_DOCUMENT  = regexp.MustCompile(`word/document\.xml`)
	re_COMMENTS  = regexp.MustCompile(`word/comments\.xml`)
	re_ENDNOTES  = regexp.MustCompile(`word/endnotes\.xml`)
	re_FOOTNOTES = regexp.MustCompile(`word/footnotes\.xml`)
	re_FOOTER    = regexp.MustCompile(`word/footer\d+\.xml`)
	re_HEADER    = regexp.MustCompile(`word/header\d+\.xml`)
	re_DOC_RELS  = regexp.MustCompile(`word/_rels/document\.xml\.rels`)
	re_CHARTS    = regexp.MustCompile(`word/charts/chart\d+\.xml`)
	re_IMAGES    = regexp.MustCompile(`word/media/image\d+\.(?:png|gif|jpg|jpeg)`)
	re_DIAGRAMS  = regexp.MustCompile(`word/diagrams/data\d+\.xml`)
)

// Open opens the specified docx file path and returns a new DocxParser instance and an error, if any.
//
// Parameters:
//   - path: a string representing the path to the docx file.
//
// Returns:
//   - *DocxParser: a pointer to the DocxParser struct.
//   - error: an error, if any.
func Open(path string) (*DocxParser, error) {
	dp := newDocxParser()
	zipRc, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	dp.zipReadCloser = zipRc

	if err = matchZipFile(dp, &zipRc.Reader); err != nil {
		dp.Close()
		return nil, err
	}

	return dp, nil
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
func OpenReader(r io.ReaderAt, n int64) (*DocxParser, error) {
	dp := newDocxParser()
	zipReader, err := zip.NewReader(r, n)
	if err != nil {
		return nil, err
	}

	if err = matchZipFile(dp, zipReader); err != nil {
		dp.Close()
		return nil, err
	}

	return dp, nil
}

// OpenURL opens the specified docx file URL and returns a DocxParser, status code, and error.
//
// Parameters:
//   - u (string): The URL to open.
//
// Returns:
//   - *DocxParser: A pointer to a DocxParser.
//   - int: The status code.
//   - error: An error object.
func OpenURL(u string) (*DocxParser, int, error) {
	resp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return nil, statusCode, err
	}

	r := bytes.NewReader(resp.Body)
	dp, err := OpenReader(r, r.Size())

	return dp, statusCode, err
}

// matchZipFile matches the zip file with the given DocxParser and zip.Reader.
//
// It populates the footerFiles, headerFiles, chartsFiles, imagesFiles, and diagramsFiles
// fields of the DocxParser based on the files found in the zip.Reader. It also sets the
// documentFile, commentsFile, endnotesFile, footnotesFile, and docRelsMap fields if the
// corresponding files are found in the zip.Reader.
//
// Parameters:
//   - dp: A pointer to the DocxParser instance.
//   - r: A pointer to the zip.Reader instance.
//
// Returns:
//   - error: An error if any occurs during the matching process.
func matchZipFile(dp *DocxParser, r *zip.Reader) error {
	dp.footerFiles = make([]*zip.File, 0, max(len(r.File)-20, 1))
	dp.headerFiles = make([]*zip.File, 0, max(len(r.File)-20, 1))
	dp.chartsFiles = make(map[string]*zip.File, 4)
	dp.imagesFiles = make(map[string]*zip.File, 4)
	dp.diagramsFiles = make(map[string]*zip.File, 4)
	for _, file := range r.File {
		switch {
		case re_DOCUMENT.MatchString(file.Name):
			dp.documentFile = file
		case re_COMMENTS.MatchString(file.Name):
			dp.commentsFile = file
		case re_ENDNOTES.MatchString(file.Name):
			dp.endnotesFile = file
		case re_FOOTNOTES.MatchString(file.Name):
			dp.footnotesFile = file
		case re_FOOTER.MatchString(file.Name):
			dp.footerFiles = append(dp.footerFiles, file)
		case re_HEADER.MatchString(file.Name):
			dp.headerFiles = append(dp.headerFiles, file)
		case re_CHARTS.MatchString(file.Name):
			dp.chartsFiles[file.Name] = file
		case re_IMAGES.MatchString(file.Name):
			dp.imagesFiles[file.Name] = file
		case re_DIAGRAMS.MatchString(file.Name):
			dp.diagramsFiles[file.Name] = file
		case re_DOC_RELS.MatchString(file.Name):
			relsMap, err := utils.ParseRelsMap(file, "word/")
			if err != nil {
				return err
			}
			dp.docRelsMap = relsMap
		}
	}

	return nil
}
