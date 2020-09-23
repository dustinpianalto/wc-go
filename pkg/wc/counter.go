package wc

import (
	"fmt"
	"io"
	"os"
	"runtime"
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

func Count(filename string, cw, cc, cl, cb, mll bool) (Counter, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return Counter{}, err
	}
	defer file.Close()

	processLine := cw || cc

	var c = &Counter{}
	numWorkers := runtime.NumCPU()

	if cl && !processLine {
		c.Lines = CountLines(file, numWorkers)
	} else {

		c = CountComplex(file, numWorkers)
	}

	if cb {
		fi, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return Counter{}, err
		}
		c.Bytes = fi.Size()
	}
	return *c, nil
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
