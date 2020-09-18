package wc

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Counter struct {
	FileReader    *bufio.Reader
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
	}
	defer file.Close()
	processLine := cw || cc
	var c = Counter{FileReader: bufio.NewReader(file)}
	if cl && !processLine {
		c.CountLines(cb)
	}
	fmt.Printf("%d %s\n", c.Lines, filename)
}

func (c *Counter) CountLines(cb bool) {
	for {
		r, s, err := c.FileReader.ReadRune()
		log.Printf("%#v, %#v, %#v", r, s, err)
		if err != nil {
			break
		}
		if r == '\n' {
			c.Lines++
		}
		if cb {
			c.Bytes += int64(s)
		}
	}
}
