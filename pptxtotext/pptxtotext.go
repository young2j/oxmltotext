// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pptxtotext

import (
	"strings"

	"github.com/young2j/oxmltotext/ocr"
	"github.com/young2j/oxmltotext/types"

	qxml "github.com/dgrr/quickxml"
)

// SetSlideSep sets slide text separator. Default is "-"x100.
func (pp *PptxParser) SetSlideSep(sep string) {
	pp.slideSep = sep
}

// SetParagraphSep sets phrase separator. Default is " ".
func (pp *PptxParser) SetPhraseSep(sep string) {
	pp.phraseSep = sep
}

// SetTableRowSep sets table row separator. Default is "\n".
func (pp *PptxParser) SetTableRowSep(sep string) {
	pp.tableRowSep = sep
}

// SetTableColSep sets table column separator. Default is "\t".
func (pp *PptxParser) SetTableColSep(sep string) {
	pp.tableColSep = sep
}

// SetParseCharts parses charts or not. Default is false.
func (pp *PptxParser) SetParseCharts(v bool) {
	pp.parseCharts = v
}

// SetParseDiagrams parses diagrams or not. Default is false.
func (pp *PptxParser) SetParseDiagrams(v bool) {
	pp.parseDiagrams = v
}

// SetParseImages parses images or not. Default is false.
// When ocr interface is not set, default tesseract-ocr will be used.
func (pp *PptxParser) SetParseImages(v bool) {
	pp.parseImages = v

	if v && pp.ocr == nil {
		pp.ocr = ocr.NewDefaultOcr()
	}
}

// SetDrawingsNoFmt sets drawings text no outline format.
func (pp *PptxParser) SetDrawingsNoFmt(v bool) {
	pp.drawingsNoFmt = v
}

// SetOcrInterface overrides default ocr interface.
func (pp *PptxParser) SetOcrInterface(ocr types.OCR) {
	pp.ocr = ocr
}

// DisableLogging disables logging.
func (pp *PptxParser) DisableLogging(v bool) {
	pp.disableLogging = v
}

// NumSlides returns the number of slides.
func (pp *PptxParser) NumSlides() int {
	return len(pp.slideFiles)
}

// Close closes the zipReader and OCR client.
// After extracting the text, please remember to call this method.
func (pp *PptxParser) Close() (err error) {
	if pp.zipReadCloser != nil {
		err = pp.zipReadCloser.Close()
		if err != nil {
			return
		}
	}
	if pp.ocr != nil {
		err = pp.ocr.Close()
		if err != nil {
			return
		}
	}

	return nil
}

// ExtractSlideTexts extracts the texts from the specified pptx slides(start 1).
//
// It takes in one or more slide numbers as parameters and returns a string
// containing the extracted texts. The function also returns an error if there
// is any issue with parsing the slides.
//
// Parameters:
//   - slides: An integer slice containing the slide numbers to extract texts from.
//
// Returns:
//   - string: A string containing the extracted texts.
//   - error: An error object if there is any issue with parsing the slides.
func (pp *PptxParser) ExtractSlideTexts(slides ...int) (string, error) {
	texts := new(strings.Builder)
	for _, slide := range slides {
		slide, err := pp.parseSlide(slide)
		if err != nil {
			return texts.String(), err
		}
		if slide.Len() > 0 {
			texts.WriteString(slide.String())
			texts.WriteString(pp.slideSep)
		}
	}

	return texts.String(), nil
}

// ExtractTexts extracts the texts from the pptx file.
//
// It iterates through each slide of the pptx file and appends the text content
// to a strings.Builder object. The extracted texts are then returned as a string.
// If there is an error encountered during the parsing of a slide, the function
// returns the extracted texts up to that point, along with the error.
//
// Returns:
//   - string: The extracted texts from the pptx file.
//   - error: An error, if any, encountered during the parsing of the slides.
func (pp *PptxParser) ExtractTexts() (string, error) {
	texts := new(strings.Builder)

	for i := 1; i <= pp.NumSlides(); i++ {
		slide, err := pp.parseSlide(i)
		if err != nil {
			return texts.String(), err
		}
		if slide.Len() > 0 {
			texts.WriteString(slide.String())
			texts.WriteString(pp.slideSep)
		}
	}

	return texts.String(), nil
}

// parseSlide parses a slide at the given index and returns the extracted texts, tables, charts, diagrams, and images.
//
// Parameters:
//   - i: the index of the slide to parse.
//
// Returns:
//   - texts: a strings.Builder containing the extracted texts.
//   - error: an error if the slide does not exist or if there was an error opening the slide file.
func (pp *PptxParser) parseSlide(i int) (*strings.Builder, error) {
	slideFile, ok := pp.slideFiles[i]
	if !ok {
		return nil, types.ErrNoSlide
	}

	rc, err := slideFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts  = new(strings.Builder)
		phrase = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:
			if e.Name() == "a:p" {
				texts.WriteString(pp.paragraphSep)
			}

		case *qxml.StartElement:
			switch e.Name() {
			case "a:t":
				r.AssignNext(&phrase)
				if !r.Next() {
					break NEXT
				}
				if len(phrase) > 0 {
					texts.WriteString(phrase)
					texts.WriteString(pp.phraseSep)
					phrase = ""
				}

			case "a:tbl":
				table := pp.extractTable(r)
				if table != nil {
					texts.WriteString(table.String())
				}

			case "c:chart":
				if !pp.parseCharts {
					continue
				}
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:id")
					chart, err := pp.extractChart(i, rIdKV.Value())
					pp.logWarn(err)
					if chart != nil {
						texts.WriteString(chart.String())
					}
				}

			case "dgm:relIds":
				if !pp.parseDiagrams {
					continue
				}
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:dm")
					diagram, err := pp.extractDiagram(i, rIdKV.Value())
					pp.logWarn(err)
					if diagram != nil {
						texts.WriteString(diagram.String())
					}
				}

			case "a:blip":
				if !pp.parseImages {
					continue
				}
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:embed")
					image, err := pp.extractImage(i, rIdKV.Value())
					pp.logWarn(err)
					if image != nil {
						texts.WriteString(image.String())
					}
				}
			}
		}
	}
	return texts, nil
}

// extractTable extracts table data from a pptx file using a qxml.Reader.
//
// Parameters:
//   - r: a qxml.Reader object used to read the XML elements.
//
// Return type:
//   - *strings.Builder: a strings.Builder object containing the extracted table data.
func (pp *PptxParser) extractTable(r *qxml.Reader) *strings.Builder {
	var (
		texts = new(strings.Builder)
		row   = new(strings.Builder)
		a_t   = ""
	)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "a:t" {
				r.AssignNext(&a_t)
				if !r.Next() {
					break NEXT
				}
				row.WriteString(a_t)
				row.WriteString(pp.tableColSep)
				a_t = ""
			}

		case *qxml.EndElement:
			switch e.Name() {
			case "a:tr":
				if row.Len() > 0 {
					texts.WriteString(row.String())
					texts.WriteString(pp.tableRowSep)
					row.Reset()
					a_t = ""
				}
			case "w:tbl":
				break NEXT
			}
		}
	}

	return texts
}

func (pp *PptxParser) logWarn(err error) {
	if pp.disableLogging {
		return
	}
	if err != nil {
		pp.logger.Warn(err.Error())
	}
}
