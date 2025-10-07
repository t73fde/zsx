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

import (
	"encoding/base64"

	"t73f.de/r/sx"
)

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

// GetBlock returns all elements of a block node.
func GetBlock(node *sx.Pair) *sx.Pair { return node.Tail() }

// MakeInlineList builds an inline node from a list of inlines.
func MakeInlineList(inlines *sx.Pair) *sx.Pair { return inlines.Cons(SymInline) }

// GetInline returns all elements of an inline node.
func GetInline(node *sx.Pair) *sx.Pair { return node.Tail() }

// MakeParaList builds a paragraph node.
func MakeParaList(inlines *sx.Pair) *sx.Pair { return inlines.Cons(SymPara) }

// GetPara returns all elements of a paragraph node.
func GetPara(node *sx.Pair) *sx.Pair { return node.Tail() }

// MakeList builds a list node.
func MakeList(sym *sx.Symbol, attrs *sx.Pair, items *sx.Pair) *sx.Pair {
	return items.Cons(attrs).Cons(sym)
}

// GetList returns the elements of a list node.
func GetList(node *sx.Pair) (*sx.Symbol, *sx.Pair, *sx.Pair) {
	if sym, isSymbol := sx.GetSymbol(node.Car()); isSymbol {
		attrsNode := node.Tail()
		blocks := attrsNode.Tail()
		return sym, attrsNode.Head(), blocks
	}
	return nil, nil, nil
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

// GetDescription returns the elements of a description list node.
func GetDescription(node *sx.Pair) (*sx.Pair, *sx.Pair) {
	attrsNode := node.Tail()
	termsVals := attrsNode.Tail()
	return attrsNode.Head(), termsVals
}

// GetTable returns the elements of a table.
func GetTable(node *sx.Pair) (*sx.Pair, *sx.Pair, *sx.Pair) {
	attrsNode := node.Tail()
	headerNode := attrsNode.Tail()
	rows := headerNode.Tail()
	return attrsNode.Head(), headerNode.Head(), rows
}

// MakeCell builds a table cell node.
func MakeCell(attrs *sx.Pair, inlines *sx.Pair) *sx.Pair { return inlines.Cons(attrs).Cons(SymCell) }

// GetCell returns all elements of a table cell node.
func GetCell(node *sx.Pair) (*sx.Pair, *sx.Pair) {
	attrsNode := node.Tail()
	inlines := attrsNode.Tail()
	return attrsNode.Head(), inlines
}

// MakeTransclusion builds a transclude node.
func MakeTransclusion(attrs *sx.Pair, ref *sx.Pair, text *sx.Pair) *sx.Pair {
	if text == nil {
		return sx.MakeList(SymTransclude, attrs, ref)
	}
	return text.Cons(ref).Cons(attrs).Cons(SymTransclude)
}

// GetTransclusion returns all elements of a tranclude node.
func GetTransclusion(node *sx.Pair) (*sx.Pair, *sx.Pair, *sx.Pair) {
	attrNode := node.Tail()
	refNode := attrNode.Tail()
	inlines := refNode.Tail()
	return attrNode.Head(), refNode.Head(), inlines
}

// MakeBLOB builds a block BLOB node.
func MakeBLOB(attrs *sx.Pair, syntax string, data []byte, description *sx.Pair) *sx.Pair {
	return description.
		Cons(sx.MakeString(encodeBinary(syntax, data))).
		Cons(sx.MakeString(syntax)).
		Cons(attrs).
		Cons(SymBLOB)
}

// GetBLOB returns all elements of a block BLOB node.
func GetBLOB(node *sx.Pair) (*sx.Pair, string, []byte, *sx.Pair) {
	attrsNode := node.Tail()
	syntaxNode := attrsNode.Tail()
	if syntax, isSyntax := sx.GetString(syntaxNode.Car()); isSyntax {
		contentNode := syntaxNode.Tail()
		if content, isContent := sx.GetString(contentNode.Car()); isContent {
			syn := syntax.GetValue()
			data := decodeBinary(syn, content.GetValue())
			description := contentNode.Tail()
			return attrsNode.Head(), syn, data, description
		}
	}
	return nil, "", nil, nil
}

// MakeBLOBuncode builds a block BLOB node, but assumes already encoded data.
func MakeBLOBuncode(attrs *sx.Pair, syntax string, data string, description *sx.Pair) *sx.Pair {
	return description.Cons(sx.MakeString(data)).Cons(sx.MakeString(syntax)).Cons(attrs).Cons(SymBLOB)
}

// GetBLOBuncode returns all elements of a block BLOB node and does not decode
// the BLOB data.
func GetBLOBuncode(node *sx.Pair) (*sx.Pair, string, string, *sx.Pair) {
	attrsNode := node.Tail()
	syntaxNode := attrsNode.Tail()
	if syntax, isSyntax := sx.GetString(syntaxNode.Car()); isSyntax {
		contentNode := syntaxNode.Tail()
		if data, isData := sx.GetString(contentNode.Car()); isData {
			description := contentNode.Tail()
			return attrsNode.Head(), syntax.GetValue(), data.GetValue(), description
		}
	}
	return nil, "", "", nil
}

// MakeText builds a text node.
func MakeText(text string) *sx.Pair {
	return sx.Cons(SymText, sx.Cons(sx.MakeString(text), sx.Nil()))
}

// GetText returns the element of a text node.
func GetText(node *sx.Pair) string {
	if s, isString := sx.GetString(node.Tail().Car()); isString {
		return s.GetValue()
	}
	return ""
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
func MakeEmbedBLOB(attrs *sx.Pair, syntax string, data []byte, inlines *sx.Pair) *sx.Pair {
	return inlines.
		Cons(sx.MakeString(encodeBinary(syntax, data))).
		Cons(sx.MakeString(syntax)).
		Cons(attrs).
		Cons(SymEmbedBLOB)
}

// GetEmbedBLOB returns all elements of an inline BLOB node.
func GetEmbedBLOB(node *sx.Pair) (*sx.Pair, string, []byte, *sx.Pair) {
	attrsNode := node.Tail()
	syntaxNode := attrsNode.Tail()
	if syntax, isSyntax := sx.GetString(syntaxNode.Car()); isSyntax {
		contentNode := syntaxNode.Tail()
		if content, isContent := sx.GetString(contentNode.Car()); isContent {
			syn := syntax.GetValue()
			data := decodeBinary(syn, content.GetValue())
			inlines := contentNode.Tail()
			return attrsNode.Head(), syn, data, inlines
		}
	}
	return nil, "", nil, nil
}

// MakeEmbedBLOBuncode builds an embedded inline BLOB node, and does not
// encode binary content data.
func MakeEmbedBLOBuncode(attrs *sx.Pair, syntax string, data string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(data)).Cons(sx.MakeString(syntax)).Cons(attrs).Cons(SymEmbedBLOB)
}

// GetEmbedBLOBuncode returns all elements of an inline BLOB node. It does not
// decode the content data.
func GetEmbedBLOBuncode(node *sx.Pair) (*sx.Pair, string, string, *sx.Pair) {
	attrsNode := node.Tail()
	syntaxNode := attrsNode.Tail()
	if syntax, isSyntax := sx.GetString(syntaxNode.Car()); isSyntax {
		contentNode := syntaxNode.Tail()
		if content, isContent := sx.GetString(contentNode.Car()); isContent {
			inlines := contentNode.Tail()
			return attrsNode.Head(), syntax.GetValue(), content.GetValue(), inlines
		}
	}
	return nil, "", "", nil
}

// MakeCite builds a node that specifies a citation.
func MakeCite(attrs *sx.Pair, text string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(text)).Cons(attrs).Cons(SymCite)
}

