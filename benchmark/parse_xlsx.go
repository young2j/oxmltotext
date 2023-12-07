// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"archive/zip"
	"context"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unsafe"

	"github.com/young2j/oxmltotext/xlsxtotext"

	"baliance.com/gooxml/spreadsheet"
	"code.sajari.com/docconv/v2"
	"github.com/google/go-tika/tika"
	"github.com/tealeg/xlsx/v3"
	"github.com/xuri/excelize/v2"
)

var xlsxPath = "../filesamples/file-sample_100kb.xlsx"

func parseXlsxByGooxml() string {
	wb, err := spreadsheet.Open(xlsxPath)
	if err != nil {
		panic(err)
	}
	defer wb.Close()
	sheets := wb.Sheets()

	res := new(strings.Builder)
	sheetSep := strings.Repeat("-", 100) + "\n"
	for _, sheet := range sheets {
		for _, row := range sheet.Rows() {
			for _, c := range row.Cells() {
				res.WriteString(c.GetString())
				res.WriteString("\t")
			}
			res.WriteString("\n")
		}

		res.WriteString(sheetSep)
	}
	return res.String()
}

func parseXlsxByTealeg() string {
	f, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		panic(err)
	}

	output, err := f.ToSlice()
	if err != nil {
		panic(err)
	}
	res := new(strings.Builder)
	sheetSep := strings.Repeat("-", 100) + "\n"
	for _, sheet := range output {
		for _, row := range sheet {
			res.WriteString(strings.Join(row, "\t"))
		}
		res.WriteString(sheetSep)
	}
	return res.String()
}

func parseXlsxByExcelize() string {
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sheets := f.GetSheetList()
	res := new(strings.Builder)
	sheetSep := strings.Repeat("-", 100) + "\n"
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		for _, row := range rows {
			res.WriteString(strings.Join(row, "\t"))
		}
		res.WriteString(sheetSep)
	}
	return res.String()
}

var RE_SHEET = regexp.MustCompile("xl/sharedStrings.xml")

// 只能解析共享字符串，不具有参考意义
func parseXlsxSharedByDocconv() string {
	zipReader, err := zip.OpenReader(xlsxPath)
	if err != nil {
		panic(err)
	}
	defer zipReader.Close()

	res := new(strings.Builder)
	sheetSep := strings.Repeat("-", 100) + "\n"
	for _, f := range zipReader.File {
		if !RE_SHEET.MatchString(f.Name) {
			continue
		}
		reader, err := f.Open()
		if err != nil {
			continue
		}
		text, _, err := docconv.ConvertXML(reader)
		if err != nil {
			continue
		}
		res.WriteString(text)
		res.WriteString(sheetSep)

		reader.Close()

	}

	return res.String()
}

func parseXlsxByOxmlToText() string {
	x, err := xlsxtotext.Open(xlsxPath)
	if err != nil {
		panic(err)
	}
	defer x.Close()

	res, err := x.ExtractTexts()
	if err != nil {
		panic(err)
	}

	return res
}

func parseXlsxByCmd() string {
	output, err := exec.Command("xlstotext", xlsxPath).Output()
	if err != nil {
		panic(err)
	}

	res := ""
	if len(output) > 0 {
		res = unsafe.String(&output[0], len(output))
	}

	return res
}

func parseXlsxByTika() string {
	f, err := os.Open(xlsxPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	client := tika.NewClient(nil, "http://localhost:9998")
	res, err := client.ParseWithHeader(ctx, f, http.Header{
		"X-Tika-OCRskipOcr": []string{"true"},
		"Accept":            []string{"text/plain"},
		"Content-Type":      []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	})
	if err != nil {
		panic(err)
	}

	return res
}
