package main

import (
	. "./terst"
	"bytes"
	"path/filepath"
	"testing"
)

func Test_zipbag(t *testing.T) {
	Terst(t)

	var output bytes.Buffer
	err := _zipbag("", "testIdentifier", "testPackage", "testFunction", &output)
	Is(err, nil)
	Is(bytes.Contains(output.Bytes(), []byte("__DogbagBag")), false)
	Is(bytes.Contains(output.Bytes(), []byte("testIdentifier")), true)
	Is(bytes.Contains(output.Bytes(), []byte("testPackage")), true)
	Is(bytes.Contains(output.Bytes(), []byte("testFunction")), true)
	Compare(output.Len(), "<", 4096)

	err = _zipbag(filepath.FromSlash("test/assets"), "testIdentifier", "testPackage", "testFunction", &output)
	Is(err, nil)
	Compare(output.Len(), ">", 4096)
}
