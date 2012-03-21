package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
//Note: you should clean the directory yourself
func createTestDir(levels int) string {
	defaultTmpDir := os.TempDir()
	targetDir := filepath.Join(defaultTmpDir, randString(6))
	if err := os.Mkdir(targetDir, 0777); err != nil {
		log.Fatal(err)
	}
	// recurisvly create the rest of the dirs
	createFiles(targetDir, 0, levels)
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

func TestFind(t *testing.T) {
	testDir := createTestDir(3)
	defer os.RemoveAll(testDir)

	// go about my bussiness
}
