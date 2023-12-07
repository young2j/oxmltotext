// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"mime"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/young2j/oxmltotext/types"
)

// MimeTypeFromURL returns the MIME type and lowercase file extension from a given URL.
//
// Parameters:
//   - u(string): The URL from which to extract the MIME type and file extension.
//
// Returns:
//   - string: The MIME type extracted from the URL.
//   - string: The file extension extracted from the URL.
func MimeTypeFromURL(u string) (string, string) {
	urL, err := url.Parse(u)
	if err != nil {
		return "", ""
	}
	ext := strings.ToLower(filepath.Ext(urL.Path))
	ct, ok := types.MIME_MAP[ext]
	if ok {
		return ct, ext
	}

	mimeType := mime.TypeByExtension(ext)

	return mimeType, ext
}
