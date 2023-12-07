// Copyright (c) 2023 young2j
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT


package main

import (
	"fmt"
	"mime"
)

func mimeType() {
	typ := mime.TypeByExtension(".pdf")
	fmt.Printf("typ: %v\n", typ)
}

func main() {
	mimeType()

	// res := parseDocByAntiword()
	// res := parseDocByTika2()

	// res := parseDocxByGooxml()
	// res := parseDocxByGodocx()
	// res := parseDocxByGooxml()
	// res := parseDocxByXml()
	// res := parseDocxByTika()
	res := parseDocxByOxmlToText()

	// res := parseXlsxByGooxml()
	// res := parseXlsxByTealeg()
	// res := parseXlsxSharedByDocconv()
	// res := parseXlsxByExcelize()
	// res := parseXlsxByQXml()
	// res := parseXlsxByOxmlToText()
	// res := parseXlsxByCmd()
	// res := parseXlsxByTika()

	// res := parsePptxByGooxml()
	// res := parsePptxByDocconv()
	// res := parsePptxByXml()
	// res := parsePptxByOxmlToText()
	// res := parsePptxByTika()

	// res := parsePdfByPdfcpu()
	// res := parsePdfByRPdf()
	// res := parsePdfByUnipdf()
	// res := parsePdfByLpdf()
	// res := parsePdfByDocconv()
	// res := parsePdfByFitz()
	// res := parsePdfByTika()

	// res := parsePptByGooxml() // no
	// res := parsePptByDocconv() // no
	// res := parsePptByTika()
	// res := parsePptByTikaCmd()
	// res := parsePptByUnoconv()

	// res := parseXls1()
	// res := parseXls2()

	// f, _ := os.OpenFile("./res.txt", os.O_CREATE|os.O_WRONLY, 0644)
	// f.WriteString(res)
	// f.Close()

	// f1, _ := os.OpenFile("./res1.txt", os.O_CREATE|os.O_WRONLY, 0644)
	// f1.WriteString(res1)
	// f1.Close()

	// res := JoinByAdd()
	// res := JoinByStringsJoin()
	// res := JoinBySprintf()
	// res := JoinByBytesBuffer()
	// res := JoinByStringsJoin()

	fmt.Printf("res: %v\n", res)
}
