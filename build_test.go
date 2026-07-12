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
	const expSlug = "slug"
	const expFrag = "fragment"
	heading := zsx.MakeHeading(expAttrs, expLevel, expText, expSlug, expFrag)
	gotAttrs, gotLevel, gotText, gotSlug, gotFrag := zsx.GetHeading(heading)
	if expAttrs != gotAttrs {
		t.Errorf("attrs: exp=%v, got=%v", expAttrs, gotAttrs)
	}
	if expLevel != gotLevel {
		t.Errorf("level: exp=%v, got=%v", expLevel, gotLevel)
	}
	if expText != gotText {
		t.Errorf("text: exp=%v, got=%v", expText, gotText)
	}
	if expSlug != gotSlug {
		t.Errorf("slug: exp=%v, got=%v", expSlug, gotSlug)
	}
	if expFrag != gotFrag {
		t.Errorf("fragment: exp=%v, got=%v", expFrag, gotFrag)
	}
}

func makeSimpleAttrs() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("attr-key"), sx.MakeString("attrs-val")))
}
