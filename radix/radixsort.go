// Copyright 2014-6 Randall Farmer. All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package radix is for playing around with radix sort APIs and implementations.
// Please don't use it in production; it's unstable and (currently) untested.
// (Look at github.com/twotwotwo/sorts instead.)
package radix

import "sort"

// Key is a uint64 by which to radix sort.
type Key uint64

// Interface describes a collection sortable with radix sort.
type Interface interface {
	sort.Interface
	Key(i int) Key
}

const radix = 8
const mask = 1<<radix - 1
const qSortCutoff = 1 << 7

// check checks data is sorted, and if it isn't, sees if it's because of a Less/Key
// inconsistency, and panics with an appropriate message. We could add it to Sort()
// if we want to be paranoid about users correctly implementing Key and Less
func check(data Interface, l int) {
	for i := 1; i < l; i++ {
		if data.Less(i, i-1) {
			if data.Key(i) > data.Key(i-1) {
				panic("sort: Less and Key do not order items the same way")
			}
			panic("sort: failed to sort data; could be nondeterministic Less or Key or race condition")
		}
	}
}

// guessInitialShift samples data to guess the highest bit in the key
// that varies, and returns a corresponding value of 'shift' to try
// in radix sort. (For, say, an array of numbers 0 to 2^32-1, the
// shift is 24 with radix 8.) A wrong guess is corrected in radixSort
// after a pass over all the keys.
func guessInitialShift(data Interface, l int) uint {
	step := l >> 5
	if l > 1<<16 {
		step = l >> 8
	}
	if step == 0 { // only for tests w/qSortCutoff lowered
		step = 1
	}
	min := data.Key(l - 1)
	max := min
	for i := 0; i < l; i += step {
		k := data.Key(i)
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	diff := min ^ max
	log2diff := 0
	for diff != 0 {
		// RF: could use 1 instead of radix as twotwotwo/sorts
		// does, and be faster when right but wrong more often.
		log2diff += radix
		diff >>= radix
	}
	shiftGuess := log2diff - radix
	if shiftGuess < 0 {
		return 0
	}

	return uint(shiftGuess)
}

// Sort sorts data
func Sort(data Interface) {
	l := data.Len()

	if l < qSortCutoff {
		qSort(data, 0, l)
		return
	}

	shift := guessInitialShift(data, l)
	radixSort(data, shift, 0, l, &[1 << radix]int{})
}

func radixSort(data Interface, shift uint, a, b int, scratch *[1 << radix]int) {
	if b-a < qSortCutoff {
		qSort(data, a, b)
		return
	}

	// use a single pass over the keys to bucket data and find min/max
	// (for skipping over bits that are always identical)
	bucketStarts := scratch
	for i := range bucketStarts {
		bucketStarts[i] = 0
	}
	min := data.Key(a)
	max := min
	for i := a; i < b; i++ {
		k := data.Key(i)
		bucketStarts[(k>>shift)&mask]++
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
		// RF: we could check here if min^max>>(shift+radix) != 0 when
		// i reaches, say, a+((b-a)>>8). That'd reduce the cost when
		// guessInitialShift guesses wrong, without making the worst
		// case that much worse.
	}

	// skip past common prefixes, bail if all keys equal
	diff := min ^ max
	if diff == 0 {
		qSort(data, a, b)
		return
	}
	if diff>>shift == 0 || diff>>(shift+radix) != 0 {
		// find highest 1 bit in diff
		log2diff := 0
		for diff != 0 {
			log2diff++
			diff >>= 1
		}
		nextShift := log2diff - radix
		if nextShift < 0 {
			nextShift = 0
		}
		radixSort(data, uint(nextShift), a, b, scratch)
		return
	}

	var bucketEnds [1 << radix]int
	pos := a
	for i, c := range bucketStarts {
		bucketStarts[i] = pos
		pos += c
		bucketEnds[i] = pos
	}

	for curBucket, bucketEnd := range bucketEnds {
		i := bucketStarts[curBucket]
		for i < bucketEnd {
			destBucket := (data.Key(i) >> shift) & mask
			if destBucket == Key(curBucket) {
				i++
				bucketStarts[destBucket]++
				continue
			}
			data.Swap(i, bucketStarts[destBucket])
			bucketStarts[destBucket]++
		}
	}

	if shift == 0 {
		pos = a
		for _, end := range bucketEnds {
			if end > pos+1 {
				qSort(data, pos, end)
			}
			pos = end
		}
		return
	}

	nextShift := shift - radix
	if shift < radix {
		nextShift = 0
	}
	pos = a
	for _, end := range bucketEnds {
		if end > pos+1 {
			radixSort(data, nextShift, pos, end, scratch)
		}
		pos = end
	}
}
