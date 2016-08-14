// Copyright 2015-6 Randall Farmer. All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package radix

import "sort"

type index struct {
	Keys []Key
	sort.Interface
}

// Less compares index elements by their Keys, falling back to Data.Less for
// equal-keyed items.
func (idx *index) Less(i, j int) bool {
	ki, kj := idx.Key(i), idx.Key(j)
	return ki < kj || (ki == kj && idx.Interface.Less(i, j))
}

// Swap swaps both the Keys and the inderlying data items at indices i and
// j.
func (idx *index) Swap(i, j int) {
	idx.Keys[i], idx.Keys[j] = idx.Keys[j], idx.Keys[i]
	idx.Interface.Swap(i, j)
}

// Key returns the uint64 key at index i.
func (idx *index) Key(i int) Key {
	return idx.Keys[i]
}

// IndexBuilder describes data that can be sorted with SortIndex. Implement it
// only when it's a net benefit to allocate space for sort keys but SortBytes
// or SortStrings can't be used.
// For example, it is worth considering for sorting struct pointers by a numeric
// key, or for sorting strings with a custom ordering.
type IndexBuilder interface {
	sort.Interface
	SetKeys([]Key, int)
}

// SetKeys above has the int parameter so we could ask for only some keys at once,
// like if we had "pages" of keys or set keys in parallel. We can remove it if
// we're not going to do either one of those.

// SortIndex allocates space to store a Key for each item then uses it to sort,
// using data.Less as a tie-breaker for equal-keyed items.
func SortIndex(data IndexBuilder) {
	l := data.Len()
	idx := &index{
		Keys:      make([]Key, l),
		Interface: data,
	}
	data.SetKeys(idx.Keys, 0)
	Sort(idx)
}

// SortBytes sorts a BytesInterface, using temporary space to speed the sort.
func SortBytes(data BytesInterface) {
	SortIndex(bytesIndexBuilder{data})
}

// SortStrings sorts a StringInterface, using temporary space to speed the sort.
func SortStrings(data StringInterface) {
	SortIndex(stringIndexBuilder{data})
}

// StringInterface describes a collection of data sortable by a string key.
type StringInterface interface {
	sort.Interface
	StringAt(int) string
}

type stringIndexBuilder struct {
	StringInterface
}

func stringKey(s string) Key {
	k := Key(0)
	shift := uint(56)
	for j := 0; j < 8 && j < len(s); j++ {
		k ^= Key(s[j]) << shift
		shift -= 8
	}
	return k
}

func (sib stringIndexBuilder) SetKeys(keys []Key, a int) {
	l := sib.Len()
	for i := range keys {
		if i+a == l {
			break
		}
		keys[i] = stringKey(sib.StringAt(i + a))
	}
}

// BytesInterface describes a collection of data sortable by a []byte key.
type BytesInterface interface {
	sort.Interface
	BytesAt(int) []byte
}

type bytesIndexBuilder struct {
	BytesInterface
}

func bytesKey(b []byte) Key {
	k := Key(0)
	shift := uint(56)
	for j := 0; j < 8 && j < len(b); j++ {
		k ^= Key(b[j]) << shift
		shift -= 8
	}
	return k
}

func (bib bytesIndexBuilder) SetKeys(keys []Key, a int) {
	l := bib.Len()
	for i := range keys {
		if i+a == l {
			break
		}
		keys[i] = bytesKey(bib.BytesAt(i + a))
	}
}
