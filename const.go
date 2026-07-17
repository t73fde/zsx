//-----------------------------------------------------------------------------
// Copyright (c) 2022-present Detlef Stern
//
// This file is part of zsx.
//
// zsx is licensed under the latest version of the EUPL (European Union Public
// License). Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2022-present Detlef Stern
//-----------------------------------------------------------------------------

package zsx

import "t73f.de/r/sx"

// Various constants for Zettel data. They are technically variables.
var (
	// Symbols for Metanodes
	SymBlock  = sx.MakeSymbol("BLOCK")
	SymInline = sx.MakeSymbol("INLINE")

	// Symbols for Zettel node types.
	SymBLOB            = sx.MakeSymbol("BLOB")
	SymCell            = sx.MakeSymbol("CELL")
	SymCite            = sx.MakeSymbol("CITE")
	SymDescription     = sx.MakeSymbol("DESCRIPTION")
	SymEmbed           = sx.MakeSymbol("EMBED")
	SymEmbedBLOB       = sx.MakeSymbol("EMBED-BLOB")
	SymEndnote         = sx.MakeSymbol("ENDNOTE")
	SymFormatEmph      = sx.MakeSymbol("FORMAT-EMPH")
	SymFormatDelete    = sx.MakeSymbol("FORMAT-DELETE")
	SymFormatInsert    = sx.MakeSymbol("FORMAT-INSERT")
	SymFormatMark      = sx.MakeSymbol("FORMAT-MARK")
	SymFormatQuote     = sx.MakeSymbol("FORMAT-QUOTE")
	SymFormatSpan      = sx.MakeSymbol("FORMAT-SPAN")
	SymFormatSub       = sx.MakeSymbol("FORMAT-SUB")
	SymFormatSuper     = sx.MakeSymbol("FORMAT-SUPER")
	SymFormatStrong    = sx.MakeSymbol("FORMAT-STRONG")
	SymHard            = sx.MakeSymbol("HARD")
	SymHeading         = sx.MakeSymbol("HEADING")
	SymListItem        = sx.MakeSymbol("ITEM")
	SymLink            = sx.MakeSymbol("LINK")
	SymListOrdered     = sx.MakeSymbol("ORDERED")
	SymListUnordered   = sx.MakeSymbol("UNORDERED")
	SymListQuote       = sx.MakeSymbol("QUOTATION")
	SymLiteralCode     = sx.MakeSymbol("LITERAL-CODE")
	SymLiteralComment  = sx.MakeSymbol("LITERAL-COMMENT")
	SymLiteralInput    = sx.MakeSymbol("LITERAL-INPUT")
	SymLiteralMath     = sx.MakeSymbol("LITERAL-MATH")
	SymLiteralOutput   = sx.MakeSymbol("LITERAL-OUTPUT")
	SymMark            = sx.MakeSymbol("MARK")
	SymPara            = sx.MakeSymbol("PARA")
	SymRegionBlock     = sx.MakeSymbol("REGION-BLOCK")
	SymRegionQuote     = sx.MakeSymbol("REGION-QUOTE")
	SymRegionVerse     = sx.MakeSymbol("REGION-VERSE")
	SymRow             = sx.MakeSymbol("ROW")
	SymSoft            = sx.MakeSymbol("SOFT")
	SymTable           = sx.MakeSymbol("TABLE")
	SymText            = sx.MakeSymbol("TEXT")
	SymThematic        = sx.MakeSymbol("THEMATIC")
	SymTransclude      = sx.MakeSymbol("TRANSCLUDE")
	SymVerbatimCode    = sx.MakeSymbol("VERBATIM-CODE")
	SymVerbatimComment = sx.MakeSymbol("VERBATIM-COMMENT")
	SymVerbatimEval    = sx.MakeSymbol("VERBATIM-EVAL")
	SymVerbatimHTML    = sx.MakeSymbol("VERBATIM-HTML")
	SymVerbatimMath    = sx.MakeSymbol("VERBATIM-MATH")
	SymVerbatimZettel  = sx.MakeSymbol("VERBATIM-ZETTEL")

	// Constant symbols for reference states.
	SymRefStateExternal = sx.MakeSymbol("EXTERNAL") // e.g. https://t73f.de/links/software
	SymRefStateHosted   = sx.MakeSymbol("HOSTED")   // e.g. ./foo ../foo /foo /foo/bar
	SymRefStateInvalid  = sx.MakeSymbol("INVALID")  // e.g. :t73f.de/r/zsx
	SymRefStateSelf     = sx.MakeSymbol("SELF")     // e.g. . .#ext #ext

	// Special symbol for later splicing nodes in a list.
	//
	// Spicing may occur, if a Walk Visitor wants to return more than one
	// value. The walking processes will (recursively) flatten a splice node
	// into a (flat) list of nodes.
	SymSpecialSplice = sx.MakeSymbol("*ZSX-SPLICE-NODES*")

	// Special attribute symbols, to be used internal
	SymSpecialID = sx.MakeSymbol("*ZSX-ID*")
)

// Constants for attributes and their values
var (
	SymAttrAlign    = sx.MakeSymbol("align")
	AttrAlignCenter = sx.MakeString("center")
	AttrAlignLeft   = sx.MakeString("left")
	AttrAlignRight  = sx.MakeString("right")
)

// Some additional constants.
const (
	// SyntaxSVG specifies the special syntax "svg", because a BLOB node
	// with this syntax will *not* encode its data.
	SyntaxSVG = "svg"
)
