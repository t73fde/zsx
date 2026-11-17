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

package gemtext_test

import (
	"bytes"
	"testing"

	"t73f.de/r/sx"
	"t73f.de/r/zsx"
	"t73f.de/r/zsx/gemtext"
)

func TestParseText(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", `(BLOCK)`},
		{"only blank lines", "\n\n", `(BLOCK (PARA (HARD)))`},
		{"one line", "a", `(BLOCK (PARA (TEXT "a")))`},
		{"lines are separate by hard break", "a\nb\n", `(BLOCK (PARA (TEXT "a") (HARD) (TEXT "b")))`},
		{"CRLF", "a\r\nb\r\n", `(BLOCK (PARA (TEXT "a") (HARD) (TEXT "b")))`},
		{"blank line between", "a\n\nb", `(BLOCK (PARA (TEXT "a") (HARD) (HARD) (TEXT "b")))`},
		{"leading white space is kept", "  a", `(BLOCK (PARA (TEXT "  a")))`},
		{"trailing white space is kept", "a \t", `(BLOCK (PARA (TEXT "a \t")))`},
		{"lone CR is text", "a\rb", `(BLOCK (PARA (TEXT "a\rb")))`},
	}
	runTable(t, testcases)
}

func TestParseHeading(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{"level 1 to 3", "#a\n## b\n###\tc\n",
			`(BLOCK (HEADING () 1 (TEXT "a")) (HEADING () 2 (TEXT "b")) (HEADING () 3 (TEXT "c")))`},
		{"empty heading", "#", `(BLOCK (HEADING () 1))`},
		{"only blanks", "## \t", `(BLOCK (HEADING () 2))`},
		{"trailing blanks removed", "# a  ", `(BLOCK (HEADING () 1 (TEXT "a")))`},
		{"level 4 is level 3 plus #x", "####x", `(BLOCK (HEADING () 3 (TEXT "#x")))`},
		{"leading blank means text", " # a", `(BLOCK (PARA (TEXT " # a")))`},
		{"heading closes list", "* a\n# h", `(BLOCK (UNORDERED () (ITEM () (PARA (TEXT "a")))) (HEADING () 1 (TEXT "h")))`},
	}
	runTable(t, testcases)
}

func TestParseList(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{"items are grouped", "* a\n* b\n",
			`(BLOCK (UNORDERED () (ITEM () (PARA (TEXT "a"))) (ITEM () (PARA (TEXT "b")))))`},
		{"empty item", "* ", `(BLOCK (UNORDERED () (ITEM ())))`},
		{"star without space is text", "*a", `(BLOCK (PARA (TEXT "*a")))`},
		{"star with tab is text", "*\ta", `(BLOCK (PARA (TEXT "*\ta")))`},
		{"text splits list", "* a\nt\n* b",
			`(BLOCK (UNORDERED () (ITEM () (PARA (TEXT "a")))) (PARA (TEXT "t")) (UNORDERED () (ITEM () (PARA (TEXT "b")))))`},
		{"blank line splits list", "* a\n\n* b",
			`(BLOCK (UNORDERED () (ITEM () (PARA (TEXT "a")))) (UNORDERED () (ITEM () (PARA (TEXT "b")))))`},
		{"text after item is kept", "*  a", `(BLOCK (UNORDERED () (ITEM () (PARA (TEXT " a")))))`},
	}
	runTable(t, testcases)
}

func TestParseQuote(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{"lines are grouped", ">a\n> b\n",
			`(BLOCK (REGION-QUOTE () ((PARA (TEXT "a") (HARD) (TEXT "b")))))`},
		{"empty quote line", ">a\n>\n>b",
			`(BLOCK (REGION-QUOTE () ((PARA (TEXT "a") (HARD) (HARD) (TEXT "b")))))`},
		{"only empty quote", ">", `(BLOCK (REGION-QUOTE () ()))`},
		{"text splits quote", ">a\nt\n>b",
			`(BLOCK (REGION-QUOTE () ((PARA (TEXT "a")))) (PARA (TEXT "t")) (REGION-QUOTE () ((PARA (TEXT "b")))))`},
		{"quote closes list", "* a\n>q",
			`(BLOCK (UNORDERED () (ITEM () (PARA (TEXT "a")))) (REGION-QUOTE () ((PARA (TEXT "q")))))`},
		{"list closes quote", ">q\n* a",
			`(BLOCK (REGION-QUOTE () ((PARA (TEXT "q")))) (UNORDERED () (ITEM () (PARA (TEXT "a")))))`},
	}
	runTable(t, testcases)
}

