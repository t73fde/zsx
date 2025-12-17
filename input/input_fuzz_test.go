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

package input_test

import (
	"testing"

	"t73f.de/r/zsx/input"
)

func FuzzNextInput(f *testing.F) {
	f.Fuzz(func(t *testing.T, src []byte) {
		t.Parallel()
		inp := input.NewInput(src)
		for inp.Ch != input.EOS {
			_ = inp.Peek()
			_ = inp.Next()
		}
	})
}
