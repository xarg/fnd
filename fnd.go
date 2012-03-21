package main

import (
	"bytes"
	"flag"
	"fmt"
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

func printFile(fileinfo os.FileInfo) {
	fmt.Printf("%s\n", fileinfo.Name())
}

func parseDir(directory, pattern string) {
	regexpPattern := ""
	if len(pattern) > 0 {
		regexpPattern = unixRegexp(pattern)
		fmt.Println(regexpPattern)
	} else if len(*regexpFlag) > 0 {
		regexpPattern = *regexpFlag
	}
	if dir, err := os.Open(directory); err == nil {
		dirInfoSlice, _ := dir.Readdir(1024)
		for _, fileinfo := range dirInfoSlice {
			filename := fileinfo.Name()
			matched, _ := regexp.Match(regexpPattern,
				[]byte(filename))
			if matched {
				printFile(fileinfo)
			}
		}
	} else {
		log.Fatal(err)
	}
}

func find() {
	pattern := ""
	directory := "."

	if flag.NArg() == 1 {
		pattern = flag.Arg(0)
	}

	if flag.NArg() == 2 {
		pattern = flag.Arg(0)
		directory = flag.Arg(1)
	}
	directory, _ = filepath.Abs(directory)
	parseDir(directory, pattern)
}

func main() {
	flag.Usage = showUsage
	flag.Parse()
	find()
}
