package main

import (
	"fmt"
	"io"
)

var newlineByte = []byte{'\n'}

type byteWriter struct {
	io.Writer
	column int
}

func (self *byteWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var offset int
	for offset = range data {
		if self.column%12 == 0 {
			self.Writer.Write(newlineByte)
			self.column = 0
			fmt.Fprint(self.Writer, "\t")
		}

		fmt.Fprintf(self.Writer, "0x%02x, ", data[offset])
		self.column++
	}

	return offset + 1, nil
}

var newlineString = []byte("\"+\n\"")

type stringWriter struct {
	io.Writer
	column int
}

func (self *stringWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var offset int
	for offset = range data {
		if self.column%16 == 0 {
			self.Writer.Write(newlineString)
			self.column = 0
		}

		fmt.Fprintf(self.Writer, "\\x%02x", data[offset])
		self.column++
	}

	return offset + 1, nil
}

type base64Writer struct {
	io.Writer
	column int
}

func (self *base64Writer) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var offset int
	for offset = range data {
		if self.column%80 == 0 {
			self.Writer.Write(newlineByte)
			self.column = 0
		}

		self.Writer.Write(data[offset : offset+1])
		self.column++
	}

	return offset + 1, nil
}
