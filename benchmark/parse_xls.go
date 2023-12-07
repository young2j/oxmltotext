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

	"github.com/google/go-tika/tika"
)

var xlsPath = "../filesamples/file-sample_100kb.xls"

// 只能按sheet一个一个的提取
// cargo install xls2txt
func parseXls1() string {
	output, err := exec.Command("xls2txt", xlsPath).Output()
	if err != nil {
		panic(err)
	}
	res := ""
	if len(output) > 0 {
		res = unsafe.String(&output[0], len(output))
	}

	return res
}

func parseXls2() string {
	output, err := exec.Command("xlstotext", xlsPath).Output()
	if err != nil {
		panic(err)
	}
	res := ""
	if len(output) > 0 {
		res = unsafe.String(&output[0], len(output))
	}

	return res
}

func parseXlsByTika() string {
	f, err := os.Open(xlsPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/vnd.ms-excel"},
	})
	if err != nil {
		panic(err)
	}

	return res
}
