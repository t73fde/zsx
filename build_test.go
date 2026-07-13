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

package zsx_test

import (
	"testing"

	"t73f.de/r/sx"
	"t73f.de/r/zsx"
)

func TestBuildHeading(t *testing.T) {
	expAttrs := makeSimpleAttrs()
	const expLevel = 2
	expText := zsx.MakeText("text")
	heading := zsx.MakeHeading(expAttrs, expLevel, expText)
	gotAttrs, gotLevel, gotText := zsx.GetHeading(heading)
	if expAttrs != gotAttrs {
		t.Errorf("attrs: exp=%v, got=%v", expAttrs, gotAttrs)
	}
	if expLevel != gotLevel {
		t.Errorf("level: exp=%v, got=%v", expLevel, gotLevel)
	}
	if expText != gotText {
		t.Errorf("text: exp=%v, got=%v", expText, gotText)
	}
}

func TestBuildMark(t *testing.T) {
	expAttrs := makeSimpleAttrs()
	const expMark = "mark"
	expText := zsx.MakeText("text")
	mark := zsx.MakeMark(expAttrs, expMark, expText)
	gotAttrs, gotMark, gotText := zsx.GetMark(mark)
	if expAttrs != gotAttrs {
		t.Errorf("attrs: exp=%v, got=%v", expAttrs, gotAttrs)
	}
	if expMark != gotMark {
		t.Errorf("level: exp=%v, got=%v", expMark, gotMark)
	}
	if expText != gotText {
		t.Errorf("text: exp=%v, got=%v", expText, gotText)
	}
}

func makeSimpleAttrs() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("attr-key"), sx.MakeString("attrs-val")))
}
