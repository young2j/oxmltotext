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

	"github.com/young2j/oxmltotext/docxtotext"

	"baliance.com/gooxml/document"
	"code.sajari.com/docconv/v2"
	"github.com/fumiama/go-docx"
	"github.com/google/go-tika/tika"
)

var docxPath = "../filesamples/file-sample_100kb.docx"

func parseDocxByGooxml() string {
	doc, err := document.Open(docxPath)
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	for _, p := range doc.Paragraphs() {
		for _, r := range p.Runs() {
			res.WriteString(r.Text())
		}
		res.WriteString("\n")
	}
	return res.String()
}

func parseDocxByGodocx() string {
	f, err := os.Open(docxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	doc, err := docx.Parse(f, info.Size())
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	for _, it := range doc.Document.Body.Items {
		if p, ok := it.(*docx.Paragraph); ok {
			res.WriteString(p.String())
			continue
		}

		if t, ok := it.(*docx.Table); ok {
			res.WriteString(t.String())
		}

	}

	return res.String()
}

func parseDocxByDocconv() string {
	f, err := os.Open(docxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	res, _, err := docconv.ConvertDocx(f)
	if err != nil {
		panic(err)
	}

	return res
}

func parseDocxByOxmlToText() string {
	docxParser, err := docxtotext.Open(docxPath)
	if err != nil {
		panic(err)
	}
	defer docxParser.Close()

	res, err := docxParser.ExtractTexts()
	if err != nil {
		panic(err)
	}

	return res
}

func parseDocxByTika() string {
	f, err := os.Open(docxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	})
	if err != nil {
		panic(err)
	}

	return res
}
