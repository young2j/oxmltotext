// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"baliance.com/gooxml/presentation"
	"code.sajari.com/docconv/v2"
	"github.com/framehack/goconv"
	"github.com/gen2brain/go-fitz"
	"github.com/google/go-tika/tika"
)

var pptPath = "../filesamples/file-sample_500kb.ppt"

// not work
func parsePptByGooxml() string {
	ppt, err := presentation.Open(pptPath)
	if err != nil {
		panic(err)
	}
	res := ""
	for _, slide := range ppt.Slides() {
		for _, placeHolder := range slide.PlaceHolders() {
			for _, p := range placeHolder.Paragraphs() {
				for _, r := range p.X().EG_TextRun {
					res += r.R.T
				}
				res += "\n"
			}
		}
		res += "----------------------------\n"
	}
	return res
}

// not work
func parsePptByDocconv() string {
	f, err := os.Open(pptPath)
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

func parsePptByTika() string {
	f, err := os.Open(pptPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/vnd.ms-powerpoint"},
	})
	if err != nil {
		panic(err)
	}

	return res
}

func parsePptByTikaCmd() string {
	output, err := exec.Command("tika", "-t", pptPath).Output()
	if err != nil {
		panic(err)
	}
	res := ""
	if len(output) > 0 {
		res = unsafe.String(&output[0], len(output))
	}

	return res
}

func parsePptByUnoconv() string {
	s := goconv.NewService()
	f, err := os.Open(pptPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 0, stat.Size())
	w := bytes.NewBuffer(buf)

	err = s.Convert(f, w)
	if err != nil {
		panic(err)
	}

	doc, err := fitz.NewFromMemory(w.Bytes())
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	for i := 1; i <= doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			continue
		}
		res.WriteString(text)
		res.WriteString("\n-----------------------\n")
	}

	return res.String()
}
