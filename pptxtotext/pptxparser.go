// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pptxtotext

import (
	"archive/zip"
	"strings"

	"github.com/young2j/oxmltotext/types"

	"go.uber.org/zap"
)

// PptxParser represents the XML file structure and settings for parsing a pptx file.
type PptxParser struct {
	zipReadCloser *zip.ReadCloser
	slideFiles    map[int]*zip.File
	chartsFiles   map[string]*zip.File
	imagesFiles   map[string]*zip.File
	diagramsFiles map[string]*zip.File
	slideRelsMap  map[int]map[string]string

	parseCharts   bool
	parseImages   bool
	parseDiagrams bool
	drawingsNoFmt bool
	ocr           types.OCR

	slideSep     string
	paragraphSep string
	phraseSep    string
	tableRowSep  string
	tableColSep  string

	logger         *zap.Logger
	disableLogging bool
}

func newPptxParser() *PptxParser {
	logger, _ := zap.NewProduction()

	return &PptxParser{
		slideSep:     strings.Repeat("-", 100) + "\n",
		paragraphSep: "\n",
		phraseSep:    " ",
		tableRowSep:  "\n",
		tableColSep:  "\t",

		logger: logger,
	}
}
