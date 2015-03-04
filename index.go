// Copyright (c) 2015, Bogdan S.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package bwt

import (
	"sort"
)

// BWT FM-index
type Index struct {
	SuffixArr              []int
	LastCol                []rune
	SymbolsOccurCountByPos map[rune][]int
	FirstOccurPos          map[rune]int
}

// Works with string, []byte, []rune
// Calculates SymbolsOccurCountByPos & FirstOccurPos
func (self *Index) Set(BWTLastCol interface{}) {
	switch BWTLastCol.(type) {
	case []rune:
		self.LastCol = BWTLastCol.([]rune)
	case string:
		self.LastCol = []rune(BWTLastCol.(string))
	case []byte:
		self.LastCol = []rune(string(BWTLastCol.([]byte)))
	}
	self.SymbolsOccurCountByPos = make(map[rune][]int)

	for i, symbol := range self.LastCol {
		for _, val := range self.SymbolsOccurCountByPos {
			val[i+1] = val[i]
		}
		if _, found := self.SymbolsOccurCountByPos[symbol]; !found {
			self.SymbolsOccurCountByPos[symbol] = make([]int, len(self.LastCol)+1)
		}
		self.SymbolsOccurCountByPos[symbol][i+1] += 1
	}

	var firstCol = RuneSlice(self.LastCol)
	sort.Sort(firstCol)
	self.FirstOccurPos = make(map[rune]int)

	for i, v := range firstCol {
		if _, found := self.FirstOccurPos[v]; !found {
			self.FirstOccurPos[v] = i
		}
	}
}

// Can't find deletions or insertions, only substitution
// Works with string, []byte, []rune
func (self *Index) Lookup(data interface{}) (result []int) {
	var (
		top    = 0
		bottom = len(self.LastCol)
	)
	result = []int{}
	var text = ReverseRunes(data)
	if len(text) == 0 {
		return
	}

	for _, symbol := range text {
		if _, ok := self.SymbolsOccurCountByPos[symbol]; !ok || top > bottom {
			return
		}
		top = self.FirstOccurPos[symbol] + self.SymbolsOccurCountByPos[symbol][top]
		bottom = self.FirstOccurPos[symbol] + self.SymbolsOccurCountByPos[symbol][bottom]
	}

	result = self.SuffixArr[top:bottom]
	return
}

// When threshold < 1 will reuse Index.Lookup method.
func (self *Index) LookupMismatches(data interface{}, threshold int) (result []int) {

	if threshold < 1 {
		return self.Lookup(data)
	}

	result = []int{}
	var text = ReverseRunes(data)
	if len(text) == 0 {
		return
	}

	// ( mismatches, pos )
	var matches = make([][2]int, len(self.LastCol))
	for i := range self.LastCol {
		matches[i] = [2]int{0, i}
	}

	for _, lookupSymbol := range text {

		matchesTmp := [][2]int{}

		for _, v := range matches {
			mismatches, pos := v[0], v[1]
			symbol := self.LastCol[pos]

			if lookupSymbol != symbol {
				mismatches += 1
			}
			if mismatches > threshold {
				continue
			}

			matchesTmp = append(matchesTmp, [2]int{
				mismatches,
				self.FirstOccurPos[symbol] + self.SymbolsOccurCountByPos[symbol][pos],
			})
		}
		matches = matchesTmp

		if len(matches) == 0 {
			return
		}
	}
	for _, v := range matches {
		result = append(result, self.SuffixArr[v[1]])
	}
	return
}
