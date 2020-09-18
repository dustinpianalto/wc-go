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
		for _, a := range goopt.Args {
			wc.Count(a, *fWords, *fChars, *fLines, *fBytes, *fMaxLineLength)
		}
	}
}
