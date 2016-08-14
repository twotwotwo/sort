package radix

import (
	"bytes"
	"math"
	"sort"
)

// Copyright 2009 The Go Authors.
// Copyright 2014-6 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Copying the helpers from package sort and adding Key funcs. index.go
// has the abbreviated-key string sort stuff.

// Float32Key provides a sort key for a float32. Use with Float32Less.
func Float32Key(f float32) Key {
	b := Key(math.Float32bits(f)) << 32
	b ^= ^(b>>63 - 1) | (1 << 63)
	return b
}

// Float32Less compares float32s, treating NaN as greater than all numbers.
func Float32Less(f, g float32) bool {
	return Float32Key(f) < Float32Key(g)
}

// Float64Key provides a sort key a float64. Use with Float64Less.
func Float64Key(f float64) Key {
	b := math.Float64bits(f)
	b ^= ^(b>>63 - 1) | (1 << 63)
	return Key(b)
}

// Float64Less compares float64s, treating NaN as greater than all numbers.
func Float64Less(f, g float64) bool {
	return Float64Key(f) < Float64Key(g)
}

// Int64Key provides a sort key for an int.
func Int64Key(i int64) Key {
	return Key(i) ^ 1<<63
}

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of ints.
func (p IntSlice) Key(i int) Key { return Int64Key(int64(p[i])) }

// Sort is a convenience method.
func (p IntSlice) Sort() { Sort(p) }

// Int32Slice attaches the methods of Interface to []int32, sorting in increasing order.
type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of int32s.
func (p Int32Slice) Key(i int) Key { return Int64Key(int64(p[i])) }

// Sort is a convenience method.
func (p Int32Slice) Sort() { Sort(p) }

// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of int64s.
func (p Int64Slice) Key(i int) Key { return Int64Key(p[i]) }

// Sort is a convenience method.
func (p Int64Slice) Sort() { Sort(p) }

// UintSlice attaches the methods of Interface to []uint, sorting in increasing order.
type UintSlice []uint

func (p UintSlice) Len() int           { return len(p) }
func (p UintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p UintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of uints.
func (p UintSlice) Key(i int) Key { return Key(p[i]) }

// Sort is a convenience method.
func (p UintSlice) Sort() { Sort(p) }

// Uint32Slice attaches the methods of Interface to []int32, sorting in increasing order.
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of uint32s.
func (p Uint32Slice) Key(i int) Key { return Key(p[i]) }

// Sort is a convenience method.
func (p Uint32Slice) Sort() { Sort(p) }

// Uint64Slice attaches the methods of Interface to []uint64, sorting in increasing order.
type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of uint64s.
func (p Uint64Slice) Key(i int) Key { return Key(p[i]) }

// Sort is a convenience method.
func (p Uint64Slice) Sort() { Sort(p) }

// Float32Slice attaches the methods of Interface to []uint32, sorting in increasing order, NaNs last.
type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return Float32Less(p[i], p[j]) }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of float32s.
func (p Float32Slice) Key(i int) Key { return Float32Key(p[i]) }

// Sort is a convenience method.
func (p Float32Slice) Sort() { Sort(p) }

// Float64Slice attaches the methods of Interface to []float64, sorting in increasing order, NaNs last.
type Float64Slice []float64

func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return Float64Less(p[i], p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Key provides a sort key for an item in a slice of float64s.
func (p Float64Slice) Key(i int) Key { return Float64Key(p[i]) }

// Sort is a convenience method.
func (p Float64Slice) Sort() { Sort(p) }

// StringSlice attaches the methods of Interface to []float64, sorting in increasing order, NaNs last.
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// StringAt returns an item from a string slice.
func (p StringSlice) StringAt(i int) string { return p[i] }

// Sort is a convenience method.
func (p StringSlice) Sort() { SortStrings(p) }

// BytesSlice attaches the methods of Interface to []float64, sorting in increasing order, NaNs last.
type BytesSlice [][]byte

func (p BytesSlice) Len() int           { return len(p) }
func (p BytesSlice) Less(i, j int) bool { return bytes.Compare(p[i], p[j]) < 0 }
func (p BytesSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// BytesAt reutrns an item from a slice of byte slices.
func (p BytesSlice) BytesAt(i int) []byte { return p[i] }

// Sort is a convenience method.
func (p BytesSlice) Sort() { SortBytes(p) }

// Ints sorts a slice of ints in increasing order.
func Ints(a []int) { IntSlice(a).Sort() }

// Int32s sorts a slice of int32s in increasing order.
func Int32s(a []int32) { Int32Slice(a).Sort() }

// Int64s sorts a slice of int64s in increasing order.
func Int64s(a []int64) { Int64Slice(a).Sort() }

// Uints sorts a slice of ints in increasing order.
func Uints(a []uint) { UintSlice(a).Sort() }

// Uint32s sorts a slice of uint64s in increasing order.
func Uint32s(a []uint32) { Uint32Slice(a).Sort() }

// Uint64s sorts a slice of uint64s in increasing order.
func Uint64s(a []uint64) { Uint64Slice(a).Sort() }

// Float32s sorts a slice of uint64s in increasing order, NaNs last.
func Float32s(a []float32) { Float32Slice(a).Sort() }

// Float64s sorts a slice of uint64s in increasing order, NaNs last.
func Float64s(a []float64) { Float64Slice(a).Sort() }

// Strings sorts a slice of strings in increasing order.
func Strings(a []string) { StringSlice(a).Sort() }

// Bytes sorts a slice of byte slices in increasing order.
func Bytes(a [][]byte) { BytesSlice(a).Sort() }

// IntsAreSorted tests whether a slice of ints is sorted in increasing order.
func IntsAreSorted(a []int) bool { return sort.IsSorted(IntSlice(a)) }

// Int32sAreSorted tests whether a slice of int32s is sorted in increasing order.
func Int32sAreSorted(a []int32) bool { return sort.IsSorted(Int32Slice(a)) }

// Int64sAreSorted tests whether a slice of int64s is sorted in increasing order.
func Int64sAreSorted(a []int64) bool { return sort.IsSorted(Int64Slice(a)) }

// UintsAreSorted tests whether a slice of ints is sorted in increasing order.
func UintsAreSorted(a []uint) bool { return sort.IsSorted(UintSlice(a)) }

// Uint32sAreSorted tests whether a slice of uint32s is sorted in increasing order.
func Uint32sAreSorted(a []uint32) bool { return sort.IsSorted(Uint32Slice(a)) }

// Uint64sAreSorted tests whether a slice of uint64s is sorted in increasing order.
func Uint64sAreSorted(a []uint64) bool { return sort.IsSorted(Uint64Slice(a)) }

// Float32sAreSorted tests whether a slice of float32s is sorted in increasing order, NaNs last.
func Float32sAreSorted(a []float32) bool { return sort.IsSorted(Float32Slice(a)) }

// Float64sAreSorted tests whether a slice of float64s is sorted in increasing order, NaNs last.
func Float64sAreSorted(a []float64) bool { return sort.IsSorted(Float64Slice(a)) }

// StringsAreSorted tests whether a slice of strings is sorted in increasing order.
func StringsAreSorted(a []string) bool { return sort.IsSorted(StringSlice(a)) }

// BytesAreSorted tests whether a slice of byte slices is sorted in increasing order.
func BytesAreSorted(a [][]byte) bool { return sort.IsSorted(BytesSlice(a)) }
