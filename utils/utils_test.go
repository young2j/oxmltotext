// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"mime"
	"os"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	p, err := CreateTempFile([]byte("test"))
	if err != nil {
		t.Error(err)
	}
	t.Log(p)

	os.Remove(p)
}

func TestMimeType(t *testing.T) {
	typ := mime.TypeByExtension(".md")
	t.Logf("mimetype:%v", typ)
}

func TestMimeTypeFromURL(t *testing.T) {
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.docx"
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.docx?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.xlsx?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.pptx?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.csv?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"
	u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.md?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"
	// u := "http://www.hbdxzj.org.cn/Uploads/detail/file/20230119/63c891e9e10c8.txt?view=detailV2&ccid=nRsNeN4J&id=020CC&&ajaxserp=0"

	typ, ext := MimeTypeFromURL(u)
	t.Log(typ, ext)
}

func TestMaxLineLength(t *testing.T) {
	l := MaxLineLen("今天的天气\nreally nice!\n")
	t.Log(l)
}

func TestFormatTarget(t *testing.T) {
	want := "ppt/diagrams/data1.xml"
	p := formatTarget("../diagrams/data1.xml", "ppt/")
	if p != want {
		t.Error(p)
	}
	p = formatTarget("../diagrams/./data1.xml", "ppt/")
	if p != want {
		t.Error(p)
	}

	p = formatTarget("./diagrams/data1.xml", "ppt/")
	if p != want {
		t.Error(p)
	}

	p = formatTarget("diagrams/data1.xml", "ppt/")
	if p != want {
		t.Error(p)
	}
}
