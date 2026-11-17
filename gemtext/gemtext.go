//-----------------------------------------------------------------------------
// Copyright (c) 2026-present Detlef Stern
//
// This file is part of zsx.
//
// zsx is licensed under the latest version of the EUPL (European Union Public
// License). Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2026-present Detlef Stern
//-----------------------------------------------------------------------------

// Package gemtext parses Gemtext (media type text/gemini) into the Sz
// encoding.
//
// The parser follows the Gemtext specification, version 0.24.0
// (gemini://geminiprotocol.net/docs/gemtext-specification.gmi). It works in
// a single pass over its input: every line is handled in isolation, depending
// only on its first characters and on a single mode bit (normal or
// pre-formatted).
//
// # Usage
//
// Parse takes the document as a byte slice and returns an Sz block, (BLOCK
// BlockElement ...). The result does not reference the input slice, so the
// caller may reuse or modify it afterwards. Parsing cannot fail: every input
// is a valid Gemtext document.
//
// # Mapping to Sz
//
//	Gemtext                 Sz
//	text line               (TEXT text), part of a PARA
//	blank line              (nothing), part of a PARA
//	=> url [text]           (LINK () (EXTERNAL url) [(TEXT text)]), part of a PARA
//	# / ## / ### text       (HEADING () level [(TEXT text)])
//	* text (consecutive)    (UNORDERED () (ITEM () [(PARA (TEXT text))]) ...)
//	> text (consecutive)    (REGION-QUOTE () ((PARA (TEXT text) (HARD) ...)))
//	```alt ... ```          (VERBATIM-CODE (("alt" . alt)) content)
//
// Gemtext has no paragraphs: every line is a stand-alone entity, and every
// blank line is vertical blank space. Therefore, consecutive text lines, link
// lines and blank lines are collected into one PARA, and the lines are
// separated by (HARD). Lines are never merged. A blank line contributes no
// element of its own, so a blank line between two text lines results in two
// consecutive (HARD) elements. Separators stand between lines, not after them:
// n lines result in n-1 (HARD) elements. A PARA ends with the first line of
// another type.
//
// Consecutive quote lines form one PARA inside one REGION-QUOTE, separated in
// the same way; an empty quote line is a blank line of that PARA. Consecutive
// list items are grouped into one list, and each item holds its own PARA.
// A list ends with the first line of another type, including a blank line.
//
// # Details
//
// The type of a line is determined by its first characters, without trimming
// leading white space. Hence, " => x" and " ```" are text lines. White space
// means space and tab only; other Unicode white space characters are part of
// the text. In link lines, URL and link text are separated by white space;
// both sides are trimmed. A link line without a URL is handled as a text line.
// URLs are not interpreted: they are always encoded as external references,
// relative references are kept as written.
//
// Heading levels end at three: of more than three leading '#' characters, the
// first three form the heading marker, the remaining ones are part of the
// heading text. The text of a heading, of a quote line, and the alt text
// are trimmed of surrounding white space; the text of text lines and list
// items is kept as it is. A PARA that would contain no element (a single blank
// line) is omitted. An empty list item is an ITEM without a PARA, and an empty
// quote is a REGION-QUOTE without blocks.
//
// Lines inside a pre-formatted block are not interpreted and keep their white
// space. They are joined by a single newline, without a trailing one. Text
// after the closing toggle line is ignored. An unterminated pre-formatted
// block ends at the end of the input and is not an error.
//
// Lines end with LF or CRLF; a CR before the LF is removed. The input is
// expected to be UTF-8, but this is not checked; invalid byte sequences are
// copied to the resulting strings unchanged. A leading UTF-8 byte order mark
// is ignored. Lines end with LF or CRLF; a CR before the LF is removed. The
// input is expected to be UTF-8, but this is not checked; invalid byte
// sequences are copied to the resulting strings unchanged.
package gemtext

import (
	"bytes"

	"t73f.de/r/sx"
	"t73f.de/r/zsx"
)

// Parse parses src as Gemtext and returns the resulting Sz block
// (BLOCK BlockElement ...). The result does not reference src.
func Parse(src []byte) *sx.Pair {
	var p parser
	return p.parse(src)
}

type group uint8

const (
	groupNone group = iota
	groupPara       // text lines, link lines, blank lines
	groupList
	groupQuote
)

type parser struct {
	blocks sx.ListBuilder

	// items collects the elements of the open group: the inline elements
	// (including hard line breaks) of a paragraph or a quote, or the ITEMs of
	// a list. It is empty if no group is open.
	items sx.ListBuilder

	grp     group  // currently open group of consecutive lines
	inPre   bool   // parser is in "pre-formatted mode"
	preAlt  string // alt text of the current pre-formatted block
	preText []byte // content of the current pre-formatted block
	preSeen bool   // at least one line was added to preText
}

var bom = []byte("\xEF\xBB\xBF")

func (p *parser) parse(src []byte) *sx.Pair {
	src = bytes.TrimPrefix(src, bom)
	for len(src) > 0 {
		var line []byte
		if i := bytes.IndexByte(src, '\n'); i >= 0 {
			line, src = src[:i], src[i+1:]
		} else {
			line, src = src, nil
		}
		if n := len(line); n > 0 && line[n-1] == '\r' { // CRLF line end
			line = line[:n-1]
		}
		p.line(line)
	}
	return p.finish()
}

// blank contains the white space characters that are recognized by the parser.
const blank = " \t"

