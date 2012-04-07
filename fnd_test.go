package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func randString(length int) string {
	now := time.Now()
	rand.Seed(now.UnixNano())

	res := ""
	maxLen := len(alphaNum)
	for i := 0; i <= length; i++ {
		res = res + string(alphaNum[rand.Intn(maxLen)])
	}
	return res
}

//recursivelly create directories and files given `directory`
func createFiles(directory string, level, maxLevel int) {
	if level == maxLevel {
		return
	}

	for i := 0; i < rand.Intn(100); i++ {
		filename := filepath.Join(directory, randString(6))
		if rand.Intn(2) == 0 { //create a file
			fd, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			fd.Close()
		} else { //create a directory
			if err := os.Mkdir(filename, 0777); err != nil {
				log.Fatal(err)
			} else {
				createFiles(filename, level+1, maxLevel)
			}
		}
	}
}

//Create a random temporary directory with `levels` and return the absolute path
func createTestDirs(levels int) string {
	targetDir := createTempDir()
	createFiles(targetDir, 0, levels)
	return targetDir
}

// just create a temp dir and return the absolute path
// Note: you should clean the directory yourself
func createTempDir() string {
	defaultTmpDir := os.TempDir()
	targetDir := filepath.Join(defaultTmpDir, randString(6))
	if err := os.Mkdir(targetDir, 0777); err != nil {
		log.Fatal(err)
	}
	return targetDir
}

func TestUnixRegexp(t *testing.T) {
	var patterns = map[string]string{
		"test*":  "test.*",
		"test.*": "test\\..*",
		"*.png":  ".*\\.png",
		"a?":     "a.",
		"hello-": "hello\\-",
	}

	for pattern, expected := range patterns {
		result := unixRegexp(pattern)
		if result != expected {
			t.Errorf("Got %s expected %s", result, expected)
		}
	}
}

// try to create the directories
func TestFindSimple(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "hello"))
	os.Create(filepath.Join(dir, "world"))

	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   "hell",
		"directory": dir,
	}, outputBuf)

	expected := "hello"
	result := outputBuf.Bytes()
	if !bytes.Contains(result, []byte(expected)) {
		t.Errorf("Got %s expected %s", result, expected)
	}
}

func TestFindRandom(t *testing.T) {
	levels := 3
	testDir := createTestDirs(levels)
	defer os.RemoveAll(testDir)

	// check if we really have 3 levels.
	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   "",
		"directory": testDir,
	}, outputBuf)
	lines := strings.Split(outputBuf.String(), "\n")

	result := 0
	for _, line := range lines {
		slices := strings.Split(line, testDir)
		if len(slices) > 1 {
			parts := strings.Split(slices[1],
				string(os.PathSeparator))
			if len(parts) > result {
				result = len(parts)
			}
			if len(parts) == levels { // we found enough levels
				return
			}
		}
	}
	t.Errorf("Expected %d levels but only got %d", levels, result)
}

func TestFindUnixPatternQue(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "foo"))
	os.Create(filepath.Join(dir, "baar"))
	os.Create(filepath.Join(dir, "food"))

	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   unixRegexp("fo?d"),
		"directory": dir,
	}, outputBuf)

	expected := "food\n"
	parts := strings.Split(outputBuf.String(), string(os.PathSeparator))
	result := parts[len(parts)-1]
	if result != expected {
		t.Errorf("Got %s expected %s", result, expected)
	}
}

func TestFindUnixPatternStar(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "foo"))
	os.Create(filepath.Join(dir, "baar"))
	os.Create(filepath.Join(dir, "food"))

	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   unixRegexp("foo*"),
		"directory": dir,
	}, outputBuf)

	output := outputBuf.String()
	if !(strings.Contains(output, "foo\n") &&
		strings.Contains(output, "food\n")) {
		t.Errorf("Didn't get foo and food")
	}
}

func TestFindRegexp(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "faroo"))
	os.Create(filepath.Join(dir, "bar"))
	os.Create(filepath.Join(dir, "food"))

	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   "ar$",
		"directory": dir,
	}, outputBuf)

	if !strings.Contains(outputBuf.String(), "bar\n") {
		t.Errorf("Didn't get foo and food")
	}
}

func TestFindCaseSensitive(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "Foo"))
	os.Create(filepath.Join(dir, "foo"))

	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":       "Foo",
		"directory":     dir,
		"caseSensitive": "true",
	}, outputBuf)

	output := outputBuf.String()
	if !strings.Contains(output, "Foo") {
		t.Errorf("Didn't find 'Foo' in %s", output)
	}

	if strings.Contains(output, "foo") {
		t.Errorf("Shouldn't get 'foo' in %s", output)
	}
}

// Find just files and directories
func TestFindFileType(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	os.Create(filepath.Join(dir, "foo_file"))
	os.Mkdir(filepath.Join(dir, "foo_dir"), 0777)

	// Find just files and not dirs
	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   "foo",
		"directory": dir,
		"filetype":  "f",
	}, outputBuf)

	output := outputBuf.String()
	if !strings.Contains(output, "foo_file") {
		t.Errorf("Must contain 'foo_file' in %s", output)
	}

	if strings.Contains(output, "foo_dir") {
		t.Errorf("Shouldn't contain 'foo_dir' in %s", output)
	}

	// Find just dirs and not files
	outputBuf = bytes.NewBufferString("") // start with a new buffer
	Find(map[string]string{
		"pattern":   "foo",
		"directory": dir,
		"filetype":  "d",
	}, outputBuf)

	output = outputBuf.String()
	if !strings.Contains(output, "foo_dir") {
		t.Errorf("Must contain 'foo_dir' in %s", output)
	}

	if strings.Contains(output, "foo_file") {
		t.Errorf("Shouldn't contain 'foo_file' in %s", output)
	}

	// Find both files and directories
	outputBuf = bytes.NewBufferString("") // start with a new buffer
	Find(map[string]string{
		"pattern":   "foo",
		"directory": dir,
		"filetype":  "f,d",
	}, outputBuf)

	output = outputBuf.String()
	if !strings.Contains(output, "foo_dir") {
		t.Errorf("Must contain 'foo_dir' in %s", output)
	}

	if !strings.Contains(output, "foo_file") {
		t.Errorf("Must contain 'foo_file' in %s", output)
	}
}

func TestFindFileTypeSymLink(t *testing.T) {

}

func BenchmarkFind(b *testing.B) {
	b.StopTimer()
	fmt.Printf("Benchmarking with %d\n", b.N)

	if b.N > 10000 {
		return
	}
	dir := createTestDirs(b.N)
	b.StartTimer()
	outputBuf := bytes.NewBufferString("")
	Find(map[string]string{
		"pattern":   "",
		"directory": dir,
	}, outputBuf)
}
