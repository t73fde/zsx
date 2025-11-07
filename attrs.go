//-----------------------------------------------------------------------------
// Copyright (c) 2020-present Detlef Stern
//
// This file is part of zsx.
//
// zsx is licensed under the latest version of the EUPL (European Union Public
// License). Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2020-present Detlef Stern
//-----------------------------------------------------------------------------

package zsx

import (
	"maps"
	"slices"
	"strings"

	"t73f.de/r/sx"
)

// Attributes store additional information about some node types.
type Attributes map[string]string

// IsEmpty returns true if there are no attributes.
func (a Attributes) IsEmpty() bool { return len(a) == 0 }

// DefaultAttribute is the value of the key of the default attribute
const DefaultAttribute = "-"

// HasDefault returns true, if the default attribute "-" has been set.
func (a Attributes) HasDefault() bool {
	if a != nil {
		_, ok := a[DefaultAttribute]
		return ok
	}
	return false
}

// RemoveDefault removes the default attribute
func (a Attributes) RemoveDefault() Attributes {
	if a != nil {
		a.Remove(DefaultAttribute)
	}
	return a
}

// Keys returns the sorted list of keys.
func (a Attributes) Keys() []string { return slices.Sorted(maps.Keys(a)) }

// Get returns the attribute value of the given key and a succes value.
func (a Attributes) Get(key string) (string, bool) {
	if a != nil {
		value, ok := a[key]
		return value, ok
	}
	return "", false
}

// Clone returns a duplicate of the attribute.
func (a Attributes) Clone() Attributes { return maps.Clone(a) }

// Set changes the attribute that a given key has now a given value.
func (a Attributes) Set(key, value string) Attributes {
	if a == nil {
		return map[string]string{key: value}
	}
	a[key] = value
	return a
}

// Remove the key from the attributes.
func (a Attributes) Remove(key string) Attributes {
	if a != nil {
		delete(a, key)
	}
	return a
}

// Add a value to an attribute key.
func (a Attributes) Add(key, value string) Attributes {
	if a == nil {
		return map[string]string{key: value}
	}
	values := a.Values(key)
	if !slices.Contains(values, value) {
		values = append(values, value)
		a[key] = strings.Join(values, " ")
	}
	return a
}

// Values are the space separated values of an attribute.
func (a Attributes) Values(key string) []string {
	if a != nil {
		if value, ok := a[key]; ok {
			return strings.Fields(value)
		}
	}
	return nil
}

// Has the attribute key a value?
func (a Attributes) Has(key, value string) bool {
	return slices.Contains(a.Values(key), value)
}

// AddClass adds a value to the class attribute.
func (a Attributes) AddClass(class string) Attributes { return a.Add("class", class) }

// GetClasses returns the class values as a string slice
func (a Attributes) GetClasses() []string { return a.Values("class") }

// HasClass returns true, if attributes contains the given class.
func (a Attributes) HasClass(s string) bool { return a.Has("class", s) }

// AsAssoc returns the attributes as an assoc list.
//
// It is partly the reverse operation of [GetAttributes]:
// `maps.Equal(GetAttributes(a.AsAssoc()), a)`.
func (a Attributes) AsAssoc() *sx.Pair {
	var lb sx.ListBuilder
	for k, v := range a {
		lb.Add(sx.Cons(sx.MakeString(k), sx.MakeString(v)))
	}
	return lb.List()
}

// GetAttributes traverses a s-expression list and returns an attribute structure.
func GetAttributes(seq *sx.Pair) (result Attributes) {
	if seq == nil {
		return nil
	}
	for obj := range seq.Values() {
		pair, isPair := sx.GetPair(obj)
		if !isPair || pair == nil {
			continue
		}
		key := pair.Car()
		if !key.IsAtom() {
			continue
		}
		val := pair.Cdr()
		if tail, isTailPair := sx.GetPair(val); isTailPair {
			val = tail.Car()
		}
		if !val.IsAtom() {
			continue
		}
		result = result.Set(GoValue(key), GoValue(val))
	}
	return result
}
