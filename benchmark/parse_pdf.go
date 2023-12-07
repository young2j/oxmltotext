// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"code.sajari.com/docconv/v2"
	"github.com/google/go-tika/tika"
	lpdf "github.com/ledongthuc/pdf"
	"github.com/moolekkari/unipdf/core"
	pdfModel "github.com/moolekkari/unipdf/model"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
	rpdf "github.com/rsc/pdf"

	"github.com/gen2brain/go-fitz"
)

var pdfPath = "../filesamples/file-sample_500kb.pdf"

// not work
func parsePdfByPdfcpu() string {
	f, err := os.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	model.ConfigPath = "disable"
	conf := model.NewDefaultConfiguration()
	conf.CheckFileNameExt = false
	conf.WriteObjectStream = false
	conf.WriteXRefStream = false
	conf.Cmd = model.EXTRACTCONTENT

	ctx, err := api.ReadContext(f, conf)
	if err != nil {
		panic(err)
	}

	if err = validate.XRefTable(ctx.XRefTable); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if err = pdfcpu.OptimizeXRefTable(ctx); err != nil {
		panic(err)
	}

	if err := ctx.EnsurePageCount(); err != nil {
		panic(err)
	}

	pages := types.IntSet{}
	for i := 1; i <= ctx.PageCount; i++ {
		pages[i] = true
	}
	res := ""
	for p, v := range pages {
		if !v {
			continue
		}

		r, err := pdfcpu.ExtractPageContent(ctx, p)
		if err != nil {
			panic(err)
		}
		if r == nil {
			continue
		}
		pageContent, err := io.ReadAll(r)
		if err != nil {
			panic(err)
		}

		res += string(pageContent)
		res += "---------------------------------\n"
	}

	return res
}

// not work
func parsePdfByRPdf() string {
	r, err := rpdf.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	res := ""
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		texts := p.Content().Text
		for _, t := range texts {
			res += t.S + "\n"
		}
		res += "---------------------------------\n"
	}

	return res
}

// not work
func parsePdfByUnipdf() string {
	f, err := os.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	parser, err := core.NewParser(f)
	if err != nil {
		panic(err)
	}

	dict, err := parser.ParseDict()
	if err != nil {
		panic(err)
	}
	pageRes, err := pdfModel.NewPdfPageResourcesFromDict(dict)
	if err != nil {
		panic(err)
	}
	s := pageRes.XObject.String()
	fmt.Printf("s: %v\n", s)

	return ""
}

func parsePdfByLpdf() string {
	_, r, err := lpdf.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		texts := p.Content().Text
		for _, t := range texts {
			res.WriteString(t.S)
		}
		res.WriteString("\n---------------------------------\n")
	}

	return res.String()
}

// 依赖pdftotext命令行
// brew install poppler
func parsePdfByDocconv() string {
	f, err := os.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	res, _, err := docconv.ConvertPDF(f)
	if err != nil {
		panic(err)
	}

	return res
}

func parsePdfByFitz() string {
	f, err := os.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	doc, err := fitz.NewFromReader(f)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	res := new(strings.Builder)
	for i := 0; i < doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			continue
		}
		res.WriteString(text)
		res.WriteString("-----------------------\n")
	}

	return res.String()
}

func parsePdfByTika() string {
	f, err := os.Open(pdfPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/pdf"},
	})
	if err != nil {
		panic(err)
	}

	return res
}
