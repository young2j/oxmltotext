// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package docxtotext

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/young2j/oxmltotext/types"
	"github.com/young2j/oxmltotext/utils"

	qxml "github.com/dgrr/quickxml"
)

// extractDrawings extracts the drawings from the given qxml.Reader.
//
// The function iterates through the XML elements of the reader and
// extracts the drawings(charts, images, and diagrams) if the corresponding flags are set.
//
// Parameters:
//
//	-r : a qxml.Reader
//
// Returns:
//   - *strings.Builder: a strings.Builder containing the extracted drawings text.
func (dp *DocxParser) extractDrawings(r *qxml.Reader) *strings.Builder {
	var texts = new(strings.Builder)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:
			if e.Name() == "w:drawing" {
				break NEXT
			}

		case *qxml.StartElement:
			switch {
			case e.Name() == "c:chart" && dp.parseCharts:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:id")
					charts, err := dp.extractChart(rIdKV.Value())
					dp.logWarn(err)
					if charts != nil {
						texts.WriteString(charts.String())
					}
				}

			case e.Name() == "a:blip" && dp.parseImages:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:embed")
					images, err := dp.extractImage(rIdKV.Value())
					dp.logWarn(err)
					if images != nil {
						texts.WriteString(images.String())
					}
				}

			case e.Name() == "dgm:relIds" && dp.parseDiagrams:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:dm")
					diagrams, err := dp.extractDiagram(rIdKV.Value())
					dp.logWarn(err)
					if diagrams != nil {
						texts.WriteString(diagrams.String())
					}
				}
			}
		}
	}

	return texts
}

// extractChart extracts text content from the chart.
//
// Parameters:
//   - rId: The reference ID of the chart.
//
// Returns:
//   - *strings.Builder: The formatted text of the chart.
//   - error: An error if the extraction fails.
func (dp *DocxParser) extractChart(rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	fname, ok := dp.docRelsMap[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := dp.chartsFiles[fname]
	if !ok {
		return nil, types.ErrNonePart
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		fmtTexts   = new(strings.Builder)
		texts      = new(strings.Builder)
		line       = new(strings.Builder)
		c_v        = ""
		lineSep    = "\n"
		space      = " "
		maxLineLen = 0
	)

	r := qxml.NewReader(rc)
	valRegex := regexp.MustCompile(`(?i)c:.?val`)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:
			if e.Name() == "c:plotArea" {
				if dp.drawingsNoFmt {
					fmtTexts.WriteString(texts.String())
					fmtTexts.WriteString(lineSep)
					return fmtTexts, nil
				}

				halfLine := bytes.Repeat([]byte("─"), max((maxLineLen-5)/2, 0))
				fmtTexts.WriteString("┌")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("chart")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("┐")
				fmtTexts.WriteString(lineSep)

				fmtTexts.WriteString(texts.String())

				fmtTexts.WriteString("└")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("─────")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("┘")
				fmtTexts.WriteString(lineSep)

				texts.Reset()
				break NEXT
			}
		case *qxml.StartElement:
			switch e.Name() {
			case "c:ser":
			INNER_NEXT:
				for r.Next() {
					switch e := r.Element().(type) {
					case *qxml.EndElement:
						if e.Name() == "c:ser" {
							break INNER_NEXT
						}

					case *qxml.StartElement:
						name := e.Name()
						switch {
						case name == "c:tx":
							if utils.FindNameIterTo(r, "c:v", "c:tx") {
								r.AssignNext(&c_v)
								if !r.Next() {
									break NEXT
								}

								line.WriteString(" [")
								line.WriteString(c_v)
								line.WriteString("]")
								c_v = ""

								if line.Len() > 0 {
									texts.WriteString(line.String())
									texts.WriteString(lineSep)
									if line.Len() > maxLineLen {
										maxLineLen = line.Len()
									}
									line.Reset()
								}
							}

						case name == "c:cat":
							for utils.FindNameIterTo(r, "c:v", "c:cat") {
								r.AssignNext(&c_v)
								if !r.Next() {
									break NEXT
								}
								line.WriteString(c_v)
								line.WriteString(space)
								c_v = ""
							}
							if line.Len() > 0 {
								texts.WriteString(space)
								texts.WriteString(line.String())
								texts.WriteString(lineSep)
								if line.Len() > maxLineLen {
									maxLineLen = line.Len()
								}
								line.Reset()
							}

						case valRegex.MatchString(name):
							for utils.MatchNameIterTo(r, "c:v", `(?i)c:.?val`) {
								r.AssignNext(&c_v)
								if !r.Next() {
									break NEXT
								}
								line.WriteString(c_v)
								line.WriteString(space)
								c_v = ""
							}
							if line.Len() > 0 {
								texts.WriteString(space)
								texts.WriteString(line.String())
								texts.WriteString(lineSep)
								if line.Len() > maxLineLen {
									maxLineLen = line.Len()
								}
								line.Reset()
							}
						}

					}
				}
			}
		}
	}

	return fmtTexts, nil
}

