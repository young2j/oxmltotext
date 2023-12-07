// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"os"
	"path/filepath"
)

// CreateTempFile creates a temporary file and writes the provided data to it.
//
// Parameters:
//   - data: the data to be written to the file.
//
// Returns:
//   - string: the absolute path of the created file.
//   - error: an error if any occurred during the process.
func CreateTempFile(data []byte) (string, error) {
	tempDir := os.TempDir()
	f, err := os.CreateTemp(tempDir, "filetotext")
	if err != nil {
		return "", err
	}
	defer f.Close()

	filePath, err := filepath.Abs(f.Name())
	if err != nil {
		filePath = filepath.Join(tempDir, f.Name())
	}

	_, err = f.Write(data)
	if err != nil {
		os.Remove(filePath)
		return "", err
	}

	return filePath, nil

}
