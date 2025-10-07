//-----------------------------------------------------------------------------
// Copyright (c) 2024-present Detlef Stern
//
// This file is part of zsx.
//
// zsx is licensed under the latest version of the EUPL (European Union Public
// License). Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2024-present Detlef Stern
//-----------------------------------------------------------------------------

package zsx

import (
	"t73f.de/r/sx"
)

// Visitor is walking the sx-based AST.
//
// First, VisitBefore is called. If it returned a false value, all child nodes
// are traversed and then VisitAfter is called. If VisitBefore returned a true
// value, no child nodes will be traversed and VisitAfter will not be called.
type Visitor interface {
	// VisitBefore is called before child nodes are traversed. If it returns
	// an object and a true value, the object is used as the result of the
	// (recursive) walking process. If any object and a false value is
	// returned, all child nodes are traversed and after that VisitAfter is
	// called.
	VisitBefore(node *sx.Pair, alst *sx.Pair) (sx.Object, bool)

	// VisitAfter is called, if the corresponding VisitBefore returned any
	// object and a false value, and when all child nodes were traversed.
	// The result of VisitAfter is the result of the walking process.
	VisitAfter(node *sx.Pair, alst *sx.Pair) sx.Object
}

// Walk a sx-based AST through a Visitor.
func Walk(v Visitor, node *sx.Pair, alst *sx.Pair) sx.Object {
	if node == nil {
		return node
	}
	if result, ok := v.VisitBefore(node, alst); ok {
		return result
	}

	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		if fn, found := walkChildren[sym]; found {
			node = fn(v, node, alst)
		}
	}
	return v.VisitAfter(node, alst)
}

// VisitorIt is for walking the sx-based AST, using only side-effects.
//
// VisitBefore and VisitAfter have the same semantic as in the Visitor
// interface, but no objects are returned.
type VisitorIt interface {
	// VisitBefore is called before child nodes are traversed. If it returns
	// a true value, the (recursive) walking process ends. If a false value is
	// returned, all child nodes are traversed and after that VisitAfter is
	// called.
	VisitBefore(node *sx.Pair, alst *sx.Pair) bool

	// VisitAfter is called, if the corresponding VisitBefore returned a
	// a false value, and when all child nodes were traversed.
	VisitAfter(node *sx.Pair, alst *sx.Pair)
}

// WalkIt walks a sx-based AST with the guidance of a Visitor. It will never
// modify the AST. Only the Visitor's side effects count. Therefore, there is
// no need to return something.
func WalkIt(v VisitorIt, node *sx.Pair, alst *sx.Pair) {
	if node == nil {
		return
	}
	if v.VisitBefore(node, alst) {
		return
	}

	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		if fn, found := walkChildrenIt[sym]; found {
			fn(v, node, alst)
		}
	}
	v.VisitAfter(node, alst)
}

// GetWalkPos returns the position of the current element in it parent list.
// It will return -1, if there is no indication about the position.
func GetWalkPos(alst *sx.Pair) int {
	if pair := alst.Assoc(symWalkPos); pair != nil {
		if i, ok := pair.Cdr().(sx.Int64); ok {
			return int(i)
		}
	}
	return -1
}

var symWalkPos = sx.MakeSymbol("walk-pos")

type walkChildrenMap map[*sx.Symbol]func(Visitor, *sx.Pair, *sx.Pair) *sx.Pair
type walkChildrenItMap map[*sx.Symbol]func(VisitorIt, *sx.Pair, *sx.Pair)

var walkChildren walkChildrenMap
var walkChildrenIt walkChildrenItMap

func init() {
	walkChildren = walkChildrenMap{
		SymBlock:         walkChildrenTail,
		SymPara:          walkChildrenTail,
		SymRegionBlock:   walkChildrenRegion,
		SymRegionQuote:   walkChildrenRegion,
		SymRegionVerse:   walkChildrenRegion,
		SymHeading:       walkHeadingChildren,
		SymListOrdered:   walkListChildren,
		SymListUnordered: walkListChildren,
		SymListQuote:     walkListChildren,
		SymDescription:   walkDescriptionChildren,
		SymTable:         walkTableChildren,
		SymCell:          walkCellChildren,
		SymTransclude:    walkChildrenInlines3,
		SymBLOB:          walkBLOBChildren,

		SymInline:       walkChildrenTail,
		SymEndnote:      walkChildrenInlines2,
		SymMark:         walkMarkChildren,
		SymLink:         walkChildrenInlines3,
		SymEmbed:        walkEmbedChildren,
		SymCite:         walkChildrenInlines3,
		SymFormatDelete: walkChildrenInlines2,
		SymFormatEmph:   walkChildrenInlines2,
		SymFormatInsert: walkChildrenInlines2,
		SymFormatMark:   walkChildrenInlines2,
		SymFormatQuote:  walkChildrenInlines2,
		SymFormatStrong: walkChildrenInlines2,
		SymFormatSpan:   walkChildrenInlines2,
		SymFormatSub:    walkChildrenInlines2,
		SymFormatSuper:  walkChildrenInlines2,
	}
	walkChildrenIt = walkChildrenItMap{
		SymBlock:         walkListIt1,
		SymPara:          walkListIt1,
		SymRegionBlock:   walkChildrenRegionIt,
		SymRegionQuote:   walkChildrenRegionIt,
		SymRegionVerse:   walkChildrenRegionIt,
		SymHeading:       walkListIt5,
		SymListOrdered:   walkListIt2,
		SymListUnordered: walkListIt2,
		SymListQuote:     walkListIt2,
		SymDescription:   walkDescriptionChildrenIt,
		SymTable:         walkTableChildrenIt,
		SymCell:          walkListIt2,
		SymTransclude:    walkListIt3,
		SymBLOB:          walkBLOBChildrenIt,

		SymInline:       walkListIt1,
		SymEndnote:      walkListIt2,
		SymMark:         walkListIt4,
		SymLink:         walkListIt3,
		SymEmbed:        walkListIt4,
		SymCite:         walkListIt3,
		SymFormatDelete: walkListIt2,
		SymFormatEmph:   walkListIt2,
		SymFormatInsert: walkListIt2,
		SymFormatMark:   walkListIt2,
		SymFormatQuote:  walkListIt2,
		SymFormatStrong: walkListIt2,
		SymFormatSpan:   walkListIt2,
		SymFormatSub:    walkListIt2,
		SymFormatSuper:  walkListIt2,
	}
}

