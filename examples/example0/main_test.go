package main

import (
	"os"
	"testing"
)

func TestMainGood(t *testing.T) {
	os.Args = []string{"main"}
	main()
}
