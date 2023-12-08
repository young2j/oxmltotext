// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlsxtotext

import (
	"image"
	"strconv"
	"strings"

	"github.com/young2j/oxmltotext/ocr"
	"github.com/young2j/oxmltotext/types"

	qxml "github.com/dgrr/quickxml"
)

// SetOnlySharedStrings sets only parsing shared strings or not. Default is false.
func (xp *XlsxParser) SetOnlySharedStrings(v bool) {
	xp.onlySharedStrings = v
}

// SetSheetSep sets the separator of the sheet text. Default is "-"x100.
func (xp *XlsxParser) SetSheetSep(sep string) {
	xp.sheetSep = sep
}

// SetRowSep sets the separator of the row text. Default is "\n".
func (xp *XlsxParser) SetRowSep(sep string) {
	xp.rowSep = sep
}

// SetColSep sets the separator of the column text. Default is "\t".
func (xp *XlsxParser) SetColSep(sep string) {
	xp.colSep = sep
}

// SetParseCharts parses charts or not. Default is false.
func (xp *XlsxParser) SetParseCharts(v bool) {
	xp.parseCharts = v
}

// SetParseDiagrams parses diagrams or not. Default is false.
func (xp *XlsxParser) SetParseDiagrams(v bool) {
	xp.parseDiagrams = v
}

// SetParseImages parses images or not. Default is false.
// When ocr interface is not set, default tesseract-ocr will be used.
func (xp *XlsxParser) SetParseImages(v bool) {
	xp.parseImages = v

	if v && xp.ocr == nil {
		xp.ocr = ocr.NewDefaultOcr()
	}
}

// SetDrawingsNoFmt sets drawings text no outline format.
func (xp *XlsxParser) SetDrawingsNoFmt(v bool) {
	xp.drawingsNoFmt = v
}

// SetOcrInterface overrides default ocr interface.
func (dp *XlsxParser) SetOcrInterface(ocr types.OCR) {
	dp.ocr = ocr
}

// SetDisableLogging sets disable logging.
func (xp *XlsxParser) SetDisableLogging(v bool) {
	xp.disableLogging = v
}

// NumSheets returns the number of sheets.
func (xp *XlsxParser) NumSheets() int {
	return len(xp.sheetFiles)
}

// Close closes the zipReader and OCR client.
// After extracting the text, please remember to call this method.
func (xp *XlsxParser) Close() (err error) {
	if xp.zipReadCloser != nil {
		err = xp.zipReadCloser.Close()
		if err != nil {
			return
		}
	}
	if xp.ocr != nil {
		err = xp.ocr.Close()
		if err != nil {
			return
		}
	}

	return nil
}

// ExtractImages extracts images from the xlsx file.
//
// Parameters:
//   - None
//
// Returns:
//   - []types.Image: a slice of images extracted from the xlsx file.
//   - error: an error if any occurred during the extraction process.
func (xp *XlsxParser) ExtractImages() ([]types.Image, error) {
	images := make([]types.Image, 0, len(xp.imagesFiles))
	for name, f := range xp.imagesFiles {
		r, err := f.Open()
		if err != nil {
			return images, err
		}
		img, format, err := image.Decode(r)
		if err != nil {
			r.Close()
			return images, err
		}
		r.Close()

		images = append(images, types.Image{
			Raw:    img,
			Name:   name,
			Format: format,
		})
	}

	return images, nil
}

// ExtractSheetTexts extracts the texts from the specified xlsx sheets(start 1).
//
// It takes in one or more sheet numbers as parameters and returns a string
// containing the extracted texts. The function also returns an error if there
// is any issue with parsing the sheets.
//
// Parameters:
//   - sheets: An integer slice containing the sheet numbers to extract texts from.
//
// Returns:
//   - string: A string containing the extracted texts.
//   - error: An error object if there is any issue with parsing the sheets.
func (xp *XlsxParser) ExtractSheetTexts(sheets ...int) (string, error) {
	if !xp.shareParsed {
		err := xp.parseSharedStrings()
		if err != nil {
			return "", err
		}
	}

	texts := new(strings.Builder)
	for _, st := range sheets {
		sheet, err := xp.parseSheet(st)
		if err != nil {
			return texts.String(), err
		}
		if sheet.Len() > 0 {
			texts.WriteString(sheet.String())
			texts.WriteString(xp.sheetSep)
		}
	}

	return texts.String(), nil
}

