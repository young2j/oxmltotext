// Copyright (c) 2023 young2j
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"archive/zip"
	"path"
	"regexp"
	"strings"

	"github.com/young2j/oxmltotext/types"

	qxml "github.com/dgrr/quickxml"
)

// ParseRelsMap parses a zip file and returns a mapping of relationship IDs to target strings.
//
// Parameters:
//   - f: *zip.File object representing the zip file to parse.
//   - prefix: string prefix used to construct full part name(target string).
//
// Returns:
//   - map[string]string: a mapping of relationship IDs to target strings.
//   - error: an error object indicating any error occurred during parsing.
func ParseRelsMap(f *zip.File, preffix string) (map[string]string, error) {
	m := make(map[string]string)
	if f == nil {
		return nil, types.ErrNilZipFile
	}
	rc, err := f.Open()
	if err != nil {
		return m, err
	}
	defer rc.Close()

	var target = new(strings.Builder)

	r := qxml.NewReader(rc)

	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			switch e.Name() {
			case "Relationship":
				attrs := e.Attrs()
				if attrs.Len() > 0 {
					rIdAttr := attrs.Get("Id")
					targetAttr := attrs.Get("Target")
					t := formatTarget(targetAttr.Value(), preffix)
					target.WriteString(t)
					m[rIdAttr.Value()] = target.String()
					target.Reset()
				}
			}
		}
	}

	return m, nil
}

// formatTarget formats the target string by adding the prefix if it doesn't have it already.
//
// Parameters:
//   - target: the string to be formatted.
//   - prefix: the prefix to be added to the target string.
//
// Returns:
//   - string: the formatted target string.
func formatTarget(target, preffix string) string {
	if strings.HasPrefix(target, preffix) {
		return target
	}
	t := path.Clean(target)
	t = strings.TrimPrefix(t, "../")

	return preffix + t
}

// MatchNameIterTo is a function that matches the name pattern and the to pattern
// iteratively using the given qxml.Reader. It returns true if the name pattern
// is matched and false if the to pattern is matched or if the end of the reader
// is reached.
//
// Parameters:
//   - r: A pointer to a qxml.Reader object
//   - namePattern: The regular expression pattern to match the name
//   - toPattern: The regular expression pattern to match the to
//
// Return:
//   - bool: true if the name pattern is matched, false otherwise
func MatchNameIterTo(r *qxml.Reader, namePattern string, toPattern string) bool {
	re_NAME := regexp.MustCompile(namePattern)
	re_TO := regexp.MustCompile(toPattern)

	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if re_NAME.MatchString(e.Name()) {
				return true
			}
		case *qxml.EndElement:
			if re_TO.MatchString(e.Name()) {
				return false
			}
		}
	}

	return false
}

// FindNameIterTo finds the given name iteratively in the qxml Reader until it reaches the specified end element.
//
// Parameters:
//   - r: a pointer to the qxml Reader.
//   - name: the name to search for in the qxml Reader.
//   - to: the end element name to stop the search.
//
// Returns:
//   - true if the name is found before reaching the end element, false otherwise.
func FindNameIterTo(r *qxml.Reader, name string, to string) bool {
	for r.Next() {
		switch e := r.Element().(type) {
		case *qxml.StartElement:
			if e.Name() == name {
				return true
			}
		case *qxml.EndElement:
			if e.Name() == to {
				return false
			}
		}
	}

	return false
}
