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
	"fmt"

	"github.com/droundy/goopt"
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
	if len(goopt.Args) == 0 {
		fmt.Println(goopt.Help())
	} else {
		if !*fWords && !*fChars && !*fLines && !*fBytes && !*fMaxLineLength {
			*fWords = true
			*fChars = false
			*fLines = true
			*fBytes = true
			*fMaxLineLength = false
		}
		var TotalCounts wc.Counter
		for _, a := range goopt.Args {
			count, err := wc.Count(a, *fWords, *fChars, *fLines, *fBytes, *fMaxLineLength)
			if err == nil {
				TotalCounts.Lines += count.Lines
				TotalCounts.Bytes += count.Bytes
				TotalCounts.Words += count.Words
				TotalCounts.Chars += count.Chars
				TotalCounts.MaxLineLength += count.MaxLineLength
				if *fLines {
					fmt.Printf("%d", count.Lines)
				}
				fmt.Printf(" ")
				if *fWords {
					fmt.Printf("%d", count.Words)
				}
				fmt.Printf(" ")
				if *fChars {
					fmt.Printf("%d", count.Chars)
				}
				fmt.Printf(" ")
				if *fBytes {
					fmt.Printf("%d", count.Bytes)
				}
				fmt.Printf(" ")
				if *fMaxLineLength {
					fmt.Printf("%d", count.MaxLineLength)
				}
				fmt.Printf(" %s\n", a)
			}

		}
		if len(goopt.Args) > 1 {
			if *fLines {
				fmt.Printf("%d", TotalCounts.Lines)
			}
			fmt.Printf(" ")
			if *fWords {
				fmt.Printf("%d", TotalCounts.Words)
			}
			fmt.Printf(" ")
			if *fChars {
				fmt.Printf("%d", TotalCounts.Chars)
			}
			fmt.Printf(" ")
			if *fBytes {
				fmt.Printf("%d", TotalCounts.Bytes)
			}
			fmt.Printf(" ")
			if *fMaxLineLength {
				fmt.Printf("%d", TotalCounts.MaxLineLength)
			}
			fmt.Println(" total")
		}
	}
}