// GetCite returns alle elements of a cite node.
func GetCite(node *sx.Pair) (*sx.Pair, string, *sx.Pair) {
	attrsNode := node.Tail()
	keyNode := attrsNode.Tail()
	if key, isString := sx.GetString(keyNode.Car()); isString {
		inlines := keyNode.Tail()
		return attrsNode.Head(), key.GetValue(), inlines
	}
	return nil, "", nil
}

// MakeEndnote builds an endnote node.
func MakeEndnote(attrs, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(attrs).Cons(SymEndnote)
}

// GetEndnote returns the elements of an endnote node.
func GetEndnote(node *sx.Pair) (*sx.Pair, *sx.Pair) {
	attrsNode := node.Tail()
	inlines := attrsNode.Tail()
	return attrsNode.Head(), inlines
}

// MakeMark builds a mark note.
func MakeMark(mark string, slug, fragment string, inlines *sx.Pair) *sx.Pair {
	return inlines.Cons(sx.MakeString(fragment)).Cons(sx.MakeString(slug)).Cons(sx.MakeString(mark)).Cons(SymMark)
}

// GetMark returns the elements of a mark node.
func GetMark(node *sx.Pair) (string, string, string, *sx.Pair) {
	markNode := node.Tail()
	if mark, isMark := sx.GetString(markNode.Car()); isMark {
		slugNode := markNode.Tail()
		if slug, isSlug := sx.GetString(slugNode.Car()); isSlug {
			fragmentNode := slugNode.Tail()
			if fragment, isFragment := sx.GetString(fragmentNode.Car()); isFragment {
				inlines := fragmentNode.Tail()
				return mark.GetValue(), slug.GetValue(), fragment.GetValue(), inlines
			}
		}
	}
	return "", "", "", nil
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

func encodeBinary(syntax string, data []byte) string {
	if syntax == SyntaxSVG {
		return string(data)
	}
	return base64.StdEncoding.EncodeToString(data)
}
func decodeBinary(syntax, content string) []byte {
	if syntax == SyntaxSVG {
		return []byte(content)
	}
	if data, err := base64.StdEncoding.DecodeString(content); err == nil {
		return data
	}
	return nil
}
