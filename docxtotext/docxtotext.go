// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package docxtotext

import (
	"image"
	"strings"

	"github.com/young2j/oxmltotext/ocr"
	"github.com/young2j/oxmltotext/types"

	qxml "github.com/dgrr/quickxml"
)

// SetParagraphSep sets paragraph separator. Default is "\n".
func (dp *DocxParser) SetParagraphSep(sep string) {
	dp.paragraphSep = sep
}

// SetPartSep sets document part(every XML file like header, footer, etc.) separator. Default is "-"x100.
func (dp *DocxParser) SetPartSep(sep string) {
	dp.partSep = sep
}

// SetTableRowSep sets table row separator. Default is "\n".
func (dp *DocxParser) SetTableRowSep(sep string) {
	dp.tableRowSep = sep
}

// SetTableColSep sets table column separator. Default is "\t".
func (dp *DocxParser) SetTableColSep(sep string) {
	dp.tableColSep = sep
}

// SetParseComments parses comments or not. Default is true.
func (dp *DocxParser) SetParseComments(v bool) {
	dp.parseComments = v
}

// SetParseEndnotes parses endnotes or not. Default is true.
func (dp *DocxParser) SetParseEndnotes(v bool) {
	dp.parseEndnotes = v
}

// SetParseFootnotes parses footnotes or not. Default is true.
func (dp *DocxParser) SetParseFootnotes(v bool) {
	dp.parseFootnotes = v
}

// SetParseFooters parses footers or not. Default is true.
func (dp *DocxParser) SetParseFooters(v bool) {
	dp.parseFooters = v
}

// SetParseHeaders parses headers or not. Default is true.
func (dp *DocxParser) SetParseHeaders(v bool) {
	dp.parseHeaders = v
}

// SetParseCharts parses charts or not. Default is false.
func (dp *DocxParser) SetParseCharts(v bool) {
	dp.parseCharts = v
}

// SetParseDiagrams parses diagrams or not. Default is false.
func (dp *DocxParser) SetParseDiagrams(v bool) {
	dp.parseDiagrams = v
}

// SetParseImages parses images or not. Default is false.
// When ocr interface is not set, default tesseract-ocr will be used.
func (dp *DocxParser) SetParseImages(v bool) {
	dp.parseImages = v

	if v && dp.ocr == nil {
		dp.ocr = ocr.NewDefaultOcr()
	}
}

// SetDrawingsNoFmt sets drawings text no outline format.
func (dp *DocxParser) SetDrawingsNoFmt(v bool) {
	dp.drawingsNoFmt = v
}

// SetOcrInterface overrides default ocr interface.
func (dp *DocxParser) SetOcrInterface(ocr types.OCR) {
	dp.ocr = ocr
}

// DisableLogging disables logging.
func (dp *DocxParser) DisableLogging(v bool) {
	dp.disableLogging = v
}

// Close closes the zipReader and OCR client.
// After extracting the text, please remember to call this method.
func (dp *DocxParser) Close() (err error) {
	if dp.zipReadCloser != nil {
		err = dp.zipReadCloser.Close()
		if err != nil {
			return
		}
	}
	if dp.ocr != nil {
		err = dp.ocr.Close()
		if err != nil {
			return
		}
	}

	return nil
}

