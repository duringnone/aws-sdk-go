package api

import (
	"bytes"
	"html"
	"regexp"
	"strings"
)

var reNewline = regexp.MustCompile(`\r?\n`)
var reMultiSpace = regexp.MustCompile(`\s+`)
var reComments = regexp.MustCompile(`<!--.*?-->`)
var reExamples = regexp.MustCompile(`<examples?>.+?<\/examples?>`)
var rePara = regexp.MustCompile(`<(?:p|h\d)>(.+?)</(?:p|h\d)>`)
var reLink = regexp.MustCompile(`<a href="(.+?)">(.+?)</a>`)
var reTag = regexp.MustCompile(`<.+?>`)
var reEndNL = regexp.MustCompile(`\n+$`)

// docstring rewrites a string to insert godocs formatting.
func docstring(doc string) string {
	doc = reNewline.ReplaceAllString(doc, "")
	doc = reMultiSpace.ReplaceAllString(doc, " ")
	doc = reComments.ReplaceAllString(doc, "")
	doc = reExamples.ReplaceAllString(doc, "")
	doc = rePara.ReplaceAllString(doc, "$1\n\n")
	doc = reLink.ReplaceAllString(doc, "$2 ($1)")
	doc = reTag.ReplaceAllString(doc, "$1")
	doc = reEndNL.ReplaceAllString(doc, "")
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return "\n"
	}

	doc = html.UnescapeString(doc)
	doc = wrap(doc, 72)

	return commentify(doc)
}

// commentify converts a string to a Go comment
func commentify(doc string) string {
	lines := strings.Split(doc, "\n")
	out := []string{}
	for i, line := range lines {
		if i > 0 && line == "" && lines[i-1] == "" {
			continue
		}
		out = append(out, "// "+line)
	}

	return strings.Join(out, "\n") + "\n"
}

// wrap returns a rewritten version of text to have line breaks
// at approximately length characters. Line breaks will only be
// inserted into whitespace.
func wrap(text string, length int) string {
	var buf bytes.Buffer
	var last rune
	var lastNL bool
	var col int

	for _, c := range text {
		switch c {
		case '\r': // ignore this
			continue // and also don't track `last`
		case '\n': // ignore this too, but reset col
			if col >= length || last == '\n' {
				buf.WriteString("\n\n")
			}
			col = 0
		case ' ', '\t': // opportunity to split
			if col >= length {
				buf.WriteByte('\n')
				col = 0
			} else {
				if !lastNL {
					buf.WriteRune(c)
				}
				col++ // count column
			}
		default:
			buf.WriteRune(c)
			col++
		}
		lastNL = c == '\n'
		last = c
	}
	return buf.String()
}
