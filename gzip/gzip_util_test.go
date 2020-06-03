package gzip

import (
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	f1, err := os.Open("../xx.txt")
	if err != nil {
		t.Fatal(err)
	}
	files := []*os.File{f1}
	Compress(files, "../xx.tar")

}
