// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import "errors"

var (
	ErrNilZipFile      = errors.New("the input zip file is nil")
	ErrEmptyRID        = errors.New("the rId is empty")
	ErrNonePart        = errors.New("the document part resolves failed or not exists")
	ErrNoSlide         = errors.New("the specified slide is not found")
	ErrNoSheet         = errors.New("the specified sheet is not found")
	ErrNoSharedStrings = errors.New("the sharedStrings.xml file is not found")
	ErrNoDocument      = errors.New("the document.xml file is not found")
	ErrNoComments      = errors.New("the comments.xml file is not found")
	ErrNoEndnotes      = errors.New("the endnotes.xml file is not found")
	ErrNoFootnotes     = errors.New("the footnotes.xml file is not found")
)
