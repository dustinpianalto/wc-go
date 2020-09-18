package wc

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

const BufferSize = 1024 * 1024

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
	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	size := fi.Size()

	processLine := cw || cc

	var c = &Counter{}
	numWorkers := runtime.NumCPU()

	chunks := make(chan []byte, numWorkers)
	counts := make(chan int64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go ConcurrentChunkCounter(chunks, counts)
	}

	if cl && !processLine {
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
			c.Lines += <-counts
		}
		close(counts)
	}

	if cb {
		c.Bytes = size
	}

	if c.Lines > 0 {
		fmt.Printf("%d ", c.Lines)
	}
	if c.Words > 0 {
		fmt.Printf("%d ", c.Words)
	}
	if c.Chars > 0 {
		fmt.Printf("%d ", c.Words)
	}
	if c.Bytes > 0 {
		fmt.Printf("%d ", c.Bytes)
	}
	if c.MaxLineLength > 0 {
		fmt.Printf("%d ", c.MaxLineLength)
	}
	fmt.Printf("%s\n", filename)
}

//func (c *Counter) CountLines() {
//	for {
//		_, p, err := c.FileReader.ReadLine()
//		if err != nil && err != io.EOF {
//			panic(err)
//		}
//		if err == io.EOF {
//			break
//		}
//		if !p {
//			c.Lines++
//		}
//	}
//	c.Lines--
//}
