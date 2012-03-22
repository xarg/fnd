package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const alphaNum = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	regexpFlag   = flag.String("e", "", "Use regexp")
	filetypeFlag = flag.String("type", "all", "Search only (f)iles, (d)irectories or (l)inks")
)

func showUsage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [pattern]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

//escape alpha numeric patterns
func regexpEscape(c byte) string {
	if bytes.IndexByte([]byte(alphaNum), c) == -1 {
		if c == '\000' {
			return "\\000"
		} else {
			return "\\" + string(c)
		}
	}
	return string(c)
}

// Translate *nix like pattern. Ex: *.py to .*\.py for regexp
func unixRegexp(pattern string) string {
	res := ""
	for _, c := range pattern {
		switch c {
		default:
			res = res + regexpEscape(byte(c))
		case '*':
			res = res + ".*"
		case '?':
			res = res + "."
		}
	}
	return res
}

func printFile(directory string, fileinfo os.FileInfo, stdout io.Writer) {
	filename := filepath.Join(directory, fileinfo.Name())

	if directory[0] != os.PathSeparator {
		filename = "." + string(os.PathSeparator) + filename
	}
	fmt.Fprintf(stdout, "%s\n", filename)
}

func parseDir(directory string, options map[string]string, stdout io.Writer) {
	if dir, err := os.Open(directory); err == nil {
		dirInfoSlice, _ := dir.Readdir(-1)
		for _, fileinfo := range dirInfoSlice {
			filename := fileinfo.Name()
			matched, _ := regexp.Match(options["pattern"], 
				[]byte(filename))
			if matched {
				printFile(directory, fileinfo, stdout)
			}
			if fileinfo.IsDir() {
				parseDir(filepath.Join(directory, filename),
					options, stdout)
			}
		}
	} else {
		log.Fatal(err)
	}
}

func Find(options map[string]string, stdout io.Writer) {
	parseDir(options["directory"], options, stdout)
}

func main() {
	options := make(map[string]string)

	flag.Usage = showUsage
	flag.Parse()

	options["pattern"] = ""
	options["directory"] = "."

	if flag.NArg() == 1 {
		if *regexpFlag != "" {
			options["pattern"] = *regexpFlag
			options["directory"] = flag.Arg(0)
		} else {
			options["pattern"] = unixRegexp(flag.Arg(0))
		}
	}

	if flag.NArg() == 2 {
		options["pattern"] = unixRegexp(flag.Arg(0))
		options["directory"] = flag.Arg(1)
	}

	Find(options, os.Stdout)
}
