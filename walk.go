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

import "t73f.de/r/sx"

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
	VisitBefore(node *sx.Pair, env *sx.Pair) (sx.Object, bool)

	// VisitAfter is called, if the corresponding VisitBefore returned any
	// object and a false value, and when all child nodes were traversed.
	// The result of VisitAfter is the result of the walking process.
	VisitAfter(node *sx.Pair, env *sx.Pair) sx.Object
}

// Walk a sx-based AST through a Visitor which does not modify the AST.
func Walk(v Visitor, node *sx.Pair, env *sx.Pair) sx.Object {
	if node == nil {
		return node
	}
	if result, ok := v.VisitBefore(node, env); ok {
		return result
	}

	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		if fn, found := walkChildren[sym]; found {
			node = fn(v, node, env)
		}
	}
	return v.VisitAfter(node, env)
}

// WalkBang a sx-based AST through a Visitor which may modify the AST in place.
func WalkBang(v Visitor, node *sx.Pair, env *sx.Pair) sx.Object {
	if node == nil {
		return node
	}
	if result, ok := v.VisitBefore(node, env); ok {
		return result
	}

	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		if fn, found := walkChildrenBang[sym]; found {
			node = fn(v, node, env)
		}
	}
	return v.VisitAfter(node, env)
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
	VisitBefore(node *sx.Pair, env *sx.Pair) bool

	// VisitAfter is called, if the corresponding VisitBefore returned a
	// a false value, and when all child nodes were traversed.
	VisitAfter(node *sx.Pair, env *sx.Pair)
}

// WalkIt walks a sx-based AST with the guidance of a Visitor. It will never
// modify the AST. Only the Visitor's side effects count. Therefore, there is
// no need to return something.
func WalkIt(v VisitorIt, node *sx.Pair, env *sx.Pair) {
	if node == nil {
		return
	}
	if v.VisitBefore(node, env) {
		return
	}

	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		if fn, found := walkChildrenIt[sym]; found {
			fn(v, node, env)
		}
	}
	v.VisitAfter(node, env)
}

// GetWalkPos returns the position of the current element in it parent list.
// It will return -1, if there is no indication about the position.
func GetWalkPos(env *sx.Pair) int {
	if pair := env.Assoc(symWalkPos); pair != nil {
		if i, ok := pair.Cdr().(sx.Int64); ok {
			return int(i)
		}
	}
	return -1
}

var symWalkPos = sx.MakeSymbol("walk-pos")

type walkChildrenMap map[*sx.Symbol]func(Visitor, *sx.Pair, *sx.Pair) *sx.Pair
type walkChildrenItMap map[*sx.Symbol]func(VisitorIt, *sx.Pair, *sx.Pair)

