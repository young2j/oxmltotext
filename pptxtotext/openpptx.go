// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pptxtotext

import (
	"archive/zip"
	"bytes"
	"io"
	"regexp"
	"strconv"

	"github.com/young2j/oxmltotext/utils"
)

var (
	re_SLIDE      = regexp.MustCompile(`ppt/slides/slide(\d+)\.xml`)
	re_SLIDE_RELS = regexp.MustCompile(`ppt/slides/_rels/slide(\d+)\.xml\.rels`)
	re_CHARTS     = regexp.MustCompile(`ppt/charts/chart\d+\.xml`)
	re_IMAGES     = regexp.MustCompile(`ppt/media/image\d+\.(?:png|gif|jpg|jpeg)`)
	re_DIAGRAMS   = regexp.MustCompile(`ppt/diagrams/data\d+\.xml`)
)

// Open opens the specified pptx file path and returns a new PptxParser instance and an error, if any.
//
// Parameters:
//   - path: a string representing the path to the pptx file.
//
// Returns:
//   - *PptxParser: a pointer to the PptxParser struct.
//   - error: an error, if any.
func Open(path string) (*PptxParser, error) {
	pp := newPptxParser()
	zipRc, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	pp.zipReadCloser = zipRc

	if err = matchZipFile(pp, &zipRc.Reader); err != nil {
		pp.Close()
		return nil, err
	}

	return pp, nil
}

// OpenReader opens a PptxParser for the given io.ReaderAt and file size.
//
// Parameters:
//   - r: The io.ReaderAt to read the pptx file from.
//   - n: The size of the pptx file.
//
// Returns:
//   - *PptxParser: The opened PptxParser object.
//   - error: Any error that occurred during the opening process.
func OpenReader(r io.ReaderAt, n int64) (*PptxParser, error) {
	pp := newPptxParser()
	zipReader, err := zip.NewReader(r, n)
	if err != nil {
		return nil, err
	}

	if err = matchZipFile(pp, zipReader); err != nil {
		pp.Close()
		return nil, err
	}

	return pp, nil
}

// OpenURL opens the specified pptx file URL and returns a PptxParser, status code, and error.
//
// Parameters:
//   - u (string): The URL to open.
//
// Returns:
//   - *PptxParser: A pointer to a PptxParser.
//   - int: The status code.
//   - error: An error object.
func OpenURL(u string) (*PptxParser, int, error) {
	resp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return nil, statusCode, err
	}

	r := bytes.NewReader(resp.Body)
	pp, err := OpenReader(r, r.Size())

	return pp, statusCode, err
}

// matchZipFile is a function that matches the files in a zip.Reader to specific categories such as slides, charts,
// images, and diagrams. It populates the relevant maps in the PptxParser struct with the matched files.
//
// Parameters:
//   - pp: a pointer to the PptxParser struct that holds the maps for slideFiles, chartsFiles, imagesFiles,
//     diagramsFiles, and slideRelsMap.
//   - r: a pointer to the zip.Reader struct that contains the files to be matched.
//
// Return:
//   - error: an error if there was a problem parsing the files or populating the maps, otherwise nil.
func matchZipFile(pp *PptxParser, r *zip.Reader) error {
	slidesNum := max((len(r.File)-16)/3, 1)
	pp.slideFiles = make(map[int]*zip.File, slidesNum)
	pp.chartsFiles = make(map[string]*zip.File, 4)
	pp.imagesFiles = make(map[string]*zip.File, 4)
	pp.diagramsFiles = make(map[string]*zip.File, 4)
	pp.slideRelsMap = make(map[int]map[string]string, slidesNum)

	for _, file := range r.File {
		switch {
		case re_CHARTS.MatchString(file.Name):
			pp.chartsFiles[file.Name] = file
		case re_IMAGES.MatchString(file.Name):
			pp.imagesFiles[file.Name] = file
		case re_DIAGRAMS.MatchString(file.Name):
			pp.diagramsFiles[file.Name] = file
		default:
			matches := re_SLIDE.FindStringSubmatch(file.Name)
			if len(matches) > 1 {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					pp.logWarn(err)
					continue
				}
				pp.slideFiles[i] = file
				continue
			}

			matches = re_SLIDE_RELS.FindStringSubmatch(file.Name)
			if len(matches) > 1 {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					pp.logWarn(err)
					continue
				}
				relsMap, err := utils.ParseRelsMap(file, "ppt/")
				if err != nil {
					return err
				}
				pp.slideRelsMap[i] = relsMap
				continue
			}
		}
	}

	return nil
}
