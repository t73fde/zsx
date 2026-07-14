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

func TestNodeSymbol(t *testing.T) {
	a := sx.MakeSymbol("a")
	testData := []struct {
		name string
		lst  *sx.Pair
		exp  *sx.Symbol
	}{
		{"nil", nil, nil},
		{"list-with-nil", sx.MakeList(sx.Nil()), nil},
		{"list-with-a", sx.MakeList(a), a},
		{"list-with-list", sx.MakeList(sx.MakeList(a)), nil},
	}
	for _, tc := range testData {
		t.Run(tc.name, func(t *testing.T) {
			if got := zsx.NodeSymbol(tc.lst); got != tc.exp {
				t.Errorf("NodeSymbol(%v) should be %v, but got %v", tc.lst, tc.exp, got)
			}
		})
	}
}
