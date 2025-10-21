//-----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of zsx.
//
// zsx is licensed under the latest version of the EUPL (European Union Public
// License). Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
//-----------------------------------------------------------------------------

package zsx_test

import (
	"testing"

	"t73f.de/r/sx"
	"t73f.de/r/zsx"
)

func TestSpliceNodes(t *testing.T) {
	node := zsx.MakeBlock(
		zsx.MakeText("a"),
		sx.MakeList(zsx.SymSpecialSplice,
			sx.MakeList(zsx.SymSpecialSplice, zsx.MakeText("b")),
			zsx.MakeText("c"),
		),
	)
	obj := zsx.Walk(spliceTestVisitor{}, node, nil)
	exp := `(BLOCK (TEXT "a") (TEXT "b") (TEXT "c"))`
	if got := obj.String(); got != exp {
		t.Error(exp)
		t.Error(got)
	}
}

type spliceTestVisitor struct{}

func (spliceTestVisitor) VisitBefore(*sx.Pair, *sx.Pair) (sx.Object, bool) { return sx.Nil(), false }
func (spliceTestVisitor) VisitAfter(node *sx.Pair, _ *sx.Pair) sx.Object   { return node }
