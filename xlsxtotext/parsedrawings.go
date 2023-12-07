// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlsxtotext

import (
	"bytes"
	"log"
	"regexp"
	"strings"

	"github.com/young2j/oxmltotext/types"
	"github.com/young2j/oxmltotext/utils"

	qxml "github.com/dgrr/quickxml"
)

// extractDrawings extracts drawings text from the specified sheet index and relationship ID.
//
// The function iterates through the XML elements of the drawing part reader and
// extracts the drawings(charts, images, and diagrams) if the corresponding flags are set.
//
// Parameters:
//   - i: the index of the sheet
//   - rId: the relationship ID of the drawing
//
// Returns:
//   - *strings.Builder: the extracted drawings text as a strings.Builder
func (xp *XlsxParser) extractDrawings(i int, rId string) *strings.Builder {
	if rId == "" {
		return nil
	}

	sheetRels, ok := xp.sheetRelsMap[i]
	if !ok {
		return nil
	}

	drawingName, ok := sheetRels[rId]
	if !ok {
		return nil
	}

	f, ok := xp.drawingsFile[drawingName]
	if !ok {
		return nil
	}
	rc, err := f.Open()
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rc.Close()

	r := qxml.NewReader(rc)
	var texts = new(strings.Builder)

	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.EndElement:

		case *qxml.StartElement:
			switch {
			case e.Name() == "c:chart" && xp.parseCharts:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:id")
					charts, err := xp.extractChart(drawingName, rIdKV.Value())
					xp.logWarn(err)
					if charts != nil {
						texts.WriteString(charts.String())
					}
				}

			case e.Name() == "dgm:relIds" && xp.parseDiagrams:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:dm")
					diagrams, err := xp.extractDiagram(drawingName, rIdKV.Value())
					xp.logWarn(err)
					if diagrams != nil {
						texts.WriteString(diagrams.String())
					}
				}

			case e.Name() == "a:blip" && xp.parseImages:
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdKV := attrs.Get("r:embed")
					images, err := xp.extractImage(drawingName, rIdKV.Value())
					xp.logWarn(err)
					if images != nil {
						texts.WriteString(images.String())
					}
				}

			}
		}
	}

	return texts
}

// extractChart extracts the chart text from the xlsx file for a given drawing part and relationship ID.
//
// Parameters:
//   - drawingName: the name of drawing part
//   - rId: the relationship ID of the chart
//
// Returns:
//   - *strings.Builder: the extracted chart text as a strings.Builder
//   - error: any error that occurred during the extraction
func (xp *XlsxParser) extractChart(drawingName string, rId string) (*strings.Builder, error) {
	drawingRels, ok := xp.drawingRelsMap[drawingName]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := drawingRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := xp.chartsFiles[fname]
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
				if xp.drawingsNoFmt {
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

// extractChart extracts the diagram text from the xlsx file for a given drawing part and relationship ID.
//
// Parameters:
//   - drawingName: the name of drawing part
//   - rId: the relationship ID of the diagram
//
// Returns:
//   - *strings.Builder: the extracted diagram text as a strings.Builder
//   - error: any error that occurred during the extraction
func (xp *XlsxParser) extractDiagram(drawingName string, rId string) (*strings.Builder, error) {
	drawingRels, ok := xp.drawingRelsMap[drawingName]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := drawingRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := xp.diagramsFiles[fname]
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
				if xp.drawingsNoFmt {
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
//   - drawingName: the name of drawing part
//   - rId: the reference id of the image.
//
// Returns:
//   - *strings.Builder: the formatted text of the extracted content.
//   - error: any error that occurred during the extraction process.
func (xp *XlsxParser) extractImage(drawingName string, rId string) (*strings.Builder, error) {
	drawingRels, ok := xp.drawingRelsMap[drawingName]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := drawingRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := xp.imagesFiles[fname]
	if !ok {
		return nil, types.ErrNonePart
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	text, err := xp.ocr.Run(rc)
	if err != nil {
		return nil, err
	}
	var (
		fmtTexts = new(strings.Builder)
		lineSep  = "\n"
	)

	if xp.drawingsNoFmt {
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