// ExtractTexts extracts the texts from the xlsx file.
//
// It iterates through each sheet of the xlsx file and appends the text content
// to a strings.Builder object. The extracted texts are then returned as a string.
//
// If onlySharedStrings is set to true, only shared strings will be extracted.
//
// If there is an error encountered during the parsing of a sheet, the function
// returns the extracted texts up to that point, along with the error.
//
// Returns:
//   - string: The extracted texts from the xlsx file.
//   - error: An error, if any, encountered during the parsing of the sheets.
func (xp *XlsxParser) ExtractTexts() (string, error) {
	if !xp.shareParsed {
		err := xp.parseSharedStrings()
		if err != nil {
			return "", err
		}
	}

	texts := new(strings.Builder)
	if xp.onlySharedStrings {
		for _, v := range xp.sharedStringsMap {
			texts.WriteString(*v)
			texts.WriteString(xp.rowSep)
		}
		return texts.String(), nil
	}

	for i := 1; i <= xp.NumSheets(); i++ {
		sheet, err := xp.parseSheet(i)
		if err != nil {
			return texts.String(), err
		}
		if sheet.Len() > 0 {
			texts.WriteString(sheet.String())
			texts.WriteString(xp.sheetSep)
		}
	}

	return texts.String(), nil
}

// parseSharedStrings parses the shared strings in the xlsx file.
//
// It opens the shared strings file and reads the XML elements to populate the
// sharedStringsMap with the values. It assigns the values to the corresponding
// indices in the sharedStringsMap. It returns an error if there is any issue
// with opening the file or parsing the XML elements.
//
// Returns:
//   - error: An error if there is any issue with opening the file or parsing
//     the XML elements.
func (xp *XlsxParser) parseSharedStrings() error {
	rc, err := xp.sharedStringsFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	var c = 0
	r := qxml.NewReader(rc)
NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			switch e.Name() {
			case "sst":
				cap := 0
				uniqueCount := e.Attrs().Get("uniqueCount")
				if uniqueCount != nil {
					cap, _ = strconv.Atoi(uniqueCount.Value())
				}
				xp.sharedStringsMap = make(map[string]*string, cap)
				for i := 0; i < cap; i++ {
					xp.sharedStringsMap[strconv.Itoa(i)] = new(string)
				}

			case "t":
				if v, ok := xp.sharedStringsMap[strconv.Itoa(c)]; ok {
					r.AssignNext(v)
					if !r.Next() {
						break NEXT
					}
				}
				c++
			}
		}
	}

	xp.shareParsed = true

	return nil
}

// parseSheet parses a sheet at the given index and returns the extracted texts, tables, charts, diagrams, and images.
//
// Parameters:
//   - i: the index of the sheet to parse.
//
// Returns:
//   - texts: a strings.Builder containing the extracted texts.
//   - error: an error if the slide does not exist or if there was an error opening the sheet file.
func (xp *XlsxParser) parseSheet(i int) (*strings.Builder, error) {
	sheetFile, ok := xp.sheetFiles[i]
	if !ok {
		return nil, types.ErrNoSheet
	}

	rc, err := sheetFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts            = new(strings.Builder)
		cellValueOrIndex = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:
			if e.Name() == "row" {
				texts.WriteString(xp.rowSep)
			}
		case *qxml.StartElement:
			switch e.Name() {
			case "v":
				r.AssignNext(&cellValueOrIndex)
				if !r.Next() {
					break NEXT
				}

				s, ok := xp.sharedStringsMap[cellValueOrIndex]
				if ok {
					texts.WriteString(*s)
				} else {
					texts.WriteString(cellValueOrIndex)
				}
				texts.WriteString(xp.colSep)
				cellValueOrIndex = ""

			case "drawing":
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:id")
					drawings := xp.extractDrawings(i, rIdKV.Value())
					if drawings != nil {
						texts.WriteString(drawings.String())
					}
				}
			}
		}
	}

	return texts, nil
}

func (xp *XlsxParser) logWarn(err error) {
	if xp.disableLogging {
		return
	}
	if err != nil {
		xp.logger.Warn(err.Error())
	}
}
