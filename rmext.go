package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/GreenRaccoon23/slices"
)

var (
	PrintBase bool
	PrintFull bool
	Paths     []string
)

func init() {

	if helpNeeded() {
		help(1)
	} else if helpWanted() {
		help(0)
	}

	flags()

	if !validArgs() {
		help(1)
	}
}

func helpNeeded() bool {
	if noArgs := (len(os.Args) == 1); noArgs {
		return true
	}
	return false
}

func helpWanted() bool {
	switch os.Args[1] {
	case "-h", "h", "help", "--help", "-H", "H", "HELP", "--HELP", "-help", "--h", "--H":
		return true
	}
	return false
}

// Print help and exit with a status code.
func help(status int) {
	defer os.Exit(status)
	fmt.Printf(
		"%s\n\n  %s\n\n  %s\n%s\n\n  %s\n%s\n%s\n%s\n%s\n\n  %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		"rmext",
		"Usage: rmext file...",
		"Description:",
		"    Print the basename of files WITHOUT their extension.",
		"Options:",
		"   -b         Print the basename of files without their path",
		"   -f         Print the full path of files",
		"   (default)  Print the path of files exactly as they were entered",
		"                (minus any file extensions)",
		"Notes:",
		"    This program relies on the system's list of known mimetypes.",
		"	 On Linux, usually the file '/etc/mime.types' stores all",
		"      MIME types known by the system.",
		"	 On Windows, the registry mess stores them.",
		"    In rare cases, the system may not have a mime entry for an",
		"      uncommon file type, and this program will not remove the",
		"      extension associated with that file type.",
		"    An example is '.go' files; Linux distributions usually do not",
		"      include a mime type for '.go' files by default (yet).",
	)
}

func flags() {

	bools := []string{"-b", "-f"}

	flag.BoolVar(&PrintBase, "b", false, "")
	flag.BoolVar(&PrintFull, "f", false, "")
	flag.Parse()

	Paths = os.Args[1:]
	Paths = slices.Filter(Paths, bools...)
}

func validArgs() bool {

	if PrintBase && PrintFull {
		return false
	}

	if len(Paths) == 0 {
		return false
	}

	return true
}

func main() {
	defer os.Exit(0)
	var err error
	for _, t := range Paths {
		switch {
		case PrintBase:
			t = filepath.Base(t)
		case PrintFull:
			t, err = filepath.Abs(t)
			chkerr(err)
		}
		base, _ := splitExt(t)
		fmt.Println(base)
	}
}

// Exit with status 1 if an error occurs.
func chkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Split the extension off a filename.
// Return the basename and the extension.
func splitExt(filename string) (base, ext string) {
	base = filepath.Clean(filename)
	for {
		testext := filepath.Ext(base)
		if testext == "" || mime.TypeByExtension(testext) == "" {
			return
		}
		ext = concat(testext, ext)
		base = strings.TrimSuffix(base, testext)
	}
}

// Concatenate strings.
func concat(slc ...string) string {
	b := bytes.NewBuffer(nil)
	defer b.Reset()
	for _, s := range slc {
		b.WriteString(s)
	}
	return b.String()
}