// extractDiagram extracts text content from the diagram.
//
// Parameters:
//   - rId: The reference ID of the diagram.
//
// Returns:
//   - *strings.Builder: The formatted text of the diagram.
//   - error: An error if the extraction fails.
func (dp *DocxParser) extractDiagram(rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	fname, ok := dp.docRelsMap[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := dp.diagramsFiles[fname]
	if !ok {
		return nil, types.ErrNonePart
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		fmtTexts   = new(strings.Builder)
		texts      = new(strings.Builder)
		line       = new(strings.Builder)
		c_v        = ""
		lineSep    = "\n"
		space      = " "
		maxLineLen = 0
	)

	r := qxml.NewReader(rc)

NEXT:
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:
			if e.Name() == "dgm:ptLst" {
				if dp.drawingsNoFmt {
					fmtTexts.WriteString(texts.String())
					fmtTexts.WriteString(lineSep)
					return fmtTexts, nil
				}

				halfLine := bytes.Repeat([]byte("─"), max((maxLineLen-7)/2, 0))
				fmtTexts.WriteString("┌")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("diagram")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("┐")
				fmtTexts.WriteString(lineSep)

				fmtTexts.WriteString(texts.String())

				fmtTexts.WriteString("└")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("───────")
				fmtTexts.Write(halfLine)
				fmtTexts.WriteString("┘")
				fmtTexts.WriteString(lineSep)

				texts.Reset()
				break NEXT
			}
		case *qxml.StartElement:
			switch e.Name() {
			case "a:p":
				for utils.FindNameIterTo(r, "a:t", "a:p") {
					r.AssignNext(&c_v)
					if !r.Next() {
						break NEXT
					}
					line.WriteString(c_v)
					line.WriteString(space)
					c_v = ""
				}
				if line.Len() > 0 {
					texts.WriteString(space)
					texts.WriteString(line.String())
					texts.WriteString(lineSep)
					if line.Len() > maxLineLen {
						maxLineLen = line.Len()
					}
					line.Reset()
				}
			}
		}
	}

	return fmtTexts, nil
}

// extractImage extracts text content from image by the ocr interface.
//
// Parameters:
//   - rId: the reference id of the image.
//
// Returns:
//   - *strings.Builder: the formatted text of the extracted content.
//   - error: any error that occurred during the extraction process.
func (dp *DocxParser) extractImage(rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	fname, ok := dp.docRelsMap[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := dp.imagesFiles[fname]
	if !ok {
		return nil, types.ErrNonePart
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	text, err := dp.ocr.Run(rc)
	if err != nil {
		return nil, err
	}
	var (
		fmtTexts = new(strings.Builder)
		lineSep  = "\n"
	)
	if dp.drawingsNoFmt {
		fmtTexts.WriteString(text)
		fmtTexts.WriteString(lineSep)
		return fmtTexts, nil
	}

	var (
		newText, maxLineLen = utils.MaxLineLenWithPrefix(text, []byte(" "))
		halfLine            = bytes.Repeat([]byte("─"), max((maxLineLen-5)/2, 0))
	)
	fmtTexts.WriteString("┌")
	fmtTexts.Write(halfLine)
	fmtTexts.WriteString("image")
	fmtTexts.Write(halfLine)
	fmtTexts.WriteString("┐")
	fmtTexts.WriteString(lineSep)

	fmtTexts.WriteString(newText)
	fmtTexts.WriteString(lineSep)

	fmtTexts.WriteString("└")
	fmtTexts.Write(halfLine)
	fmtTexts.WriteString("─────")
	fmtTexts.Write(halfLine)
	fmtTexts.WriteString("┘")
	fmtTexts.WriteString(lineSep)

	return fmtTexts, nil
}