// line handles a single line, without its line end. The line type is
// determined by the first characters of the line, without trimming.
func (p *parser) line(line []byte) {
	if p.inPre {
		p.preLine(line)
		return
	}

	// A line that consists only of space and tab characters is a blank line.
	if len(bytes.Trim(line, blank)) == 0 {
		p.addLine(groupPara, nil)
		return
	}

	switch {
	case line[0] == '`' && len(line) > 2 && line[1] == '`' && line[2] == '`':
		p.closeGroup()
		p.beginPre(line[3:])
	case line[0] == '=' && len(line) > 1 && line[1] == '>':
		p.link(line)
	case line[0] == '*' && len(line) > 1 && line[1] == ' ':
		p.listItem(line[2:])
	case line[0] == '>':
		p.addLine(groupQuote, makeText(bytes.TrimLeft(line[1:], blank)))
	case line[0] == '#':
		p.closeGroup()
		p.heading(line)
	default:
		p.addLine(groupPara, makeText(line))
	}
}

// addLine adds one line to the open paragraph or quote of group g, which is
// opened if necessary. elem may be nil for a line without content. Lines are
// separated by a hard line break.
func (p *parser) addLine(g group, elem *sx.Pair) {
	if p.grp == g {
		p.items.Add(zsx.MakeHard())
	} else {
		p.openGroup(g)
	}
	if elem != nil {
		p.items.Add(elem)
	}
}

func (p *parser) finish() *sx.Pair {
	if p.inPre {
		p.endPre()
	}
	p.closeGroup()
	return zsx.MakeBlockList(p.blocks.List())
}

func (p *parser) emit(block *sx.Pair) { p.blocks.Add(block) }

// makePara returns a paragraph for text, or nil if text is empty.
func makePara(text []byte) *sx.Pair {
	if len(text) == 0 {
		return nil
	}
	return zsx.MakePara(zsx.MakeText(string(text)))
}

// makeText returns a text element for text, or nil if text is empty.
func makeText(text []byte) *sx.Pair {
	if len(text) == 0 {
		return nil
	}
	return zsx.MakeText(string(text))
}

// makeTextList returns a list with a single text element for text, or nil if
// text is empty.
func makeTextList(text []byte) *sx.Pair {
	if len(text) > 0 {
		return sx.Cons(zsx.MakeText(string(text)), nil)
	}
	return nil
}

// ----- heading

// maxHeadingLevel is the highest heading level of Gemtext.
const maxHeadingLevel = 3

func headingLevel(line []byte) int {
	limit := min(maxHeadingLevel, len(line))
	n := 0
	for n < limit && line[n] == '#' {
		n++
	}
	return n
}

func (p *parser) heading(line []byte) {
	level := headingLevel(line)
	text := bytes.Trim(line[level:], blank)
	p.emit(zsx.MakeHeading(nil, level, makeTextList(text)))
}

// ----- link

// link handles a line starting with "=>".
func (p *parser) link(line []byte) {
	url, text := cutLink(line[2:])
	if len(url) == 0 { // no reference: handle as text line
		p.addLine(groupPara, makeText(line))
		return
	}
	ref := zsx.MakeReference(zsx.SymRefStateExternal, string(url))
	p.addLine(groupPara, zsx.MakeLink(nil, ref, makeTextList(text)))
}

// cutLink splits the remainder of a link line into URL and link text.
func cutLink(rest []byte) (url, text []byte) {
	rest = bytes.TrimLeft(rest, blank)
	i := bytes.IndexAny(rest, blank)
	if i < 0 {
		return rest, nil
	}
	return rest[:i], bytes.Trim(rest[i:], blank)
}

// ----- groups: list items and quote lines

func (p *parser) openGroup(g group) {
	if p.grp != g {
		p.closeGroup()
		p.grp = g
	}
}

func (p *parser) listItem(text []byte) {
	p.openGroup(groupList)
	if para := makePara(text); para != nil {
		p.items.Add(zsx.MakeListItem(nil, sx.Cons(para, nil)))
	} else {
		p.items.Add(zsx.MakeListItem(nil, nil))
	}
}

func (p *parser) closeGroup() {
	switch p.grp {
	case groupNone:
		return
	case groupPara:
		if inl := p.items.List(); inl != nil {
			p.emit(sx.Cons(zsx.SymPara, inl))
		}
	case groupList:
		p.emit(zsx.MakeList(zsx.SymListUnordered, nil, p.items.List()))
	case groupQuote:
		var blocks *sx.Pair
		if inl := p.items.List(); inl != nil {
			blocks = sx.Cons(sx.Cons(zsx.SymPara, inl), nil)
		}
		p.emit(zsx.MakeRegion(zsx.SymRegionQuote, nil, blocks, nil))
	}
	p.grp = groupNone
	p.items.Reset()
}

// ----- pre-formatted blocks

func (p *parser) beginPre(alt []byte) {
	p.inPre = true
	p.preAlt = string(bytes.Trim(alt, blank))
	p.preText = p.preText[:0]
	p.preSeen = false
}

var prefixPre = []byte("```")

func (p *parser) preLine(line []byte) {
	if bytes.HasPrefix(line, prefixPre) { // text after the toggle is ignored
		p.endPre()
		return
	}
	if p.preSeen {
		p.preText = append(p.preText, '\n')
	}
	p.preText = append(p.preText, line...)
	p.preSeen = true
}

func (p *parser) endPre() {
	attrs := sx.Nil()
	if p.preAlt != "" {
		attrs = sx.Cons(sx.Cons(sx.MakeString("alt"), sx.MakeString(p.preAlt)), nil)
	}
	p.emit(zsx.MakeVerbatim(zsx.SymVerbatimCode, attrs, string(p.preText)))
	p.inPre = false
}
