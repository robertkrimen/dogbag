package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	pth "path"
	"path/filepath"
	"strings"
)

func __DogbagBagFunction() *__DogbagBagType {
	return New__DogbagBagType(strings.NewReader(__DogbagBagTypeData))
}

type __DogbagBagType struct {
	data    __DogbagBagTypeDataReader
	name    string
	hash    string
	base    string
	archive *zip.Reader
}

type __DogbagBagTypeDataReader interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Seek(int64, int) (int64, error)
	Len() int
}

func New__DogbagBagType(data __DogbagBagTypeDataReader) *__DogbagBagType {
	return &__DogbagBagType{
		data: data,
		name: __DogbagBagName,
	}
}

func (self *__DogbagBagType) SetBase(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	self.base = path
	return nil
}

func (self *__DogbagBagType) Deploy(path string) error {
	err := self.SetBase(filepath.Join(path, self.name+"."+self.digest()))
	if err != nil {
		return err
	}

	archive, err := self.Archive()
	if err != nil {
		return err
	}

	for _, file := range archive.File {
		name := file.FileHeader.Name
		path := self.Path(name)

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		os.MkdirAll(filepath.Dir(path), 0777)

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self __DogbagBagType) Archive() (*zip.Reader, error) {
	if self.archive == nil {
		archive, err := zip.NewReader(self.data, int64(self.data.Len()))
		if err != nil {
			return nil, err
		}
		self.archive = archive
	}
	return self.archive, nil
}

func (self __DogbagBagType) Path(path ...string) string {
	if self.base == "" {
		return ""
	}
	{
		path := filepath.FromSlash(pth.Join(path...))
		path = filepath.Join(self.base, path)
		return path
	}
}

func (self __DogbagBagType) Open(path ...string) (io.ReadCloser, error) {
	if self.base == "" {
		archive, err := self.Archive()
		if err != nil {
			return nil, err
		}
		path := pth.Join(path...)
		for _, file := range archive.File {
			if file.Name == path {
				return file.Open()
			}
		}
		return nil, os.ErrNotExist
	}

	return os.Open(self.Path(path...))
}

func (self __DogbagBagType) digest() string {
	if self.hash == "" {
		hash := sha1.New()
		_, err := io.Copy(hash, self.data)
		self.data.Seek(0, 0)
		if err != nil {
			return ""
		}
		self.hash = hex.EncodeToString(hash.Sum(nil))
	}
	return self.hash
}
