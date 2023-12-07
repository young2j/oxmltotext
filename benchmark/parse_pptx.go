// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"net/http"
	"os"
	"strings"

	"github.com/young2j/oxmltotext/pptxtotext"

	"baliance.com/gooxml/presentation"
	"code.sajari.com/docconv/v2"
	"github.com/google/go-tika/tika"
)

var pptxPath = "../filesamples/file-sample_500kb.pptx"

func parsePptxByGooxml() string {
	ppt, err := presentation.Open(pptxPath)
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	for _, slide := range ppt.Slides() {
		for _, placeHolder := range slide.PlaceHolders() {
			for _, p := range placeHolder.Paragraphs() {
				for _, r := range p.X().EG_TextRun {
					res.WriteString(r.R.T)
				}
				res.WriteString("\n")
			}
		}
		res.WriteString("----------------------------\n")
	}
	return res.String()
}

func parsePptxByDocconv() string {
	f, err := os.Open(pptxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	res, _, err := docconv.ConvertPptx(f)
	if err != nil {
		panic(err)
	}

	return res
}

func parsePptxByOxmlToText() string {
	ppt, err := pptxtotext.Open(pptxPath)
	if err != nil {
		panic(err)
	}
	defer ppt.Close()

	res, err := ppt.ExtractTexts()
	if err != nil {
		panic(err)
	}

	return res
}

func parsePptxByTika() string {
	f, err := os.Open(pptxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	})
	if err != nil {
		panic(err)
	}
	return res
}
