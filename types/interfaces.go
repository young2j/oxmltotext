// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import "io"

type OCR interface {
	Run(r io.Reader) (string, error)
	Close() error
}
