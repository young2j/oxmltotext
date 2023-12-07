// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ppttotext

import (
	"bytes"
	"io"
	"os"

	"github.com/young2j/oxmltotext/utils"
)

// ExtractFromPathByTika extracts text data from a ppt file located at the given path using Tika server.
//
// Parameters:
//   - path: The path of the ppt file to extract data from.
//   - tikaServerURL: The URL of the Tika server to use for extraction.
//
// Returns:
//   - string: The extracted data.
//   - int: The status code of the HTTP response from the Tika server.
//   - error: An error if any occurred during extraction.
func ExtractFromPathByTika(path string, tikaServerURL string) (string, int, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return "", 0, err
	}

	return ExtractFromReaderByTika(f, int(finfo.Size()), tikaServerURL)
}

// ExtractFromURLByTika extracts text data from the specified URL using a Tika server.
//
// Parameters:
//   - u: the URL to extract data from.
//   - tikaServerURL: the URL of the Tika server to use.
//
// Returns:
//   - string: the extracted data.
//   - int: the status code of the URL or Tika Server HTTP response.
//   - error: an error if any occurred.
func ExtractFromURLByTika(u string, tikaServerURL string) (string, int, error) {
	fileResp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(fileResp)
	if err != nil {
		return "", statusCode, err
	}
	r := bytes.NewReader(fileResp.Body)

	return ExtractFromReaderByTika(r, len(fileResp.Body), tikaServerURL)
}

// ExtractFromReaderByTika extracts text data from an io.Reader using the Tika server.
//
// Parameters:
//   - r: The io.Reader from which to extract data.
//   - size: The size of the data in the io.Reader.
//   - tikaServerURL: The URL of the Tika server.
//
// Returns:
//   - string: The extracted data.
//   - int: The status code of HTTP response from the Tika server.
//   - error: An error if any occurred during extraction.
func ExtractFromReaderByTika(r io.Reader, size int, tikaServerURL string) (string, int, error) {
	resp, err := utils.FastPut(
		tikaServerURL,
		utils.WithBodyStream(r, size),
		utils.WithHeaders(map[string]string{
			"X-Tika-OCRskipOcr": "true",
			"Accept":            "text/plain",
			"Content-Type":      "application/vnd.ms-powerpoint",
		}))
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return "", statusCode, err
	}
	res := utils.BytesToString(resp.Body)

	return res, statusCode, nil
}
