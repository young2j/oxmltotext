// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

const (
	// 微软 Office Word 格式（Microsoft Word 97 - 2004 document）
	EXT_DOC = ".doc"
	CT_DOC  = "application/msword"
	// 微软 Office Word 文档格式
	EXT_DOCX = ".docx"
	CT_DOCX  = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	// 微软 Office Excel 格式（Microsoft Excel 97 - 2004 Workbook）
	EXT_XLS = ".xls"
	CT_XLS  = "application/vnd.ms-excel"
	// 微软 Office Excel 文档格式
	EXT_XLSX = ".xlsx"
	CT_XLSX  = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	// 微软 Office PowerPoint 格式（Microsoft PowerPoint 97 - 2003 演示文稿）
	EXT_PPT = ".ppt"
	CT_PPT  = "application/vnd.ms-powerpoint"
	// 微软 Office PowerPoint 文稿格式
	EXT_PPTX = ".pptx"
	CT_PPTX  = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	// GZ 压缩文件格式
	EXT_GZ   = ".gz"
	CT_GZ    = "application/x-gzip"
	EXT_GZIP = ".gzip"
	CT_GZIP  = CT_GZ
	// ZIP 压缩文件格式
	EXT_ZIP  = ".zip"
	CT_ZIP   = "application/zip"
	EXT_7ZIP = ".7zip"
	CT_7ZIP  = CT_ZIP
	// RAR 压缩文件格式
	EXT_RAR = ".rar"
	CT_RAR  = "application/rar"
	// TAR 压缩文件格式
	EXT_TAR = ".tar"
	CT_TAR  = "application/x-tar"
	EXT_TGZ = ".tgz"
	CT_TGZ  = CT_TAR
	// PDF 是 Portable Document Format 的简称，即便携式文档格式
	EXT_PDF = ".pdf"
	CT_PDF  = "application/pdf"
	// RTF 是指 Rich Text Format，即通常所说的富文本格式
	EXT_RTF = ".rtf"
	CT_RTF  = "application/rtf"
	// GIF 图像格式
	EXT_GIF = ".gif"
	CT_GIF  = "image/gif"
	// JPG(JPEG) 图像格式
	EXT_JPEG = ".jpeg"
	CT_JPEG  = "image/jpeg"
	EXT_JPG  = ".jpg"
	CT_JPG   = CT_JPEG
	// JPG2 图像格式
	EXT_JPG2 = ".jpg2"
	CT_JPG2  = "image/jp2"
	// PNG 图像格式
	EXT_PNG = ".png"
	CT_PNG  = "image/png"
	// TIF(TIFF) 图像格式
	EXT_TIFF = ".tiff"
	CT_TIFF  = "image/tiff"
	EXT_TIF  = ".tif"
	CT_TIF   = CT_TIFF
	// BMP 图像格式（位图格式）
	EXT_BMP = ".bmp"
	CT_BMP  = "image/bmp"
	// SVG 图像格式
	EXT_SVG  = ".svg"
	CT_SVG   = "image/svg+xml"
	EXT_SVGZ = ".svgz"
	CT_SVGZ  = CT_SVG
	// WebP 图像格式
	EXT_WEBP = ".webp"
	CT_WEBP  = "image/webp"
	// ico 图像格式，通常用于浏览器 Favicon 图标
	EXT_ICO = ".ico"
	CT_ICO  = "image/x-icon"
	// 金山 Office 文字排版文件格式
	EXT_WPS = ".wps"
	CT_WPS  = "application/kswps"
	// 金山 Office 表格文件格式
	EXT_ET = ".et"
	CT_ET  = "application/kset"
	// 金山 Office 演示文稿格式
	EXT_DPS = ".dps"
	CT_DPS  = "application/ksdps"
	// Photoshop 源文件格式
	EXT_PSD = ".psd"
	CT_PSD  = "application/x-photoshop"
	// Coreldraw 源文件格式
	EXT_CDR = ".cdr"
	CT_CDR  = "application/x-coreldraw"
	// Adobe Flash 源文件格式
	EXT_SWF = ".swf"
	CT_SWF  = "application/x-shockwave-flash"
	// 普通文本格式
	EXT_TXT = ".txt"
	CT_TXT  = "text/plain"
	EXT_MD  = ".md"
	CT_MD   = CT_TXT
	// Javascript 文件类型
	EXT_JS = ".js"
	CT_JS  = "application/x-javascript"
	// 表示 CSS 样式表
	EXT_CSS = ".css"
	CT_CSS  = "text/css"
	// csv
	EXT_CSV = ".csv"
	CT_CSV  = "text/csv"
	// HTML 文件格式
	EXT_HTML  = ".html"
	CT_HTML   = "text/html"
	EXT_HTM   = ".htm"
	CT_HTM    = CT_HTML
	EXT_SHTML = ".shtml"
	CT_SHTML  = CT_HTML
	// XHTML 文件格式
	EXT_XHTML = ".xhtml"
	CT_XHTML  = "application/xhtml+xml"
	EXT_XHT   = ".xht"
	CT_XHT    = CT_XHTML
	// XML 文件格式
	EXT_XML = ".xml"
	CT_XML  = "text/xml"
	// VCF 文件格式
	EXT_VCF = ".vcf"
	CT_VCF  = "text/x-vcard"
	// PHP 文件格式
	EXT_PHP   = ".php"
	CT_PHP    = "application/x-httpd-php"
	EXT_PHP3  = ".php3"
	CT_PHP3   = CT_PHP
	EXT_PHP4  = ".php4"
	CT_PHP4   = CT_PHP
	EXT_PHTML = ".phtml"
	CT_PHTML  = CT_PHP
	// Java 归档文件格式
	EXT_JAR = ".jar"
	CT_JAR  = "application/java-archive"
	// Android 平台包文件格式
	EXT_APK = ".apk"
	CT_APK  = "application/vnd.android.package-archive"
	// Windows 系统可执行文件格式
	EXT_EXE = ".exe"
	CT_EXE  = "application/octet-stream"
	// PEM 文件格式
	EXT_PEM  = ".pem"
	CT_PEM   = "application/x-x509-user-cert"
	EXT_CERT = ".cert"
	CT_CERT  = CT_PEM
	// mpeg 音频格式
	EXT_MP3 = ".mp3"
	CT_MP3  = "audio/mpeg"
	// mid 音频格式
	EXT_MID  = ".mid"
	CT_MID   = "audio/midi"
	EXT_MIDI = ".midi"
	CT_MIDI  = CT_MID
	// wav 音频格式
	EXT_WAV = ".wav"
	CT_WAV  = "audio/x-wav"
	// m3u 音频格式
	EXT_M3U = ".m3u"
	CT_M3U  = "audio/x-mpegurl"
	// m4a 音频格式
	EXT_M4A = ".m4a"
	CT_M4A  = "audio/x-m4a"
	// ogg 音频格式
	EXT_OGG = ".ogg"
	CT_OGG  = "audio/ogg"
	// Real Audio 音频格式
	EXT_RA = ".ra"
	CT_RA  = "audio/x-realaudio"
	// mp4 视频格式
	EXT_MP4 = ".mp4"
	CT_MP4  = "video/mp4"
	// mpeg 视频格式
	EXT_MPEG = ".mpeg"
	CT_MPEG  = "video/mpeg"
	EXT_MPG  = ".mpg"
	CT_MPG   = CT_MPEG
	EXT_MPE  = ".mpe"
	CT_MPE   = CT_MPEG
	// QuickTime 视频格式
	EXT_MOV = ".mov"
	CT_MOV  = "video/quicktime"
	EXT_QT  = ".qt"
	CT_QT   = CT_MOV
	// m4v 视频格式
	EXT_M4V = ".m4v"
	CT_M4V  = "video/x-m4v"
	// wmv 视频格式（Windows 操作系统上的一种视频格式）
	EXT_WMV = ".wmv"
	CT_WMV  = "video/x-ms-wmv"
	// avi 视频格式
	EXT_AVI = ".avi"
	CT_AVI  = "video/x-msvideo"
	// webm 视频格式
	EXT_WEBM = ".webm"
	CT_WEBM  = "video/webm"
	// 一种基于 flash 技术的视频格式
	EXT_FLV = ".flv"
	CT_FLV  = "video/x-flv"
)

