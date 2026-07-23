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

// Visitor is walking the SZ-based AST.
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
//
// Walk returns an [sx.Object] (and not a *[sx.Pair]) because a [Visitor] may
// return something completely different to a list. For example, it may return
// a user-defined object, which consumes an SZ-based syntax tree.
func Walk(v Visitor, node *sx.Pair, alst *sx.Pair) sx.Object {
	if node == nil {
		return node
	}
	if result, ok := v.VisitBefore(node, alst); ok {
		return result
	}

	if sym := NodeSymbol(node); sym != nil {
		if fn, found := mapWalkChildren[sym]; found {
			node = fn(v, node, alst)
		}
	}
	return v.VisitAfter(node, alst)
}

// VisitorIt is for walking the SZ-based AST, but only for side-effects.
//
// VisitItBefore and VisitItAfter have the same semantic as VisitBefore /
// VisitAfter in the Visitor interface, but no objects are returned.
type VisitorIt interface {
	// VisitItBefore is called before child nodes are traversed. If it returns
	// a true value, the (recursive) walking process ends. If a false value is
	// returned, all child nodes are traversed and after that VisitItAfter is
	// called.
	VisitItBefore(node *sx.Pair, alst *sx.Pair) bool

	// VisitItAfter is called, if the corresponding VisitBefore returned a
	// a false value, and when all child nodes were traversed.
	VisitItAfter(node *sx.Pair, alst *sx.Pair)
}

