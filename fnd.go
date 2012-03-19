package main

import (
	"bytes"
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"regexp"
)

var (
	regexp_flag = flag.String("e", "", "Use regexp") 
	filetype_flag = flag.String("type", "all", "Search only (f)iles, (d)irectories or (l)inks") 
	alpha_num = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func show_usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [pattern]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

//escape alpha numeric patterns
func regexp_escape(c byte) string{
	if bytes.IndexByte([]byte(alpha_num), c) == -1 {
		if c == '\000'{
			return string("\\000")
		} else {
			return string('\\' + c)
		}
	}
	return string(c)
}

// Translate *nix like pattern. Ex: *.py to .*\.py for regexp
func unix2regexp(pat string) string{
	res := ""
	i := 0
	for i < len(pat) {
		c := pat[i]
		i = i + 1
		if c == '*' {
			res = res + ".*"
		} else if c == '?' {
			res = res + "."
		} else {
			res = res + regexp_escape(c)
		}
	}
	return res + "$"
}

func print_file(file_info os.FileInfo) {
	fmt.Printf("%s\n", file_info.Name())
}

func parse_dir(directory, pattern string){
	regexp_pattern := ""
	if len(pattern) > 0 {
		regexp_pattern = unix2regexp(pattern)
		fmt.Println(regexp_pattern)
	} else if len(*regexp_flag) > 0 {
		regexp_pattern = *regexp_flag
	}
	if directory_fp, err := os.Open(directory); err == nil {
		dir_info_slice, _ := directory_fp.Readdir(1024)
		for _, file_info := range dir_info_slice {
			file_name := file_info.Name()
			matched, _ := regexp.Match(regexp_pattern, 
						[]byte(file_name))
			if matched {
				print_file(file_info)
			}
		}
	} else {
		panic(fmt.Sprintf("Failed to open %s. Error: %s", directory, err))
	} 
}

func find(){
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
	parse_dir(directory, pattern)
}

func main(){
	flag.Usage = show_usage
	flag.Parse()
	find()
}