var MIME_MAP = map[string]string{
	EXT_DOC:   CT_DOC,
	EXT_DOCX:  CT_DOCX,
	EXT_XLS:   CT_XLS,
	EXT_XLSX:  CT_XLSX,
	EXT_PPT:   CT_PPT,
	EXT_PPTX:  CT_PPTX,
	EXT_GZ:    CT_GZ,
	EXT_GZIP:  CT_GZIP,
	EXT_ZIP:   CT_ZIP,
	EXT_7ZIP:  CT_7ZIP,
	EXT_RAR:   CT_RAR,
	EXT_TAR:   CT_TAR,
	EXT_TGZ:   CT_TGZ,
	EXT_PDF:   CT_PDF,
	EXT_RTF:   CT_RTF,
	EXT_GIF:   CT_GIF,
	EXT_JPEG:  CT_JPEG,
	EXT_JPG:   CT_JPG,
	EXT_JPG2:  CT_JPG2,
	EXT_PNG:   CT_PNG,
	EXT_TIFF:  CT_TIFF,
	EXT_TIF:   CT_TIF,
	EXT_BMP:   CT_BMP,
	EXT_SVG:   CT_SVG,
	EXT_SVGZ:  CT_SVGZ,
	EXT_WEBP:  CT_WEBP,
	EXT_ICO:   CT_ICO,
	EXT_WPS:   CT_WPS,
	EXT_ET:    CT_ET,
	EXT_DPS:   CT_DPS,
	EXT_PSD:   CT_PSD,
	EXT_CDR:   CT_CDR,
	EXT_SWF:   CT_SWF,
	EXT_TXT:   CT_TXT,
	EXT_MD:    CT_MD,
	EXT_JS:    CT_JS,
	EXT_CSS:   CT_CSS,
	EXT_CSV:   CT_CSV,
	EXT_HTML:  CT_HTML,
	EXT_HTM:   CT_HTM,
	EXT_SHTML: CT_SHTML,
	EXT_XHTML: CT_XHTML,
	EXT_XHT:   CT_XHT,
	EXT_XML:   CT_XML,
	EXT_VCF:   CT_VCF,
	EXT_PHP:   CT_PHP,
	EXT_PHP3:  CT_PHP3,
	EXT_PHP4:  CT_PHP4,
	EXT_PHTML: CT_PHTML,
	EXT_JAR:   CT_JAR,
	EXT_APK:   CT_APK,
	EXT_EXE:   CT_EXE,
	EXT_PEM:   CT_PEM,
	EXT_CERT:  CT_CERT,
	EXT_MP3:   CT_MP3,
	EXT_MID:   CT_MID,
	EXT_MIDI:  CT_MIDI,
	EXT_WAV:   CT_WAV,
	EXT_M3U:   CT_M3U,
	EXT_M4A:   CT_M4A,
	EXT_OGG:   CT_OGG,
	EXT_RA:    CT_RA,
	EXT_MP4:   CT_MP4,
	EXT_MPEG:  CT_MPEG,
	EXT_MPG:   CT_MPG,
	EXT_MPE:   CT_MPE,
	EXT_MOV:   CT_MOV,
	EXT_QT:    CT_QT,
	EXT_M4V:   CT_M4V,
	EXT_WMV:   CT_WMV,
	EXT_AVI:   CT_AVI,
	EXT_WEBM:  CT_WEBM,
	EXT_FLV:   CT_FLV,
}
