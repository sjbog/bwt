// Copyright (c) 2015, Bogdan S.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package bwt

import (
	"reflect"
	"sort"
	"testing"
)

func Test_BWT_SetLastCol(t *testing.T) {
	t.Parallel()

	var (
		bwt      = Index{}
		expected = map[rune][]int{
			'a': []int{0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 5, 6},
			'$': []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
			'm': []int{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'p': []int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			's': []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'b': []int{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'n': []int{0, 0, 0, 1, 1, 1, 2, 3, 3, 3, 3, 3, 3, 3, 3},
		}
		expectedFirstOccurPos = map[rune]int{'a': 1, 'p': 12, 's': 13, 'b': 7, 'm': 8, '$': 0, 'n': 9}
		emptyLastCol          = make([]rune, 0)
	)
	bwt.LastCol = emptyLastCol
	bwt.Set("smnpbnnaaaaa$a")

	if reflect.DeepEqual(bwt.LastCol, emptyLastCol) {
		t.Error("Index.Set shoult work with string")
	}

	bwt.LastCol = emptyLastCol
	bwt.Set([]byte("smnpbnnaaaaa$a"))

	if reflect.DeepEqual(bwt.LastCol, emptyLastCol) {
		t.Error("Index.Set shoult work with []byte")
	}

	bwt.LastCol = emptyLastCol
	bwt.Set([]rune("smnpbnnaaaaa$a"))

	if reflect.DeepEqual(bwt.LastCol, emptyLastCol) {
		t.Error("Index.Set shoult work with []rune")
	}
	if !reflect.DeepEqual(bwt.SymbolsOccurCountByPos, expected) {
		t.Error("Index( \"smnpbnnaaaaa$a\" ).SymbolsOccurCountByPos is incorrect")
	}
	if !reflect.DeepEqual(bwt.FirstOccurPos, expectedFirstOccurPos) {
		t.Errorf("Index( \"smnpbnnaaaaa$a\" ).FirstOccurPos is incorrect\nGot: %v\nExpected: %v", bwt.FirstOccurPos, expectedFirstOccurPos)
	}
}

func Test_BWT_Lookups(t *testing.T) {
	t.Parallel()

	var bwt = Index{
		LastCol:       []rune("smnpbnnaaaaa$a"),
		SuffixArr:     []int{13, 5, 3, 1, 7, 9, 11, 6, 4, 2, 8, 10, 0, 12},
		FirstOccurPos: map[rune]int{'p': 12, 's': 13, 'a': 1, 'b': 7, 'm': 8, '$': 0, 'n': 9},
		SymbolsOccurCountByPos: map[rune][]int{
			'a': []int{0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 5, 6},
			'$': []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
			'm': []int{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'p': []int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			's': []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'b': []int{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			'n': []int{0, 0, 0, 1, 1, 1, 2, 3, 3, 3, 3, 3, 3, 3, 3},
		},
	}

	var (
		result   = bwt.Lookup("ana")
		expected = []int{1, 7, 9}
	)
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) || !reflect.DeepEqual(result, bwt.LookupMismatches("ana", 0)) {
		t.Errorf("Index( 'panamabananas$' ).Lookup( 'ana' ) = %v\n\tGot: %v", expected, result)
	}

	result = bwt.Lookup([]byte("ana"))
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Index.Lookup should work with []byte")
	}

	result = bwt.Lookup([]rune("ana"))
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Index.Lookup should work with []rune")
	}

	result = bwt.LookupMismatches("ana", 1)
	sort.Ints(result)
	expected = []int{1, 3, 5, 7, 9}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Index( 'panamabananas$' ).LookupMismatches( 'ana' ) = %v\n\tGot: %v", expected, result)
	}

	result = bwt.LookupMismatches([]byte("ana"), 1)
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Index.LookupMismatches should work with []byte")
	}

	result = bwt.LookupMismatches([]rune("ana"), 1)
	sort.Ints(result)

	if !reflect.DeepEqual(result, expected) {
		t.Error("Index.LookupMismatches should work with []rune")
	}

	result = bwt.LookupMismatches("belle", 2)

	if !reflect.DeepEqual(result, []int{}) ||
		!reflect.DeepEqual(bwt.Lookup("anabelle"), []int{}) ||
		!reflect.DeepEqual(bwt.Lookup("anab"), []int{}) ||
		!reflect.DeepEqual(bwt.Lookup(""), []int{}) ||
		!reflect.DeepEqual(bwt.LookupMismatches("", 1), []int{}) {
		t.Error(`Index.Lookup should return empty result`)
	}
}
