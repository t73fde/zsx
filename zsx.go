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

// Package zsx contains zettel data handling as sx expressions.
package zsx

import "t73f.de/r/sx"

// GoValue returns the string value of the sx.Object suitable for Go processing.
func GoValue(obj sx.Object) string {
	switch o := obj.(type) {
	case sx.String:
		return o.GetValue()
	case *sx.Symbol:
		return o.GetValue()
	}
	return obj.String()
}