var walkChildren, walkChildrenBang walkChildrenMap
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
	walkChildrenBang = walkChildrenMap{
		SymBlock:         walkChildrenTailBang,
		SymPara:          walkChildrenTailBang,
		SymRegionBlock:   walkChildrenRegionBang,
		SymRegionQuote:   walkChildrenRegionBang,
		SymRegionVerse:   walkChildrenRegionBang,
		SymHeading:       walkHeadingChildrenBang,
		SymListOrdered:   walkListChildrenBang,
		SymListUnordered: walkListChildrenBang,
		SymListQuote:     walkListChildrenBang,
		SymDescription:   walkDescriptionChildrenBang,
		SymTable:         walkTableChildrenBang,
		SymCell:          walkCellChildrenBang,
		SymTransclude:    walkChildrenInlines3Bang,
		SymBLOB:          walkBLOBChildrenBang,

		SymInline:       walkChildrenTailBang,
		SymEndnote:      walkChildrenInlines2Bang,
		SymMark:         walkMarkChildrenBang,
		SymLink:         walkChildrenInlines3Bang,
		SymEmbed:        walkEmbedChildrenBang,
		SymCite:         walkChildrenInlines3Bang,
		SymFormatDelete: walkChildrenInlines2Bang,
		SymFormatEmph:   walkChildrenInlines2Bang,
		SymFormatInsert: walkChildrenInlines2Bang,
		SymFormatMark:   walkChildrenInlines2Bang,
		SymFormatQuote:  walkChildrenInlines2Bang,
		SymFormatStrong: walkChildrenInlines2Bang,
		SymFormatSpan:   walkChildrenInlines2Bang,
		SymFormatSub:    walkChildrenInlines2Bang,
		SymFormatSuper:  walkChildrenInlines2Bang,
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

func walkChildrenTail(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	lb.Add(node.Car())
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range node.Tail().Pairs() {
		obj := Walk(v, n.Head(), env.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if !sx.IsNil(obj) {
			lb.Add(obj)
		}
	}
	return lb.List()
}
func walkChildrenTailBang(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	hasNil := false
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range node.Tail().Pairs() {
		obj := WalkBang(v, n.Head(), env.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if sx.IsNil(obj) {
			hasNil = true
		}
		n.SetCar(obj)
	}
	if !hasNil {
		return node
	}
	for n := node; ; {
		next := n.Tail()
		if next == nil {
			break
		}
		if sx.IsNil(next.Car()) {
			n.SetCdr(next.Cdr())
			continue
		}
		n = next
	}
	return node
}

func walkChildrenList(v Visitor, lst *sx.Pair, env *sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range lst.Pairs() {
		obj := Walk(v, n.Head(), env.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if !sx.IsNil(obj) {
			lb.Add(obj)
		}
	}
	return lb.List()
}
func walkChildrenListBang(v Visitor, lst *sx.Pair, env *sx.Pair) *sx.Pair {
	hasNil := false
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range lst.Pairs() {
		obj := WalkBang(v, n.Head(), env.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
		if sx.IsNil(obj) {
			hasNil = true
		}
		n.SetCar(obj)
	}
	if !hasNil {
		return lst
	}
	var result sx.ListBuilder
	for obj := range lst.Values() {
		if !sx.IsNil(obj) {
			result.Add(obj)
		}
	}
	return result.List()
}

// WalkItList will WalkIt for all elements of the list, after skipping the first elements.
func WalkItList(v VisitorIt, lst *sx.Pair, skip int, env *sx.Pair) {
	for range skip {
		lst = lst.Tail()
	}
	pos, pair := 0, sx.Cons(symWalkPos, sx.Int64(0))
	for n := range lst.Pairs() {
		WalkIt(v, n.Head(), env.Cons(pair))
		pos++
		pair.SetCdr(sx.Int64(pos))
	}
}
func walkListIt1(v VisitorIt, node *sx.Pair, env *sx.Pair) { WalkItList(v, node, 1, env) }
func walkListIt2(v VisitorIt, node *sx.Pair, env *sx.Pair) { WalkItList(v, node, 2, env) }
func walkListIt3(v VisitorIt, node *sx.Pair, env *sx.Pair) { WalkItList(v, node, 3, env) }
func walkListIt4(v VisitorIt, node *sx.Pair, env *sx.Pair) { WalkItList(v, node, 4, env) }
func walkListIt5(v VisitorIt, node *sx.Pair, env *sx.Pair) { WalkItList(v, node, 5, env) }

func walkChildrenRegion(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	next = next.Tail()
	blocks := next.Head()
	newBlocks := walkChildrenList(v, blocks, env)
	inlines := next.Tail()
	newInlines := walkChildrenList(v, inlines, env)
	if blocks != newBlocks || inlines != newInlines {
		var lb sx.ListBuilder
		lb.AddN(sym, attrs, newBlocks)
		lb.ExtendBang(newInlines)
		return lb.List()
	}
	return node
}
func walkChildrenRegionBang(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := node.Car()
	next := node.Tail()
	// attrs := next.Car()
	next = next.Tail()
	next.SetCar(walkChildrenListBang(v, next.Head(), env))
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return node
}
func walkChildrenRegionIt(v VisitorIt, node *sx.Pair, env *sx.Pair) {
	// sym := node.Car()
	next := node.Tail()
	// attrs := next.Car()
	next = next.Tail()
	WalkItList(v, next.Head(), 0, env)
	WalkItList(v, next.Tail(), 0, env)
}

func walkHeadingChildren(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	level := next.Car()
	next = next.Tail()
	attrs := next.Car()
	next = next.Tail()
	slug := next.Car()
	next = next.Tail()
	fragment := next.Car()
	inlines := next.Tail()
	if newInlines := walkChildrenList(v, inlines, env); inlines != newInlines {
		var lb sx.ListBuilder
		lb.AddN(sym, level, attrs, slug, fragment)
		lb.ExtendBang(newInlines)
		return lb.List()
	}
	return node
}
func walkHeadingChildrenBang(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := node.Car()
	next := node.Tail()
	// level := next.Car()
	next = next.Tail()
	// attrs := next.Car()
	next = next.Tail()
	// slug := next.Car()
	next = next.Tail()
	// fragment := next.Car()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return node
}

func walkListChildren(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	items := next.Tail()
	if newItems := walkChildrenList(v, items, env); items != newItems {
		var lb sx.ListBuilder
		lb.AddN(sym, attrs)
		lb.ExtendBang(newItems)
		return lb.List()
	}
	return node
}
func walkListChildrenBang(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := node.Car()
	next := node.Tail()
	// attrs := next.Car()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return node
}

func walkDescriptionChildren(v Visitor, dn *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := dn.Car(), dn.Tail()
	attrs := next.Car()

	var lb sx.ListBuilder
	lb.AddN(sym, attrs)

	for n := next.Tail(); n != nil; n = n.Tail() {
		lb.Add(walkChildrenList(v, n.Head(), env)) // Term
		if n = n.Tail(); n == nil {
			break
		}
		lb.Add(Walk(v, n.Head(), env)) // Values
	}
	return lb.List()
}
func walkDescriptionChildrenBang(v Visitor, dn *sx.Pair, env *sx.Pair) *sx.Pair {
	for n := dn.Tail().Tail(); n != nil; n = n.Tail() {
		n.SetCar(walkChildrenListBang(v, n.Head(), env))
		n = n.Tail()
		if n == nil {
			break
		}
		n.SetCar(WalkBang(v, n.Head(), env))
	}
	return dn
}
func walkDescriptionChildrenIt(v VisitorIt, dn *sx.Pair, env *sx.Pair) {
	for n := dn.Tail().Tail(); n != nil; n = n.Tail() {
		WalkItList(v, n.Head(), 0, env)
		n = n.Tail()
		if n == nil {
			break
		}
		WalkIt(v, n.Head(), env)
	}
}

func walkTableChildren(v Visitor, tn *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := tn.Car(), tn.Tail()

	// TODO: attrs
	var lb sx.ListBuilder
	lb.Add(sym)
	for row := range next.Pairs() {
		lb.Add(walkChildrenList(v, row.Head(), env))
	}
	return lb.List()
}
func walkTableChildrenBang(v Visitor, tn *sx.Pair, env *sx.Pair) *sx.Pair {
	for row := range tn.Tail().Pairs() {
		row.SetCar(walkChildrenListBang(v, row.Head(), env))
	}
	return tn
}
func walkTableChildrenIt(v VisitorIt, tn *sx.Pair, env *sx.Pair) {
	for row := range tn.Tail().Pairs() {
		WalkItList(v, row.Head(), 0, env)
	}
}

func walkCellChildren(v Visitor, cn *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := cn.Car(), cn.Tail()
	attrs, inlines := next.Head(), next.Tail()
	if newInlines := walkChildrenList(v, inlines, env); inlines != newInlines {
		var lb sx.ListBuilder
		lb.AddN(sym, attrs)
		lb.ExtendBang(newInlines)
		return lb.List()
	}
	return cn
}
func walkCellChildrenBang(v Visitor, cn *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := cn.Car()
	next := cn.Tail()
	// attrs := next.Head()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return cn
}

func walkBLOBChildren(v Visitor, bn *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := bn.Car(), bn.Tail()
	attrs := next.Car()
	next = next.Tail()
	inlines := next.Head()
	next = next.Tail()
	s1 := next.Car()
	s2 := next.Tail().Car()

	if newInlines := walkChildrenList(v, inlines, env); newInlines != inlines {
		return sx.MakeList(sym, attrs, newInlines, s1, s2)
	}
	return bn
}
func walkBLOBChildrenBang(v Visitor, bn *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := bn.Car()
	next := bn.Tail()
	// attrs := next.Car()
	next = next.Tail()
	// description := next.Car()
	next.SetCar(walkChildrenListBang(v, next.Head(), env))
	return bn
}
func walkBLOBChildrenIt(v VisitorIt, bn *sx.Pair, env *sx.Pair) {
	// sym := bn.Car()
	next := bn.Tail()
	// attrs := next.Car()
	next = next.Tail()
	// description := next.Car()
	WalkItList(v, next.Head(), 0, env)
}

func walkMarkChildren(v Visitor, mn *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := mn.Car(), mn.Tail()
	mark := next.Car()
	next = next.Tail()
	slug := next.Car()
	next = next.Tail()
	fragment := next.Car()
	inlines := next.Tail()
	if newInlines := walkChildrenList(v, inlines, env); newInlines != inlines {
		var lb sx.ListBuilder
		lb.AddN(sym, mark, slug, fragment)
		lb.ExtendBang(newInlines)
		return lb.List()
	}
	return mn
}
func walkMarkChildrenBang(v Visitor, mn *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := mn.Car()
	next := mn.Tail()
	// mark := next.Car()
	next = next.Tail()
	// slug := next.Car()
	next = next.Tail()
	// fragment := next.Car()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return mn
}

func walkEmbedChildren(v Visitor, en *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := en.Car(), en.Tail()
	attr := next.Car()
	next = next.Tail()
	ref := next.Car()
	next = next.Tail()
	syntax := next.Car()

	inlines := next.Tail()
	if newInlines := walkChildrenList(v, inlines, env); newInlines != inlines {
		var lb sx.ListBuilder
		lb.AddN(sym, attr, ref, syntax)
		lb.ExtendBang(newInlines)
		return lb.List()
	}
	return en
}
func walkEmbedChildrenBang(v Visitor, en *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := en.Car()
	next := en.Tail()
	// attr := next.Car()
	next = next.Tail()
	// ref := next.Car()
	next = next.Tail()
	// syntax := next.Car()

	// text-list := next.Tail()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return en
}

func walkChildrenInlines3(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	next = next.Tail()
	val3 := next.Car()
	elems := next.Tail()
	if newElems := walkChildrenList(v, elems, env); newElems != elems {
		var lb sx.ListBuilder
		lb.AddN(sym, attrs, val3)
		lb.ExtendBang(newElems)
		return lb.List()
	}
	return node
}
func walkChildrenInlines3Bang(v Visitor, ln *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := ln.Car()
	next := ln.Tail()
	// attrs := next.Car()
	next = next.Tail()
	// val3 := next.Car()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return ln
}

func walkChildrenInlines2(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	sym, next := node.Car(), node.Tail()
	attrs := next.Car()
	elems := next.Tail()
	if newElems := walkChildrenList(v, elems, env); newElems != elems {
		var lb sx.ListBuilder
		lb.AddN(sym, attrs)
		lb.ExtendBang(newElems)
		return lb.List()
	}
	return node
}
func walkChildrenInlines2Bang(v Visitor, node *sx.Pair, env *sx.Pair) *sx.Pair {
	// sym := node.Car()
	next := node.Tail() // Attrs
	// attrs := next.Car()
	next.SetCdr(walkChildrenListBang(v, next.Tail(), env))
	return node
}
