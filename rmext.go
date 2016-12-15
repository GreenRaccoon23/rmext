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

	if err := modPaths(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	defer os.Exit(0)

	for _, fpath := range Paths {
		name, _ := splitExt(fpath)
		fmt.Println(name)
	}
}

func helpNeeded() bool {

	// fmt.Println("helpNeeded()")

	if noArgs := (len(os.Args) == 1); noArgs {
		return true
	}

	return false
}

func helpWanted() bool {

	// fmt.Println("helpWanted()")

	switch os.Args[1] {
	case "-h", "h", "help", "--help", "-H", "H", "HELP", "--HELP", "-help", "--h", "--H":
		return true
	}

	return false
}

func help(status int) {

	// fmt.Println("help()")
	// fmt.Printf("  status: %v\n", status)

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

	// fmt.Println("flags()")

	bools := []string{"-b", "-f"}

	flag.BoolVar(&PrintBase, "b", false, "")
	flag.BoolVar(&PrintFull, "f", false, "")
	flag.Parse()

	Paths = os.Args[1:]
	Paths = slices.Filter(Paths, bools...)
}

func validArgs() bool {

	// fmt.Println("validArgs()")

	if PrintBase && PrintFull {
		return false
	}

	if len(Paths) == 0 {
		return false
	}

	return true
}

func modPaths() error {

	// fmt.Println("modPaths()")

	if PrintBase {
		Paths = basepaths(Paths)
		return nil
	} else if PrintFull {
		var err error
		Paths, err = abspaths(Paths)
		return err
	} else {
		return nil
	}
}

func basepaths(paths []string) []string {

	// fmt.Println("basepaths()")
	// fmt.Printf("  paths: %v\n", paths)

	var based []string

	for _, fpath := range paths {
		basename := filepath.Base(fpath)
		based = append(based, basename)
	}

	return based
}

func abspaths(paths []string) ([]string, error) {

	// fmt.Println("abspaths()")
	// fmt.Printf("  paths: %v\n", paths)

	var absed []string

	for _, fpath := range paths {

		abs, err := filepath.Abs(fpath)
		if err != nil {
			return nil, err
		}

		absed = append(absed, abs)
	}

	return absed, nil
}

func splitExt(filename string) (name, ext string) {

	// fmt.Println("splitExt()")
	// fmt.Printf("  filename: %v\n", filename)

	name = filepath.Clean(filename)

	for {

		extTest := filepath.Ext(name)
		mimeType := mime.TypeByExtension(extTest)

		// base case
		if noMoreExts := (extTest == "" || mimeType == ""); noMoreExts {
			return
		}

		// recursive case
		ext = concat(extTest, ext)               // ".tar" + ".xz"
		name = strings.TrimSuffix(name, extTest) // "file.tar" > "file"
	}
}

func concat(slc ...string) string {

	// fmt.Println("concat()")
	// fmt.Printf("  slc: %v\n", slc)

	b := bytes.NewBuffer(nil)
	defer b.Reset()

	for _, s := range slc {
		b.WriteString(s)
	}

	return b.String()
}
