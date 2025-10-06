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

package zsx

import "t73f.de/r/sx"

// MakeBlock builds a block node.
func MakeBlock(blocks ...*sx.Pair) *sx.Pair {
	var lb sx.ListBuilder
	lb.Add(SymBlock)
	for _, block := range blocks {
		lb.Add(block)
	}
	return lb.List()
}

// MakeBlockList builds a block node from a list of blocks.
func MakeBlockList(blocks *sx.Pair) *sx.Pair { return blocks.Cons(SymBlock) }

// MakeInlineList builds an inline node from a list of inlines.
func MakeInlineList(inlines *sx.Pair) *sx.Pair { return inlines.Cons(SymInline) }

// MakeParaList builds a paragraph node.
func MakeParaList(inlines *sx.Pair) *sx.Pair { return inlines.Cons(SymPara) }

// MakeList builds a list node.
func MakeList(sym *sx.Symbol, attrs *sx.Pair, items *sx.Pair) *sx.Pair {
	return items.Cons(attrs).Cons(sym)
}

// MakeVerbatim builds a node for verbatim text.
func MakeVerbatim(sym *sx.Symbol, attrs *sx.Pair, content string) *sx.Pair {
	return sx.MakeList(sym, attrs, sx.MakeString(content))
}

// GetVerbatim returns the elements of a verbatim node.
func GetVerbatim(node *sx.Pair) (*sx.Symbol, *sx.Pair, string) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		if s, isString := sx.GetString(attrsNode.Tail().Car()); isString {
			return sym, attrsNode.Head(), s.GetValue()
		}
	}
	return nil, nil, ""
}

// MakeRegion builds a region node.
func MakeRegion(sym *sx.Symbol, attrs *sx.Pair, blocks *sx.Pair, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(blocks).Cons(attrs).Cons(sym)
}

// GetRegion returns the elements of a region node.
func GetRegion(node *sx.Pair) (*sx.Symbol, *sx.Pair, *sx.Pair, *sx.Pair) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		blocksNode := attrsNode.Tail()
		inlines := blocksNode.Tail()
		return sym, attrsNode.Head(), blocksNode.Head(), inlines
	}
	return nil, nil, nil, nil
}

// MakeHeading builds a heading node.
func MakeHeading(level int, attrs, text *sx.Pair, slug, fragment string) *sx.Pair {
	return text.
		Cons(sx.MakeString(fragment)).
		Cons(sx.MakeString(slug)).
		Cons(attrs).
		Cons(sx.Int64(level)).
		Cons(SymHeading)
}

// GetHeading returns the elements of a heading node.
func GetHeading(node *sx.Pair) (int, *sx.Pair, *sx.Pair, string, string) {
	levelNode := node.Tail()
	if levelNum, isNum := sx.GetNumber(levelNode.Car()); isNum {
		if level, isInt := levelNum.(sx.Int64); isInt {
			attrsNode := levelNode.Tail()
			slugNode := attrsNode.Tail()
			if slug, isSlug := sx.GetString(slugNode.Car()); isSlug {
				fragmentNode := slugNode.Tail()
				if fragment, isFragment := sx.GetString(fragmentNode.Car()); isFragment {
					inlines := fragmentNode.Tail()
					return int(level), attrsNode.Head(), inlines, slug.GetValue(), fragment.GetValue()
				}
			}
		}
	}
	return 0, nil, nil, "", ""
}

// MakeThematic builds a node to implement a thematic break.
func MakeThematic(attrs *sx.Pair) *sx.Pair {
	return sx.Cons(SymThematic, sx.Cons(attrs, sx.Nil()))
}

// GetThematic returns the elements of a thematic break node.
func GetThematic(node *sx.Pair) *sx.Pair { return node.Tail().Head() }

// GetList returns the elements of a list node.
func GetList(node *sx.Pair) (*sx.Symbol, *sx.Pair, *sx.Pair) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		blocks := attrsNode.Tail()
		return sym, attrsNode.Head(), blocks
	}
	return nil, nil, nil
}

// MakeCell builds a table cell node.
func MakeCell(attrs *sx.Pair, inlines *sx.Pair) *sx.Pair { return inlines.Cons(attrs).Cons(SymCell) }

// MakeTransclusion builds a transclusion node.
func MakeTransclusion(attrs *sx.Pair, ref *sx.Pair, text *sx.Pair) *sx.Pair {
	if text == nil {
		return sx.MakeList(SymTransclude, attrs, ref)
	}
	return text.Cons(ref).Cons(attrs).Cons(SymTransclude)
}

