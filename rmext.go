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
)

var (
	doBasename bool
	doFullpath bool
	targets    []string
)

func init() {
	chkHelp()
	flags()
}

// Check whether user requested help.
func chkHelp() {
	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "-h", "h", "help", "--help", "-H", "H", "HELP", "--HELP", "-help", "--h", "--H":
		help(0)
	}

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

// Parse user arguments and modify global variables accordingly.
func flags() {
	// Program requires at least one user argument.
	// Print help and exit with status 1 if none were received.
	if len(os.Args) < 2 {
		help(1)
	}

	// Parse commandline arguments.
	flag.BoolVar(&doBasename, "b", false, "")
	flag.BoolVar(&doFullpath, "f", false, "")
	flag.Parse()

	// Modify global variables based on commandline arguments.
	targets = os.Args[1:]
	if doBasename && doFullpath {
		help(1)
	}
	if !doBasename && !doFullpath {
		return
	}
	bools := []string{"-b", "-f"}
	targets = filter(targets, bools...)
	if len(targets) == 0 {
		help(1)
	}
	return
}

// Remove elements in a slice (if they exist).
// Only remove EXACT matches.
func filter(slc []string, args ...string) (filtered []string) {
	for _, s := range slc {
		if slcHas(slc, s) {
			continue
		}
		filtered = append(filtered, s)
	}
	return
}

// Check whether a slice contains a string.
// Only return true if an element in the slice EXACTLY matches the string.
// If testing for more than one string,
//   return true if ANY of them match an element in the slice.
func slcHas(slc []string, args ...string) bool {
	for _, s := range slc {
		for _, a := range args {
			if s == a {
				return true
			}
		}
	}
	return false
}

func main() {
	defer os.Exit(0)
	var err error
	for _, t := range targets {
		switch {
		case doBasename:
			t = filepath.Base(t)
		case doFullpath:
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