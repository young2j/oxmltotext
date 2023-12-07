// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import "image"

type Image struct {
	Raw    image.Image
	Name   string
	Format string
}
