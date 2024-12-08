package main

import (
	"os"
	"testing"
)

func TestMainGood(t *testing.T) {
	os.Args = []string{"main"}
	_ = os.Setenv("PASSWORD", "bogus")
	main()
}

func TestMainBad(t *testing.T) {
	os.Args = []string{"main"}
	main()
}
