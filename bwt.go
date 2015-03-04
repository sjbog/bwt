// Copyright (c) 2015, Bogdan S.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

// Implements FM-Index (https://en.wikipedia.org/wiki/FM-index),
// full-text substring index based on the Burrows-Wheeler transform.
// One should use "index/suffixarray" pkg
package bwt

import (
	"container/heap"
	"unicode/utf8"
)

func New(data interface{}, maxDepth int) *Index {
	var LastCol, SuffixArr = BWT(data, maxDepth)
	var index = &Index{
		SuffixArr: SuffixArr,
	}
	index.Set(LastCol)
	return index
}

// Naive string matrix sorting (via priority queue)
// Uses U+0003(\u0003) End-of-text (ETX) character to indicate text end and separate rotations
// maxDepth < 1 will sort full text, otherwise prefixes are sorted
func BWT(data interface{}, maxDepth int) (bwt []rune, suffixArray []int) {
	var (
		pq   = new(pq)
		text []rune
	)
	switch data.(type) {
	case []rune:
		text = data.([]rune)
	case string:
		text = []rune(data.(string))
	case []byte:
		text = []rune(string(data.([]byte)))
	}
	// \u0003 ETX is a stop rune - otherwise "line" -> will match rotated "eli"
	text = append(text, 0x3)

	bwt = make([]rune, len(text))
	suffixArray = make([]int, len(text))

	if maxDepth < 1 || len(text) <= maxDepth {
		// full string
		for pos, _ := range text {
			// text[ pos : ] + text[ : pos ]
			row := text[pos:]
			row = append(row, text[:pos]...)

			heap.Push(pq, bwtItem{
				Pos:   pos,
				Value: string(row),
			})
		}
	} else {
		// prefix
		for pos, _ := range text {
			// text[ pos : ] + text[ : pos ]
			row := text[pos:]
			row = append(row, text[:pos]...)

			// row[ : maxDepth ] + row[ -1 ]
			lastV := row[len(row)-1]
			row = row[:maxDepth]
			row = append(row, lastV)

			heap.Push(pq, bwtItem{
				Pos:   pos,
				Value: string(row),
			})
		}
	}

	for i := 0; pq.Len() > 0; i++ {
		item := heap.Pop(pq).(bwtItem)
		suffixArray[i] = item.Pos
		bwt[i], _ = utf8.DecodeLastRuneInString(item.Value)
	}
	return
}

type bwtItem struct {
	Value string
	Pos   int
}

type pq []bwtItem

func (self *pq) Len() int {
	return len(*self)
}

func (self *pq) Less(i, j int) bool {
	return (*self)[i].Value < (*self)[j].Value
}

func (self *pq) Swap(i, j int) {
	(*self)[i], (*self)[j] = (*self)[j], (*self)[i]
}

func (self *pq) Push(x interface{}) {
	*self = append(*self, x.(bwtItem))
}

func (self *pq) Pop() interface{} {
	if self.Len() == 0 {
		return nil
	}
	var last_index = self.Len() - 1
	var item = (*self)[last_index]

	(*self) = (*self)[0:last_index]

	return item
}

type RuneSlice []rune

func (s RuneSlice) Len() int           { return len(s) }
func (s RuneSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s RuneSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Works with string, []byte, []rune
func ReverseRunes(data interface{}) (res []rune) {
	switch data.(type) {
	case string:
		res = make([]rune, utf8.RuneCountInString(data.(string)))
		var i = len(res)

		for _, symbol := range data.(string) {
			i--
			res[i] = symbol
		}
	case []byte:
		res = make([]rune, utf8.RuneCount(data.([]byte)))
		var i = len(res)

		for _, symbol := range string(data.([]byte)) {
			i--
			res[i] = symbol
		}
	case []rune:
		res = make([]rune, len(data.([]rune)))
		var i = len(res)

		for _, symbol := range data.([]rune) {
			i--
			res[i] = symbol
		}
	}
	return
}
