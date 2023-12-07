// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package docxtotext

import (
	"archive/zip"
	"strings"

	"github.com/young2j/oxmltotext/types"

	"go.uber.org/zap"
)

// DocxParser represents the XML file structure and settings for parsing a docx file.
type DocxParser struct {
	zipReadCloser *zip.ReadCloser
	documentFile  *zip.File
	commentsFile  *zip.File
	headerFiles   []*zip.File
	footerFiles   []*zip.File
	footnotesFile *zip.File
	endnotesFile  *zip.File
	chartsFiles   map[string]*zip.File
	imagesFiles   map[string]*zip.File
	diagramsFiles map[string]*zip.File
	docRelsMap    map[string]string
	ocr           types.OCR

	parseComments  bool
	parseHeaders   bool
	parseFooters   bool
	parseFootnotes bool
	parseEndnotes  bool
	parseCharts    bool
	parseImages    bool
	parseDiagrams  bool
	drawingsNoFmt  bool

	paragraphSep string
	partSep      string
	tableRowSep  string
	tableColSep  string

	logger         *zap.Logger
	disableLogging bool
}

func newDocxParser() *DocxParser {
	logger, _ := zap.NewProduction()

	return &DocxParser{
		parseComments:  true,
		parseEndnotes:  true,
		parseFootnotes: true,
		parseFooters:   true,
		parseHeaders:   true,
		paragraphSep:   "\n",
		partSep:        strings.Repeat("-", 100) + "\n",
		tableRowSep:    "\n",
		tableColSep:    "\t",
		logger:         logger,
	}
}
