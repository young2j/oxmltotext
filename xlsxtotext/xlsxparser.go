// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlsxtotext

import (
	"archive/zip"
	"strings"

	"github.com/young2j/oxmltotext/types"

	"go.uber.org/zap"
)

// XlsxParser represents the XML file structure and settings for parsing a xlsx file.
type XlsxParser struct {
	zipReadCloser     *zip.ReadCloser
	sharedStringsFile *zip.File
	sharedStringsMap  map[string]*string
	sheetFiles        map[int]*zip.File
	chartsFiles       map[string]*zip.File
	imagesFiles       map[string]*zip.File
	diagramsFiles     map[string]*zip.File
	drawingsFile      map[string]*zip.File
	sheetRelsMap      map[int]map[string]string
	drawingRelsMap    map[string]map[string]string

	parseCharts   bool
	parseImages   bool
	parseDiagrams bool
	drawingsNoFmt bool
	ocr           types.OCR

	onlySharedStrings bool
	sheetSep          string
	rowSep            string
	colSep            string
	shareParsed       bool

	logger         *zap.Logger
	disableLogging bool
}

func newXlsxParser() *XlsxParser {
	logger, _ := zap.NewProduction()

	return &XlsxParser{
		sharedStringsMap: make(map[string]*string, 0),
		sheetSep:         strings.Repeat("-", 100) + "\n",
		rowSep:           "\n",
		colSep:           "\t",
		logger:           logger,
	}
}