// MakeBLOB builds a block BLOB node.
func MakeBLOB(attrs, description *sx.Pair, syntax, content string) *sx.Pair {
	return sx.Cons(SymBLOB,
		sx.Cons(attrs,
			sx.Cons(description,
				sx.Cons(sx.MakeString(syntax),
					sx.Cons(sx.MakeString(content), sx.Nil())))))
}

// MakeText builds a text node.
func MakeText(text string) *sx.Pair {
	return sx.Cons(SymText, sx.Cons(sx.MakeString(text), sx.Nil()))
}

// MakeSoft builds a node for a soft line break.
func MakeSoft() *sx.Pair { return sx.Cons(SymSoft, sx.Nil()) }

// MakeHard builds a node for a hard line break.
func MakeHard() *sx.Pair { return sx.Cons(SymHard, sx.Nil()) }

// MakeLink builds a link node.
func MakeLink(attrs *sx.Pair, ref *sx.Pair, text *sx.Pair) *sx.Pair {
	return text.Cons(ref).Cons(attrs).Cons(SymLink)
}

// GetLink returns the elements of a link node.
func GetLink(node *sx.Pair) (*sx.Pair, *sx.Pair, *sx.Pair) {
	attrs := node.Tail()
	ref := attrs.Tail()
	inlines := ref.Tail()
	return attrs.Head(), ref.Head(), inlines
}

// MakeEmbed builds a embed node.
func MakeEmbed(attrs *sx.Pair, ref sx.Object, syntax string, text *sx.Pair) *sx.Pair {
	return text.Cons(sx.MakeString(syntax)).Cons(ref).Cons(attrs).Cons(SymEmbed)
}

// GetEmbed returns the elements of an embed node.
func GetEmbed(node *sx.Pair) (*sx.Pair, *sx.Pair, string, *sx.Pair) {
	attrs := node.Tail()
	ref := attrs.Tail()
	syntax := ref.Tail()
	inlines := syntax.Tail()
	syntaxVal, isString := sx.GetString(syntax.Car())
	if !isString {
		syntaxVal = sx.MakeString("")
	}
	return attrs.Head(), ref.Head(), syntaxVal.GetValue(), inlines
}

// MakeEmbedBLOB builds an embedded inline BLOB node.
func MakeEmbedBLOB(attrs *sx.Pair, syntax string, content string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(content)).Cons(sx.MakeString(syntax)).Cons(attrs).Cons(SymEmbedBLOB)
}

// MakeCite builds a node that specifies a citation.
func MakeCite(attrs *sx.Pair, text string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(text)).Cons(attrs).Cons(SymCite)
}

// MakeEndnote builds an endnote node.
func MakeEndnote(attrs, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(attrs).Cons(SymEndnote)
}

// MakeMark builds a mark note.
func MakeMark(mark string, slug, fragment string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(fragment)).Cons(sx.MakeString(slug)).Cons(sx.MakeString(mark)).Cons(SymMark)
}

// MakeFormat builds an inline formatting node.
func MakeFormat(sym *sx.Symbol, attrs, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(attrs).Cons(sym)
}

// GetFormat returns the elements of a formatting node.
func GetFormat(node *sx.Pair) (*sx.Symbol, *sx.Pair, *sx.Pair) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		inlines := attrsNode.Tail()
		return sym, attrsNode.Head(), inlines
	}
	return nil, nil, nil
}

// MakeLiteral builds a inline node with literal text.
func MakeLiteral(sym *sx.Symbol, attrs *sx.Pair, text string) *sx.Pair {
	return sx.Cons(sym, sx.Cons(attrs, sx.Cons(sx.MakeString(text), sx.Nil())))
}

// GetLiteral returns the elements of a literal node.
func GetLiteral(node *sx.Pair) (*sx.Symbol, *sx.Pair, string) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		if s, isString := sx.GetString(attrsNode.Tail().Car()); isString {
			return sym, attrsNode.Head(), s.GetValue()
		}
	}
	return nil, nil, ""
}

// MakeReference builds a reference node.
func MakeReference(sym *sx.Symbol, val string) *sx.Pair {
	return sx.Cons(sym, sx.Cons(sx.MakeString(val), sx.Nil()))
}

// GetReference returns the reference symbol and value.
func GetReference(ref *sx.Pair) (*sx.Symbol, string) {
	if ref != nil {
		if sym, isSymbol := sx.GetSymbol(ref.Car()); isSymbol {
			val, isString := sx.GetString(ref.Cdr())
			if !isString {
				val, isString = sx.GetString(ref.Tail().Car())
			}
			if isString {
				return sym, val.GetValue()
			}
		}
	}
	return nil, ""
}
