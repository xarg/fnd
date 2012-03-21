package main

import (
	"testing"
)


func Test_unix2regexp(t *testing.T) {
	var patterns = map[string] string {
		"test*": "test.*",
		"test.*": "test\\..*",
		"*.png": ".*\\.png",
		"a?": "a.",
		"hello-": "hello\\-",
	}

	for pattern, expected := range patterns {
		result := unix2regexp(pattern)
		if result != expected {
			t.Errorf("Got %s expected %s", result, expected)
		}
	}
}
