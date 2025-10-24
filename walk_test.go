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
	"fmt"
	"strconv"
	"testing"

	"t73f.de/r/sx"
	"t73f.de/r/zsx"
)

func TestWalkPosList(t *testing.T) {
	node := zsx.MakeBlock(zsx.MakeText("0"), zsx.MakeText("1"), zsx.MakeText("2"))
	v := posListTestVisitor{}
	_ = zsx.Walk(&v, node, nil)
	if v.err != nil {
		t.Error("Walk:", v.err)
		return
	}
	zsx.WalkIt(&v, node, nil)
	if v.err != nil {
		t.Error("WalkIt:", v.err)
		return
	}
}

type posListTestVisitor struct {
	err error
}

func (v *posListTestVisitor) VisitBefore(node *sx.Pair, alst *sx.Pair) (sx.Object, bool) {
	if v.err == nil {
		v.err = checkPosList(node, alst)
	}
	return sx.Nil(), false
}
func (v *posListTestVisitor) VisitAfter(node *sx.Pair, alst *sx.Pair) sx.Object {
	if v.err == nil {
		v.err = checkPosList(node, alst)
	}
	return node
}
func (v *posListTestVisitor) VisitItBefore(node *sx.Pair, alst *sx.Pair) bool {
	if v.err == nil {
		v.err = checkPosList(node, alst)
	}
	return false
}
func (v *posListTestVisitor) VisitItAfter(node *sx.Pair, alst *sx.Pair) {
	if v.err == nil {
		v.err = checkPosList(node, alst)
	}
}
func checkPosList(node *sx.Pair, alst *sx.Pair) error {
	if zsx.SymText.IsEqual(node.Car()) {
		n := zsx.GetWalkList(alst)
		if n == nil {
			return fmt.Errorf("no parent node for: %v", node)
		}
		if n.Car() != node {
			return fmt.Errorf("parent node %v is not node %v", n.Car(), node)
		}
		text := zsx.GetText(node)
		spos := strconv.Itoa(zsx.GetWalkPos(alst))
		if text != spos {
			return fmt.Errorf("pos of node %v should be %v, but got %v", node, text, spos)
		}
	}
	return nil
}

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
