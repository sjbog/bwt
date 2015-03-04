// Copyright (c) 2015, Bogdan S.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package bwt

import (
	"reflect"
	"testing"
)

func Test_BWT(t *testing.T) {
	t.Parallel()

	var (
		data                           = "panamabananas"
		resultLastCol, resultSuffixArr = BWT(data, 0)
		expectedLastCol                = []rune("smnpbnnaaaaa\u0003a")
		expectedSuffixArr              = []int{13, 5, 3, 1, 7, 9, 11, 6, 4, 2, 8, 10, 0, 12}
	)
	if !reflect.DeepEqual(resultLastCol, expectedLastCol) {
		t.Errorf("BWT( \"%v\" ) LastCol is incorrect\nGot: %v\nExpected: %v", data, resultLastCol, expectedLastCol)
	}
	if !reflect.DeepEqual(resultSuffixArr, expectedSuffixArr) {
		t.Errorf("BWT( \"%v\" ) SuffixArr is incorrect\nGot: %v\nExpected: %v", data, resultSuffixArr, expectedSuffixArr)
	}
}

func Test_BWT_Rotation(t *testing.T) {
	t.Parallel()

	var (
		data     = "line"
		index    = New(data, -1)
		result   = index.Lookup("eli")
		expected = []int{}
	)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("BWT( \"%v\" ) stop rune is incorrect, Index.Lookup( \"%v\" )\nGot: %v\nExpected: %v", data, "eli", result, expected)
	}
}
