// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"unsafe"

	"code.sajari.com/docconv/v2"
	"github.com/google/go-tika/tika"
)

var docPath = "../filesamples/file-sample_100kb.doc"

// not work
func parseDocByDocconv() string {
	f, err := os.Open(docPath)
	if err != nil {
		panic(err)
	}

	res, _, err := docconv.ConvertDoc(f)
	if err != nil {
		panic(err)
	}

	return res
}

func parseDocByAntiword() string {
	output, err := exec.Command("antiword", docPath).Output()
	if err != nil {
		panic(err)
	}
	res := ""
	if len(output) > 0 {
		res = unsafe.String(&output[0], len(output))
	}

	return res
}

func parseDocByTika() string {
	f, err := os.Open(docPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/msword"},
	})
	if err != nil {
		panic(err)
	}

	return res
}
