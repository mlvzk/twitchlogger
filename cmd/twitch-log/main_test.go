package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadBulk(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("ch1\nch2\nch3"))

	channels, end := readBulk(scanner, 2)
	if len(channels) != 2 {
		t.Errorf("Invalid channel length: %d", len(channels))
	}
	if end {
		t.Errorf("Ended too quickly")
	}

	channels, end = readBulk(scanner, 2)
	if len(channels) != 1 {
		t.Errorf("Invalid channel length: %d", len(channels))
	}
	if !end {
		t.Errorf("Didn't end")
	}
}

func TestDeleteLastLine(t *testing.T) {
	bytes := []byte{'t', '\r', '\n', 's'}

	count := deleteLastLine(bytes, len(bytes))
	if count != 3 {
		t.Errorf("Wrong count, got: %d", count)
	}
}
