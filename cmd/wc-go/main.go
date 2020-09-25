package main

// Print  newline,  word, and byte counts for each FILE, and a total line if more than
//       one FILE is specified.  A word is a non-zero-length sequence of  characters  delimited by white space.
//
//       With no FILE, or when FILE is -, read standard input.
//
//       The  options  below  may  be used to select which counts are printed, always in the
//       following order: newline, word, character, byte, maximum line length.
//
//       -c, --bytes
//              print the byte counts
//
//       -m, --chars
//              print the character counts
//
//       -l, --lines
//              print the newline counts
//
//       --files0-from=F
//              read input from the files specified by NUL-terminated names in file F; If  F
//              is - then read names from standard input
//
//       -L, --max-line-length
//              print the maximum display width
//
//       -w, --words
//              print the word counts
//
//       --help display this help and exit
//
//       --version
//              output version information and exit

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/droundy/goopt"
	"github.com/dustinpianalto/quotearg"
	"github.com/dustinpianalto/wc-go/pkg/wc"
)

var (
	fBytes = goopt.Flag([]string{"-c", "--bytes"},
		[]string{"--no-bytes"},
		"print the byte count",
		"exclude the byte count")
	fChars = goopt.Flag([]string{"-m", "--chars"},
		[]string{"--no-chars"},
		"print the character counts",
		"exlude the character counts")
	fLines = goopt.Flag([]string{"-l", "--lines"},
		[]string{"--no-lines"},
		"print the newline count",
		"exclude the newline count")
	fFilesFrom = goopt.String([]string{"--files0-from"},
		"",
		"read input from the files specified by NUL-terminated names "+
			"in file; If file is '-' the read names from standard input")
	fMaxLineLength = goopt.Flag([]string{"-L", "--max-line-length"},
		[]string{"--no-max-line-lenght"},
		"print the maximum line length of the input",
		"exclude the maximum line length")
	fWords = goopt.Flag([]string{"-w", "--words"},
		[]string{"--no-words"},
		"print the word count",
		"exclude the word count")
)

func main() {
	goopt.Version = "v0.0.0a"
	goopt.Parse(nil)
	if len(goopt.Args) == 0 && *fFilesFrom == "" {
		fmt.Println(goopt.Help())
	} else {

		if !*fWords && !*fChars && !*fLines && !*fBytes && !*fMaxLineLength {
			*fWords = true
			*fChars = false
			*fLines = true
			*fBytes = true
			*fMaxLineLength = false
		}

		var files []string
		if *fFilesFrom != "" {
			fFile, err := os.Open(*fFilesFrom)
			if err != nil {
				log.Fatalf("Cannot open file %s: %s", fFilesFrom, err.Error())
			}
			reader := bufio.NewReader(fFile)
			for {
				s, err := reader.ReadString(0x00)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
				files = append(files, strings.TrimRight(s, "\x00"))
			}

			fFile.Close()
		}

		files = append(files, goopt.Args...)

		var maxBytes int64
		if len(files) > 1 ||
			((*fWords && *fChars) ||
				(*fWords && *fLines) ||
				(*fWords && *fBytes) ||
				(*fWords && *fMaxLineLength) ||
				(*fChars && *fLines) ||
				(*fChars && *fBytes) ||
				(*fChars && *fMaxLineLength) ||
				(*fLines && *fBytes) ||
				(*fLines && *fMaxLineLength) ||
				(*fBytes && *fMaxLineLength)) {
			for _, f := range files {
				fp, err := os.Open(f)
				if err != nil {
					continue
				}
				fi, err := fp.Stat()
				if err != nil {
					continue
				}
				b := fi.Size()
				maxBytes += b
			}
		}

		var maxStrLen = intLen(maxBytes)

		var TotalCounts wc.Counter
		for _, a := range files {
			count, err := wc.Count(a, *fWords, *fChars, *fLines, *fBytes, *fMaxLineLength)
			if err == nil {
				TotalCounts.Lines += count.Lines
				TotalCounts.Bytes += count.Bytes
				TotalCounts.Words += count.Words
				TotalCounts.Chars += count.Chars
				TotalCounts.MaxLineLength += count.MaxLineLength
				if *fLines {
					fmt.Printf("%*d ", maxStrLen, count.Lines)
				}
				if *fWords {
					fmt.Printf("%*d ", maxStrLen, count.Words)
				}
				if *fChars {
					fmt.Printf("%*d ", maxStrLen, count.Chars)
				}
				if *fBytes {
					fmt.Printf("%*d ", maxStrLen, count.Bytes)
				}
				if *fMaxLineLength {
					fmt.Printf("%*d ", maxStrLen, count.MaxLineLength)
				}
				if strings.Contains(a, "\n") {
					n := quotearg.Quote([]rune(a), quotearg.ShellEscapeAlwaysQuotingStyle, 0, 0, '\'', '\'')
					fmt.Printf("%s\n", string(n))
				} else {
					fmt.Printf("%s\n", a)

				}
			}
		}

		if len(files) > 1 {
			if *fLines {
				fmt.Printf("%*d ", maxStrLen, TotalCounts.Lines)
			}
			if *fWords {
				fmt.Printf("%*d ", maxStrLen, TotalCounts.Words)
			}
			if *fChars {
				fmt.Printf("%*d ", maxStrLen, TotalCounts.Chars)
			}
			if *fBytes {
				fmt.Printf("%*d ", maxStrLen, TotalCounts.Bytes)
			}
			if *fMaxLineLength {
				fmt.Printf("%*d ", maxStrLen, TotalCounts.MaxLineLength)
			}
			fmt.Println("total")
		}
	}
}

func intLen(i int64) int64 {
	var count int64
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
