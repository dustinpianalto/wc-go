package wc

import (
	"bytes"
	"unicode"
)

type ComplexCount struct {
	CharCount     int64
	WordCount     int64
	LineCount     int64
	MaxLineLength int64
}

func GetComplexCount(chunk []byte) ComplexCount {
	var count = ComplexCount{}
	word := false
	var lineLength int64
	runes := bytes.Runes(chunk)
	for _, b := range runes {
		count.CharCount++
		if b == '\n' {
			if lineLength > count.MaxLineLength {
				count.MaxLineLength = lineLength
			}
			lineLength = 0
			count.LineCount++
			if word {
				word = false
				count.WordCount++
			}
		} else if unicode.IsSpace(b) {
			lineLength++
			if word {
				word = false
				count.WordCount++
			}
		} else {
			lineLength++
			word = true
		}
	}
	return count
}

func ConcurrentComplexChunkCounter(chunks <-chan []byte, counts chan<- ComplexCount) {
	var totalCount ComplexCount
	for chunk := range chunks {
		count := GetComplexCount(chunk)
		totalCount.CharCount += count.CharCount
		totalCount.WordCount += count.WordCount
		totalCount.LineCount += count.LineCount
		if count.MaxLineLength > totalCount.MaxLineLength {
			totalCount.MaxLineLength = count.MaxLineLength
		}
	}
	counts <- totalCount
}