func TestParseLink(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{"with text", "=> gemini://x.org/ Example text  \n",
			`(BLOCK (PARA (LINK () (EXTERNAL "gemini://x.org/") (TEXT "Example text"))))`},
		{"without space after arrow", "=>gemini://x.org/",
			`(BLOCK (PARA (LINK () (EXTERNAL "gemini://x.org/"))))`},
		{"tabs as separator", "=>\tfoo/bar\tText",
			`(BLOCK (PARA (LINK () (EXTERNAL "foo/bar") (TEXT "Text"))))`},
		{"many blanks", "=>  \t u \t  t",
			`(BLOCK (PARA (LINK () (EXTERNAL "u") (TEXT "t"))))`},
		{"relative reference", "=> foo/bar.gmi",
			`(BLOCK (PARA (LINK () (EXTERNAL "foo/bar.gmi"))))`},
		{"no reference", "=>", `(BLOCK (PARA (TEXT "=>")))`},
		{"only blanks", "=>   ", `(BLOCK (PARA (TEXT "=>   ")))`},
		{"leading blank means text", " => x", `(BLOCK (PARA (TEXT " => x")))`},
		{"link closes quote", ">q\n=> u",
			`(BLOCK (REGION-QUOTE () ((PARA (TEXT "q")))) (PARA (LINK () (EXTERNAL "u"))))`},
	}
	runTable(t, testcases)
}

func TestParsePreformatted(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		in   string
		want sx.Object
	}{
		{"simple", "```\na\nb\n```\n", block(verbatim("", "a\nb"))},
		{"with alt text", "```go\nx\n```", block(verbatim("go", "x"))},
		{"alt text is trimmed", "``` \t alt \t\n```", block(verbatim("alt", ""))},
		{"empty block", "```\n```", block(verbatim("", ""))},
		{"lines are not interpreted", "```\n# x\n=> y\n* z\n> q\n```", block(verbatim("", "# x\n=> y\n* z\n> q"))},
		{"white space is kept", "```\n\ta \n  b\n```", block(verbatim("", "\ta \n  b"))},
		{"empty lines are kept", "```\na\n\nb\n```", block(verbatim("", "a\n\nb"))},
		{"leading blank does not close", "```\n ```\n```", block(verbatim("", " ```"))},
		{"closing text is ignored", "```\na\n``` ignored\nb", block(verbatim("", "a"), para("b"))},
		{"unterminated block", "```\na", block(verbatim("", "a"))},
		{"CRLF", "```\r\na\r\n```\r\n", block(verbatim("", "a"))},
		{"closes list", "* a\n```\n```", block(
			sx.MakeList(sx.MakeSymbol("UNORDERED"), sx.Nil(), sx.MakeList(sx.MakeSymbol("ITEM"), sx.Nil(), para("a"))),
			verbatim("", ""))},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := gemtext.Parse([]byte(tc.in))
			if !tc.want.IsEqual(got) {
				t.Errorf("input %q:\nwant: %v\n got: %v", tc.in, tc.want, got)
			}
		})
	}
}

// TestParseDoesNotAliasInput checks that the result is independent of the input slice.
func TestParseDoesNotAliasInput(t *testing.T) {
	t.Parallel()
	src := []byte("text\n# head\n* item\n>quote\n=> url name\n```alt\npre\n```\n")
	want := gemtext.Parse(bytes.Clone(src))
	got := gemtext.Parse(src)
	clear(src)
	if !want.IsEqual(got) {
		t.Errorf("result changed after modifying input:\nwant: %v\n got: %v", want, got)
	}
}

func runTable(t *testing.T, testcases []struct{ name, in, want string }) {
	t.Helper()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := gemtext.Parse([]byte(tc.in)).String(); got != tc.want {
				t.Errorf("input %q:\nwant: %s\n got: %s", tc.in, tc.want, got)
			}
		})
	}
}

// Helpers to construct expected values.

func block(elems ...sx.Object) sx.Object {
	return zsx.MakeBlockList(sx.MakeList(elems...))
}

func para(text string) sx.Object {
	return zsx.MakePara(zsx.MakeText(text))
}

func verbatim(alt, content string) sx.Object {
	var attrs *sx.Pair
	if alt != "" {
		attrs = sx.MakeList(sx.Cons(sx.MakeString("alt"), sx.MakeString(alt)))
	}
	return zsx.MakeVerbatim(zsx.SymVerbatimCode, attrs, content)
}