func walkChildrenTail(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	lb.Add(node.Car())
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range node.Tail().Pairs() {
		obj := Walk(v, n.Head(), alst.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if !sx.IsNil(obj) {
			lb.Add(obj)
		}
	}
	return lb.List()
}

func walkChildrenList(v Visitor, lst *sx.Pair, alst *sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range lst.Pairs() {
		obj := Walk(v, n.Head(), alst.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if !sx.IsNil(obj) {
			lb.Add(obj)
		}
	}
	return lb.List()
}

// WalkItList will WalkIt for all elements of the list, after skipping the first elements.
func WalkItList(v VisitorIt, lst *sx.Pair, skip int, alst *sx.Pair) {
	for range skip {
		lst = lst.Tail()
	}
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range lst.Pairs() {
		WalkIt(v, n.Head(), alst.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
	}
}
func walkListIt1(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 1, alst) }
func walkListIt2(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 2, alst) }
func walkListIt3(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 3, alst) }
func walkListIt4(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 4, alst) }
func walkListIt5(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 5, alst) }

func walkChildrenRegion(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, attrs, blocks, inlines := GetRegion(node)
	newBlocks := walkChildrenList(v, blocks, alst)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeRegion(sym, attrs, newBlocks, newInlines)
}
func walkChildrenRegionIt(v VisitorIt, node *sx.Pair, alst *sx.Pair) {
	_, _, blocks, inlines := GetRegion(node)
	WalkItList(v, blocks, 0, alst)
	WalkItList(v, inlines, 0, alst)
}

func walkHeadingChildren(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	level, attrs, inlines, slug, fragment := GetHeading(node)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeHeading(level, attrs, newInlines, slug, fragment)
}

func walkListChildren(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, attrs, items := GetList(node)
	newItems := walkChildrenList(v, items, alst)
	return MakeList(sym, attrs, newItems)
}

func walkDescriptionChildren(v Visitor, dn *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := dn.Car(), dn.Tail()
	attrs := next.Car()

	var lb sx.ListBuilder
	lb.AddN(sym, attrs)

	for n := next.Tail(); n != nil; n = n.Tail() {
		lb.Add(walkChildrenList(v, n.Head(), alst)) // Term
		if n = n.Tail(); n == nil {
			break
		}
		lb.Add(Walk(v, n.Head(), alst)) // Values
	}
	return lb.List()
}
func walkDescriptionChildrenIt(v VisitorIt, dn *sx.Pair, alst *sx.Pair) {
	for n := dn.Tail().Tail(); n != nil; n = n.Tail() {
		WalkItList(v, n.Head(), 0, alst)
		n = n.Tail()
		if n == nil {
			break
		}
		WalkIt(v, n.Head(), alst)
	}
}

func walkTableChildren(v Visitor, tn *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := tn.Car(), tn.Tail()
	attrs := next.Head()
	next = next.Tail()

	var lb sx.ListBuilder
	lb.AddN(sym, attrs)
	for row := range next.Pairs() {
		lb.Add(walkChildrenList(v, row.Head(), alst))
	}
	return lb.List()
}
func walkTableChildrenIt(v VisitorIt, tn *sx.Pair, alst *sx.Pair) {
	for row := range tn.Tail().Tail().Pairs() {
		WalkItList(v, row.Head(), 0, alst)
	}
}

func walkCellChildren(v Visitor, cn *sx.Pair, alst *sx.Pair) *sx.Pair {
	attrs, inlines := GetCell(cn)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeCell(attrs, newInlines)
}

func walkBLOBChildren(v Visitor, bn *sx.Pair, alst *sx.Pair) *sx.Pair {
	attrs, syntax, content, inlines := GetBLOBuncode(bn)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeBLOBuncode(attrs, syntax, content, newInlines)
}
func walkBLOBChildrenIt(v VisitorIt, bn *sx.Pair, alst *sx.Pair) {
	_, _, _, inlines := GetBLOBuncode(bn)
	WalkItList(v, inlines, 0, alst)
}

func walkMarkChildren(v Visitor, mn *sx.Pair, alst *sx.Pair) *sx.Pair {
	mark, slug, fragment, inlines := GetMark(mn)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeMark(mark, slug, fragment, newInlines)
}

func walkEmbedChildren(v Visitor, en *sx.Pair, alst *sx.Pair) *sx.Pair {
	attrs, ref, syntax, inlines := GetEmbed(en)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeEmbed(attrs, ref, syntax, newInlines)
}

func walkChildrenInlines3(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	next = next.Tail()
	val3 := next.Car()
	elems := next.Tail()
	newElems := walkChildrenList(v, elems, alst)
	var lb sx.ListBuilder
	lb.AddN(sym, attrs, val3)
	lb.ExtendBang(newElems)
	return lb.List()
}

func walkChildrenInlines2(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	elems := next.Tail()
	newElems := walkChildrenList(v, elems, alst)
	var lb sx.ListBuilder
	lb.AddN(sym, attrs)
	lb.ExtendBang(newElems)
	return lb.List()
}
