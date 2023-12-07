// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pptxtotext

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/young2j/oxmltotext/types"
	"github.com/young2j/oxmltotext/utils"

	qxml "github.com/dgrr/quickxml"
)

// extractChart extracts the chart text from the pptx file for a given slide index and relationship ID.
//
// Parameters:
//   - i: the index of the slide
//   - rId: the relationship ID of the chart
//
// Returns:
//   - *strings.Builder: the extracted chart text as a strings.Builder
//   - error: any error that occurred during the extraction
func (pp *PptxParser) extractChart(i int, rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	slideRels, ok := pp.slideRelsMap[i]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := slideRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := pp.chartsFiles[fname]
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
				if pp.drawingsNoFmt {
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

// extractChart extracts the diagram text from the pptx file for a given slide index and relationship ID.
//
// Parameters:
//   - i: the index of the slide
//   - rId: the relationship ID of the diagram
//
// Returns:
//   - *strings.Builder: the extracted diagram text as a strings.Builder
//   - error: any error that occurred during the extraction
func (pp *PptxParser) extractDiagram(i int, rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	slideRels, ok := pp.slideRelsMap[i]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := slideRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := pp.diagramsFiles[fname]
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
				if pp.drawingsNoFmt {
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
//   - i: the index of the slide
//   - rId: the reference id of the image.
//
// Returns:
//   - *strings.Builder: the formatted text of the extracted content.
//   - error: any error that occurred during the extraction process.
func (pp *PptxParser) extractImage(i int, rId string) (*strings.Builder, error) {
	if rId == "" {
		return nil, types.ErrEmptyRID
	}

	slideRels, ok := pp.slideRelsMap[i]
	if !ok {
		return nil, types.ErrNonePart
	}

	fname, ok := slideRels[rId]
	if !ok {
		return nil, types.ErrNonePart
	}

	f, ok := pp.imagesFiles[fname]
	if !ok {
		return nil, types.ErrNonePart
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	text, err := pp.ocr.Run(rc)
	if err != nil {
		return nil, err
	}

	var (
		fmtTexts = new(strings.Builder)
		lineSep  = "\n"
	)

	if pp.drawingsNoFmt {
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
