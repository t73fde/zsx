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

package gemtext

import (
	"bytes"
	"testing"
)

func TestHeadingLevel(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 0},
		{"#", 1},
		{"##a", 2},
		{"### a", 3},
		{"#### a", 3},
		{"######## a", 3},
		{"# #", 1},
	}
	for _, tc := range testcases {
		if got := headingLevel([]byte(tc.in)); got != tc.want {
			t.Errorf("headingLevel(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestCutLink(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		in       string
		url, txt string
	}{
		{"", "", ""},
		{"  ", "", ""},
		{"u", "u", ""},
		{" u", "u", ""},
		{"u t", "u", "t"},
		{"u \t t u \t", "u", "t u"},
		{"\tu\tt", "u", "t"},
		{"u\u00a0t", "u\u00a0t", ""}, // only space and tab separate
	}
	for _, tc := range testcases {
		url, txt := cutLink([]byte(tc.in))
		if !bytes.Equal(url, []byte(tc.url)) || !bytes.Equal(txt, []byte(tc.txt)) {
			t.Errorf("cutLink(%q) = (%q, %q), want (%q, %q)", tc.in, url, txt, tc.url, tc.txt)
		}
	}
}