// ExtractImages extracts images from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - []types.Image: a slice of images extracted from the docx file.
//   - error: an error if any occurred during the extraction process.
func (dp *DocxParser) ExtractImages() ([]types.Image, error) {
	images := make([]types.Image, 0, len(dp.imagesFiles))
	for name, f := range dp.imagesFiles {
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

// ExtractTexts extracts the texts from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - string: The extracted texts.
//   - error: An error if any.
func (dp *DocxParser) ExtractTexts() (string, error) {
	texts, err := dp.extractDocument()
	if err != nil {
		return "", err
	}

	if dp.parseComments && dp.commentsFile != nil {
		comments, err := dp.extractComments()
		if err != nil {
			return texts.String(), err
		}
		if comments.Len() > 0 {
			texts.WriteString(dp.partSep)
			texts.WriteString(comments.String())
		}
	}

	if dp.parseHeaders {
		headers, err := dp.extractHeaders()
		if err != nil {
			return texts.String(), err
		}
		if headers.Len() > 0 {
			texts.WriteString(dp.partSep)
			texts.WriteString(headers.String())
		}
	}

	if dp.parseFooters {
		footers, err := dp.extractFooters()
		if err != nil {
			return texts.String(), err
		}
		if footers.Len() > 0 {
			texts.WriteString(dp.partSep)
			texts.WriteString(footers.String())
		}
	}

	if dp.parseFootnotes && dp.footnotesFile != nil {
		footnotes, err := dp.extractFootnotes()
		if err != nil {
			return texts.String(), err
		}
		if footnotes.Len() > 0 {
			texts.WriteString(dp.partSep)
			texts.WriteString(footnotes.String())
		}
	}

	if dp.parseEndnotes && dp.endnotesFile != nil {
		endnotes, err := dp.extractEndnotes()
		if err != nil {
			return texts.String(), err
		}
		if endnotes.Len() > 0 {
			texts.WriteString(dp.partSep)
			texts.WriteString(endnotes.String())
		}
	}

	return texts.String(), nil
}

// extractDocument extracts the text of document part from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted document text.
//   - error: an error if any.
func (dp *DocxParser) extractDocument() (*strings.Builder, error) {
	if dp.documentFile == nil {
		dp.logWarn(types.ErrNoDocument)
		return new(strings.Builder), nil
	}

	rc, err := dp.documentFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts     = new(strings.Builder)
		paragraph = new(strings.Builder)
		w_t       = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			switch e.Name() {
			case "w:t":
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					paragraph.WriteString(w_t)
					w_t = ""
				}

			case "w:tbl":
				table := dp.extractTable(r)
				if table != nil {
					texts.WriteString(table.String())
				}

			case "w:drawing":
				drawings := dp.extractDrawings(r)
				if drawings != nil {
					texts.WriteString(drawings.String())
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:p" {
				if paragraph.Len() > 0 {
					texts.WriteString(paragraph.String())
					texts.WriteString(dp.paragraphSep)
					paragraph.Reset()
					w_t = ""
				}
			}
		}
	}

	return texts, nil
}

// extractTable extracts the table from the given qxml.Reader and returns a strings.Builder with the extracted table contents.
//
// Parameters:
//   - r: a qxml.Reader instance from which the table is extracted.
//
// Return:
//   - texts: a strings.Builder instance containing the extracted table contents.
func (dp *DocxParser) extractTable(r *qxml.Reader) *strings.Builder {
	var (
		texts = new(strings.Builder)
		row   = new(strings.Builder)
		w_t   = ""
	)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				row.WriteString(w_t)
				row.WriteString(dp.tableColSep)
				w_t = ""
			}

		case *qxml.EndElement:
			switch e.Name() {
			case "w:tr":
				if row.Len() > 0 {
					texts.WriteString(row.String())
					texts.WriteString(dp.tableRowSep)
					row.Reset()
					w_t = ""
				}
			case "w:tbl":
				break NEXT
			}
		}
	}

	return texts
}

// / extractComments extracts the text of comments part from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted comments text.
//   - error: an error if any.
func (dp *DocxParser) extractComments() (*strings.Builder, error) {
	if dp.commentsFile == nil {
		dp.logWarn(types.ErrNoComments)
		return new(strings.Builder), nil
	}

	rc, err := dp.commentsFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts = new(strings.Builder)
		w_t   = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					texts.WriteString(w_t)
					w_t = ""
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:comment" {
				texts.WriteString(dp.paragraphSep)
			}
		}
	}

	return texts, nil
}

