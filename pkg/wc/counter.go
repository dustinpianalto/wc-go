package wc

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const BufferSize = 1024 * 1024 * 4

type Counter struct {
	Words         int64
	Chars         int64
	Lines         int64
	Bytes         int64
	MaxLineLength int64
}

func Count(filename string, cw, cc, cl, cb, mll bool) {
	if !cw && !cc && !cl && !cb && !mll {
		cw = true
		cc = false
		cl = true
		cb = true
		mll = false
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	processLine := cw || cc

	var c = &Counter{}
	numWorkers := 2 //runtime.NumCPU()

	if cl && !processLine {
		c.Lines = CountLines(file, numWorkers)
	} else {

		c = CountComplex(file, numWorkers)
	}

	if cb {
		fi, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Bytes = fi.Size()
	}

	if cl {
		fmt.Printf("%d ", c.Lines)
	}
	if cw {
		fmt.Printf("%d ", c.Words)
	}
	if cc {
		fmt.Printf("%d ", c.Chars)
	}
	if cb {
		fmt.Printf("%d ", c.Bytes)
	}
	if mll {
		fmt.Printf("%d ", c.MaxLineLength)
	}
	fmt.Printf("%s\n", filename)
}

func CountLines(file *os.File, numWorkers int) int64 {
	var lines int64

	chunks := make(chan []byte, numWorkers)
	counts := make(chan int64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go ConcurrentChunkCounter(chunks, counts)
	}

	for {
		buf := make([]byte, BufferSize)
		bytes, err := file.Read(buf)
		chunks <- buf[:bytes]
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}
	close(chunks)
	for i := 0; i < numWorkers; i++ {
		lines += <-counts
	}
	close(counts)
	return lines
}

func CountComplex(file *os.File, numWorkers int) *Counter {
	counter := Counter{}

	chunks := make(chan ComplexChunk, numWorkers)
	counts := make(chan ComplexCount, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go ConcurrentComplexChunkCounter(chunks, counts)
	}
	var lastRune rune = ' ' // Fake the first char being a space so that the first word is counted
	for {
		buf := make([]byte, BufferSize)
		count, err := file.Read(buf)
		chunks <- ComplexChunk{lastRune, buf[:count]}
		lastRune, _ = utf8.DecodeLastRune(buf[:count])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}
	close(chunks)
	for i := 0; i < numWorkers; i++ {
		count := <-counts
		counter.Lines += count.LineCount
		counter.Words += count.WordCount
		counter.Chars += count.CharCount
		if count.MaxLineLength > counter.MaxLineLength {
			counter.MaxLineLength = count.MaxLineLength
		}
	}
	close(counts)
	return &counter
}