// WalkIt walks a sx-based AST with the guidance of a Visitor. It will never
// modify the AST. Only the Visitor's side effects count. Therefore, there is
// no need to return something.
func WalkIt(v VisitorIt, node *sx.Pair, alst *sx.Pair) {
	if node == nil {
		return
	}
	if v.VisitItBefore(node, alst) {
		return
	}

	if sym := NodeSymbol(node); sym != nil {
		if fn, found := mapWalkChildrenIt[sym]; found {
			fn(v, node, alst)
		}
	}
	v.VisitItAfter(node, alst)
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

// GetWalkList returns the list of the current node and all subsequent nodes.
func GetWalkList(alst *sx.Pair) *sx.Pair {
	if pair := alst.Assoc(symWalkList); pair != nil {
		return pair.Tail()
	}
	return nil
}

var symWalkList = sx.MakeSymbol("walk-list")

type walkChildrenMap map[*sx.Symbol]func(Visitor, *sx.Pair, *sx.Pair) *sx.Pair
type walkChildrenItMap map[*sx.Symbol]func(VisitorIt, *sx.Pair, *sx.Pair)

var mapWalkChildren walkChildrenMap
var mapWalkChildrenIt walkChildrenItMap

func init() {
	mapWalkChildren = walkChildrenMap{
		SymBlock:         walkChildrenTail,
		SymPara:          walkChildrenTail,
		SymRegionBlock:   walkChildrenRegion,
		SymRegionQuote:   walkChildrenRegion,
		SymRegionVerse:   walkChildrenRegion,
		SymHeading:       walkChildrenTail3,
		SymListOrdered:   walkChildrenTail2,
		SymListUnordered: walkChildrenTail2,
		SymListQuote:     walkChildrenTail2,
		SymListItem:      walkChildrenTail2,
		SymDescription:   walkChildrenTail2,
		SymTerm:          walkChildrenTail2,
		SymDetail:        walkChildrenTail,
		SymEntry:         walkChildrenTail2,
		SymTable:         walkTableChildren,
		SymRow:           walkChildrenTail2,
		SymCell:          walkChildrenTail2,
		SymTransclude:    walkChildrenTail3,
		SymBLOB:          walkBLOBChildren,

		SymInline:       walkChildrenTail,
		SymEndnote:      walkChildrenTail2,
		SymMark:         walkChildrenTail3,
		SymLink:         walkChildrenTail3,
		SymEmbed:        walkEmbedChildren,
		SymCite:         walkChildrenTail3,
		SymFormatDelete: walkChildrenTail2,
		SymFormatEmph:   walkChildrenTail2,
		SymFormatInsert: walkChildrenTail2,
		SymFormatMark:   walkChildrenTail2,
		SymFormatQuote:  walkChildrenTail2,
		SymFormatStrong: walkChildrenTail2,
		SymFormatSpan:   walkChildrenTail2,
		SymFormatSub:    walkChildrenTail2,
		SymFormatSuper:  walkChildrenTail2,
	}
	mapWalkChildrenIt = walkChildrenItMap{
		SymBlock:         walkListIt1,
		SymPara:          walkListIt1,
		SymRegionBlock:   walkChildrenRegionIt,
		SymRegionQuote:   walkChildrenRegionIt,
		SymRegionVerse:   walkChildrenRegionIt,
		SymHeading:       walkListIt3,
		SymListOrdered:   walkListIt2,
		SymListUnordered: walkListIt2,
		SymListQuote:     walkListIt2,
		SymListItem:      walkListIt2,
		SymDescription:   walkListIt2,
		SymTerm:          walkListIt2,
		SymDetail:        walkListIt1,
		SymEntry:         walkListIt2,
		SymTable:         walkListIt2,
		SymRow:           walkListIt2,
		SymCell:          walkListIt2,
		SymTransclude:    walkListIt3,
		SymBLOB:          walkBLOBChildrenIt,

		SymInline:       walkListIt1,
		SymEndnote:      walkListIt2,
		SymMark:         walkListIt3,
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
	return walkChildren(v, node.Tail(), alst, &lb)
}
func walkChildrenTail2(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	newElems := walkChildrenList(v, next.Tail(), alst)
	return newElems.Cons(next.Car()).Cons(sym)
}
func walkChildrenTail3(v Visitor, node *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs, next := next.Car(), next.Tail()
	newElems := walkChildrenList(v, next.Tail(), alst)
	return newElems.Cons(next.Car()).Cons(attrs).Cons(sym)
}

func walkChildrenList(v Visitor, lst *sx.Pair, alst *sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	return walkChildren(v, lst, alst, &lb)
}

func walkChildren(v Visitor, lst *sx.Pair, alst *sx.Pair, lb *sx.ListBuilder) *sx.Pair {
	pos, pairPos, pairList := 0, sx.Cons(symWalkPos, sx.Int64(0)), sx.Cons(symWalkList, sx.Nil())
	alst = alst.Cons(pairPos).Cons(pairList)
	for n := range lst.Pairs() {
		pairList.SetCdr(n)
		obj := Walk(v, n.Head(), alst)
		flattenChildren(lb, obj)
		pos++
		pairPos.SetCdr(sx.Int64(pos))
	}
	return lb.List()
}

func flattenChildren(lb *sx.ListBuilder, obj sx.Object) {
	if sx.IsNil(obj) {
		return
	}
	if pair, isPair := sx.GetPair(obj); isPair {
		if sym := NodeSymbol(pair); SymSpecialSplice.IsEqualSymbol(sym) {
			for child := range pair.Tail().Values() {
				flattenChildren(lb, child)
			}
			return
		}
	}
	lb.Add(obj)
}

// WalkItList will WalkIt for all elements of the list, after skipping the first elements.
func WalkItList(v VisitorIt, lst *sx.Pair, skip int, alst *sx.Pair) {
	for range skip {
		lst = lst.Tail()
	}
	pos, pairPos, pairList := 0, sx.Cons(symWalkPos, sx.Int64(0)), sx.Cons(symWalkList, sx.Nil())
	alst = alst.Cons(pairPos).Cons(pairList)
	for n := range lst.Pairs() {
		pairList.SetCdr(n)
		WalkIt(v, n.Head(), alst)
		pos++
		pairPos.SetCdr(sx.Int64(pos))
	}
}
func walkListIt1(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 1, alst) }
func walkListIt2(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 2, alst) }
func walkListIt3(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 3, alst) }
func walkListIt4(v VisitorIt, node *sx.Pair, alst *sx.Pair) { WalkItList(v, node, 4, alst) }

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

func walkTableChildren(v Visitor, tn *sx.Pair, alst *sx.Pair) *sx.Pair {
	sym, next := tn.Car(), tn.Tail()
	attrs := next.Head()
	next = next.Tail()

	var lb sx.ListBuilder
	lb.AddN(sym, attrs)
	for row := range next.Pairs() {
		lb.Add(Walk(v, row.Head(), alst))
	}
	return lb.List()
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

func walkEmbedChildren(v Visitor, en *sx.Pair, alst *sx.Pair) *sx.Pair {
	attrs, ref, syntax, inlines := GetEmbed(en)
	newInlines := walkChildrenList(v, inlines, alst)
	return MakeEmbed(attrs, ref, syntax, newInlines)
}