// extractEndnotes extracts the text of endnotes part from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted endnotes text.
//   - error: an error if any.
func (dp *DocxParser) extractEndnotes() (*strings.Builder, error) {
	if dp.endnotesFile == nil {
		dp.logWarn(types.ErrNoEndnotes)
		return new(strings.Builder), nil
	}

	rc, err := dp.endnotesFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts   = new(strings.Builder)
		endnote = new(strings.Builder)
		w_t     = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					endnote.WriteString(w_t)
					w_t = ""
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:endnote" {
				if endnote.Len() > 0 {
					texts.WriteString(endnote.String())
					texts.WriteString(dp.paragraphSep)
					endnote.Reset()
					w_t = ""
				}
			}
		}
	}

	return texts, nil

}

// extractFootnotes extracts the text of footnotes part from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted footnotes text.
//   - error: an error if any.
func (dp *DocxParser) extractFootnotes() (*strings.Builder, error) {
	if dp.footnotesFile == nil {
		dp.logWarn(types.ErrNoFootnotes)
		return new(strings.Builder), nil
	}

	rc, err := dp.footnotesFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts    = new(strings.Builder)
		footnote = new(strings.Builder)
		w_t      = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					footnote.WriteString(w_t)
					w_t = ""
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:footnote" {
				if footnote.Len() > 0 {
					texts.WriteString(footnote.String())
					texts.WriteString(dp.paragraphSep)
					footnote.Reset()
					w_t = ""
				}
			}
		}
	}

	return texts, nil
}

// extractFooters extracts the text of footer parts from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted footers text.
//   - error: an error if any.
func (dp *DocxParser) extractFooters() (*strings.Builder, error) {
	texts := new(strings.Builder)
	for i := range dp.footerFiles {
		footer, err := dp.extractFooter(i)
		if err != nil {
			return nil, err
		}
		texts.WriteString(footer.String())
	}
	return texts, nil
}

// extractFooter extracts the footer from the specified index of the DocxParser's footerFiles.
//
// Parameters:
//   - i: the index of the footer file
//
// Returns:
//   - *strings.Builder: the extracted footer
//   - error: any error that occurred during the extraction process
func (dp *DocxParser) extractFooter(i int) (*strings.Builder, error) {
	rc, err := dp.footerFiles[i].Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts = new(strings.Builder)
		ftr   = new(strings.Builder)
		w_t   = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					ftr.WriteString(w_t)
					w_t = ""
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:p" {
				if ftr.Len() > 0 {
					texts.WriteString(ftr.String())
					texts.WriteString(dp.paragraphSep)
					ftr.Reset()
					w_t = ""
				}
			}
		}
	}

	return texts, nil
}

// extractDocument extracts the text of header parts from the docx file.
//
// Parameters:
//   - None
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted headers text.
//   - error: an error if any.
func (dp *DocxParser) extractHeaders() (*strings.Builder, error) {
	texts := new(strings.Builder)
	for i := range dp.headerFiles {
		header, err := dp.extractHeader(i)
		if err != nil {
			return nil, err
		}
		texts.WriteString(header.String())
	}
	return texts, nil
}

// extractHeader extracts the header from the specified index of the DocxParser's headerFiles.
//
// Parameters:
//   - i: The index of the header file to extract.
//
// Returns:
//   - *strings.Builder: The extracted header text.
//   - error: An error if there was a problem opening or processing the header file.
func (dp *DocxParser) extractHeader(i int) (*strings.Builder, error) {
	rc, err := dp.headerFiles[i].Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		texts = new(strings.Builder)
		hdr   = new(strings.Builder)
		w_t   = ""
	)
	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == "w:t" {
				r.AssignNext(&w_t)
				if !r.Next() {
					break NEXT
				}
				if len(w_t) > 0 {
					hdr.WriteString(w_t)
					w_t = ""
				}
			}

		case *qxml.EndElement:
			if e.Name() == "w:p" {
				if hdr.Len() > 0 {
					texts.WriteString(hdr.String())
					texts.WriteString(dp.paragraphSep)
					hdr.Reset()
					w_t = ""
				}
			}
		}
	}

	return texts, nil
}

func (dp *DocxParser) logWarn(err error) {
	if dp.disableLogging {
		return
	}
	if err != nil {
		dp.logger.Warn(err.Error())
	}
}
