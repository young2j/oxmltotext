// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xlstotext

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/young2j/oxmltotext/utils"
)

// ExtractFromPath extracts text by "xlstotext" cmd from the given xls file path.
//
// Parameters:
// - path: the xls file path to extract text from.
//
// Returns:
// - string: the extracted text.
// - error: any error that occurred during the extraction process.
func ExtractFromPath(path string) (string, error) {
	output, err := exec.Command("xlstotext", path).Output()
	if err != nil {
		return "", err
	}
	res := utils.BytesToString(output)

	return res, nil
}

// ExtractFromURL extracts text data by "xlstotext" cmd from a given xls file URL.
//
// Parameters:
//   - u: a string representing the URL to extract data from.
//
// Returns:
//   - string: the extracted data.
//   - int: the HTTP status code.
//   - error: any error that occurred during the extraction process.
func ExtractFromURL(u string) (string, int, error) {
	resp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return "", statusCode, err
	}

	path, err := utils.CreateTempFile(resp.Body)
	if err != nil {
		return "", statusCode, err
	}
	defer os.Remove(path)

	output, err := exec.Command("xlstotext", path).Output()
	if err != nil {
		return "", statusCode, err
	}
	res := utils.BytesToString(output)

	return res, statusCode, nil
}

// ExtractFromReader extracts text data from an io.Reader.
//
// It reads the data from the provided io.Reader and stores it in a temporary file.
// Then it uses the "xlstotext" command to extract the text from the temporary file.
// The extracted text is returned as a string.
//
// Parameters:
//   - r: An io.Reader from which the data will be read.
//
// Returns:
//   - string: The extracted text.
//   - error: An error if any occurred during the extraction process.
func ExtractFromReader(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	path, err := utils.CreateTempFile(data)
	if err != nil {
		return "", err
	}
	defer os.Remove(path)

	output, err := exec.Command("xlstotext", path).Output()
	if err != nil {
		return "", err
	}
	res := utils.BytesToString(output)

	return res, nil
}

// ExtractFromPathByTika extracts text content from a xls file specified by the given path
// using the Tika server located at the provided URL.
//
// Parameters:
//   - path: The path to the xls file.
//   - tikaServerURL: The URL of the Tika server.
//
// Returns:
//   - string: The extracted text content.
//   - int: the HTTP status code from Tika server.
//   - error: An error if any occurred during the extraction process.
func ExtractFromPathByTika(path string, tikaServerURL string) (string, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return "", 0, err
	}

	return ExtractFromReaderByTika(f, int(finfo.Size()), tikaServerURL)
}

// ExtractFromURLByTika extracts text data from a given xls file URL using the Tika server.
//
// Parameters:
//   - u (string): The xls file URL from which to extract the data.
//   - tikaServerURL (string): The URL of the Tika server.
//
// Returns:
//   - string: The extracted data.
//   - int: The status code of the HTTP response from the URL or Tika server.
//   - error: Any error that occurred during the extraction process.
func ExtractFromURLByTika(u string, tikaServerURL string) (string, int, error) {
	fileResp, err := utils.FastGet(u)
	statusCode := utils.FastStatusCode(fileResp)
	if err != nil {
		return "", statusCode, err
	}
	r := bytes.NewReader(fileResp.Body)

	return ExtractFromReaderByTika(r, len(fileResp.Body), tikaServerURL)
}

// ExtractFromReaderByTika extracts text data from a reader using Tika server.
//
// Parameters:
//   - r: an io.Reader representing the input data.
//   - size: an int representing the size of the input data.
//   - tikaServerURL: a string representing the URL of the Tika server.
//
// Returns:
//   - string: the extracted data.
//   - int: the status code of the Tika server response.
//   - error: an error, if any occurred.
func ExtractFromReaderByTika(r io.Reader, size int, tikaServerURL string) (string, int, error) {
	resp, err := utils.FastPut(
		tikaServerURL,
		utils.WithBodyStream(r, size),
		utils.WithHeaders(map[string]string{
			"X-Tika-OCRskipOcr": "true",
			"Accept":            "text/plain",
			"Content-Type":      "application/vnd.ms-excel",
		}))
	statusCode := utils.FastStatusCode(resp)
	if err != nil {
		return "", statusCode, err
	}
	res := utils.BytesToString(resp.Body)

	return res, statusCode, nil
}
