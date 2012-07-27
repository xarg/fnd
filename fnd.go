package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	alphaNum     = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	cacheDirname = "fnd.cache"
)

var (
	regexpFlag        = flag.String("e", "", "Use regexp")
	filetypeFlag      = flag.String("type", "all", "Search only (f)iles, (d)irectories or (l)inks. Comma separated.")
	caseSensitiveFlag = flag.Bool("s", false, "Case sensitive search")
	cacheFlag         = flag.Bool("c", false, "Use cache")
)

func showUsage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [pattern]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func readConfigFile(filePath string) map[string]string {
	var result map[string]string
	return result
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

func sha1Hex(str string) string {
	dirHash := sha1.New()
	io.WriteString(dirHash, str)
	return fmt.Sprintf("%x", dirHash.Sum(nil))
}

// print the filename, take format into consideration
func printFile(directory string, filename string, stdout io.Writer) {
	filename = filepath.Join(directory, filename)

	if directory[0] != os.PathSeparator {
		filename = "." + string(os.PathSeparator) + filename
	}
	fmt.Fprintf(stdout, "%s\n", filename)
}

// If we matched the pattern to a file, then print that filename
func printIfMached(options map[string]string, directory string, filename string, stdout io.Writer) {
	matched, _ := regexp.Match(options["pattern"], []byte(filename))
	if matched {
		printFile(directory, filename, stdout)
	}
}

func parseDir(options map[string]string, directory string, stdout io.Writer, cacheChan chan []byte) {
	dir, err := os.Open(directory)
	if err != nil { // can't open? just ignore it
		return
	}

	// this reads the whole content of the dir. It may not be a good idea.
	dirInfoSlice, _ := dir.Readdir(-1)
	for _, fileinfo := range dirInfoSlice {
		filename := fileinfo.Name()
		if options["caseSensitive"] == "false" {
			filename = strings.ToLower(filename)
		}
		ok := false
		if options["filetype"] != "all" {
			if options["filetype_f"] == "true" && !fileinfo.IsDir() && fileinfo.Mode() != os.ModeSymlink {
				ok = true
			}
			if options["filetype_d"] == "true" && fileinfo.IsDir() {
				ok = true
			}
			if options["filetype_l"] == "true" && fileinfo.Mode() == os.ModeSymlink {
				ok = true
			}
		} else {
			ok = true
		}
		if ok {
			filename := fileinfo.Name()
			cacheChan <- []byte(filename) // send to cache
			printIfMached(options, directory, filename, stdout)
		}
		if fileinfo.IsDir() {
			parseDir(options, filepath.Join(directory, filename),
				stdout, cacheChan)
		}
	}
}

// For a given dirname and a list of filenames write to the cache
func setCache(dirname string, cacheChan chan []byte) {
	cacheDir := filepath.Join(os.TempDir(), cacheDirname)
	// create if does not exist
	if fi, err := os.Stat(cacheDir); fi == nil && err != nil {
		if err := os.Mkdir(cacheDir, 0777); err != nil {
			log.Println(err) // just warn the user; no cache will be used.
			return
		}
	}
	cacheFilename := filepath.Join(cacheDir, sha1Hex(dirname))
	if fd, err := os.Create(cacheFilename); err == nil {
		defer fd.Close()
		if err != nil {
			log.Println(err) // just warn the user; no cache will be used.
			return
		}
		for filename := range cacheChan {
			fd.Write(append(filename, byte('\n')))
		}
	}
}

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// This will try to read the cache of a certain directory. 
// If we have it and it's old just delete the cache.
func printCache(options map[string]string, dirname string, stdout io.Writer) {
	cacheDir := filepath.Join(os.TempDir(), cacheDirname)
	cacheFilename := filepath.Join(cacheDir, sha1Hex(dirname))
	if lines, err := readLines(cacheFilename); err == nil {
		for _, filename := range lines {
			printIfMached(options, dirname, filename, stdout)
		}
	} else {
		log.Println(err) // just warn the user; no cache will be used.
	}
}

func Find(options map[string]string, stdout io.Writer) {
	if options["caseSensitive"] == "false" {
		options["pattern"] = strings.ToLower(options["pattern"])
	}
	options["filetype_f"] = "false"
	options["filetype_d"] = "false"
	options["filetype_l"] = "false"

	if options["filetype"] == "" { // just in case we don't have any option
		options["filetype"] = "all"
	}
	if options["filetype"] != "all" {
		for _, filetype := range strings.Split(options["filetype"], ",") {
			options["filetype_"+filetype] = "true"
		}
	}
	if options["cache"] == "true" {
		printCache(options, options["directory"], stdout)
		return // no need to parse the dir again
	}
	cacheChan := make(chan []byte, 1024)
	// Run a goroutine for writing the cache
	go setCache(options["directory"], cacheChan)
	parseDir(options, options["directory"], stdout, cacheChan)
}

func main() {
	options := make(map[string]string)

	flag.Usage = showUsage
	flag.Parse()

	options["pattern"] = ""
	options["directory"] = "."
	options["caseSensitive"] = "false"
	options["cache"] = "false"

	if *caseSensitiveFlag {
		options["caseSensitive"] = "true"
	}

	if *cacheFlag {
		options["cache"] = "true"
	}

	options["filetype"] = *filetypeFlag

	if flag.NArg() == 1 {
		if *regexpFlag != "" { // fnd -e <regexp> <dir>
			options["pattern"] = *regexpFlag
			options["directory"] = flag.Arg(0)
		} else { // fnd <pattern>
			options["pattern"] = unixRegexp(flag.Arg(0))
		}
	}

	if flag.NArg() == 2 { // fnd <pattern> <dir>
		if *regexpFlag != "" {
			log.Fatal("Can't use both regexp and pattern")
		}
		options["pattern"] = unixRegexp(flag.Arg(0))
		options["directory"] = flag.Arg(1)
	}

	Find(options, os.Stdout)
}
