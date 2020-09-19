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

type ComplexChunk struct {
	PrevRune rune
	Chunk    []byte
}

func GetComplexCount(chunk ComplexChunk) ComplexCount {
	var count = ComplexCount{}
	var lineLength int64
	runes := bytes.Runes(chunk.Chunk)
	prevRuneIsSpace := unicode.IsSpace(chunk.PrevRune)
	var linepos int64
	for _, b := range runes {
		count.CharCount++
		if b == '\n' || b == '\r' || b == '\f' {
			if linepos > lineLength {
				lineLength = linepos
			}
			linepos = 0
			if b == '\n' {
				count.LineCount++
			}
		}
		if unicode.IsSpace(b) {
			if b == '\t' {
				linepos += 8 - (linepos % 8)
			} else if b != '\n' && b != '\r' && b != '\f' && b != '\v' {
				linepos++
			}
			prevRuneIsSpace = true
		} else {
			linepos++
			if prevRuneIsSpace {
				count.WordCount++
			}
			prevRuneIsSpace = false
		}
	}
	count.MaxLineLength = lineLength
	return count
}

func ConcurrentComplexChunkCounter(chunks <-chan ComplexChunk, counts chan<- ComplexCount) {
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
