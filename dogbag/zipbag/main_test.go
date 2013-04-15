package main

import (
	. "../terst"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	Terst(t)

	zipbag := __DogbagBagFunction()
	IsNot(zipbag, nil)

	archive, err := zipbag.Archive()
	Is(err, nil)
	Is(len(archive.File), 8)

	Is(zipbag.Path("abc."), "")

	file, err := zipbag.Open("abc.")
	Is(err, nil)

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	Is(err, nil)
	Is(hex.EncodeToString(hash.Sum(nil)), "9a46fcda0e43dc0cc3297621c21b9007fd5364ba")

}

func TestDeploy(t *testing.T) {
	Terst(t)

	zipbag := __DogbagBagFunction()
	IsNot(zipbag, nil)

	scratch, err := ioutil.TempDir("", "zipbag")
	Is(err, nil)
	if err != nil {
		defer os.RemoveAll(scratch)
	}

	Is(zipbag.Path("abc."), "")

	err = zipbag.Deploy(scratch)
	Is(err, nil)

	path := zipbag.Path("abc.")
	IsNot(path, "")

	_, err = os.Stat(path)
	Is(err, nil)

	{
		file, err := zipbag.Open("abc.")
		Is(err, nil)

		hash := sha1.New()
		_, err = io.Copy(hash, file)
		Is(err, nil)
		Is(hex.EncodeToString(hash.Sum(nil)), "9a46fcda0e43dc0cc3297621c21b9007fd5364ba")
	}

	err = ioutil.WriteFile(path, []byte("Hello, World.\n"), 0600)
	Is(err, nil)

	{
		file, err := zipbag.Open("abc.")
		Is(err, nil)

		hash := sha1.New()
		_, err = io.Copy(hash, file)
		Is(err, nil)
		Is(hex.EncodeToString(hash.Sum(nil)), "2bb1581f2a40bebe39c074a4d4a12b018f89f0a0")
	}

	path = zipbag.Path("def/def.")
	IsNot(path, "")

	_, err = os.Stat(path)
	Is(err, nil)

	err = zipbag.Deploy(scratch)
	Is(err, nil)

	{
		file, err := zipbag.Open("abc.")
		Is(err, nil)

		hash := sha1.New()
		_, err = io.Copy(hash, file)
		Is(err, nil)
		Is(hex.EncodeToString(hash.Sum(nil)), "9a46fcda0e43dc0cc3297621c21b9007fd5364ba")
	}

}

type Bag interface {
	Path(...string) string
	Open(...string) (io.ReadCloser, error)
}

func TestInterface(t *testing.T) {
	Terst(t)

	var zipbag Bag = __DogbagBagFunction()
	IsNot(zipbag, nil)

	Is(zipbag.Path("abc."), "")

	file, err := zipbag.Open("abc.")
	Is(err, nil)

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	Is(err, nil)
	Is(hex.EncodeToString(hash.Sum(nil)), "9a46fcda0e43dc0cc3297621c21b9007fd5364ba")
}
