package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var newline = []byte{'\n'}

type ByteWriter struct {
	io.Writer
	column int
}

func (self *ByteWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var offset int
	for offset = range data {
		if self.column%12 == 0 {
			self.Writer.Write(newline)
			self.column = 0
		}

		fmt.Fprintf(self.Writer, "0x%02x,", data[offset])
		self.column++
	}

	return offset + 1, nil
}

var line = []byte("\"+\n\"")

type StringWriter struct {
	io.Writer
	column int
}

func (self *StringWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var offset int
	for offset = range data {
		if self.column%16 == 0 {
			self.Writer.Write(line)
			self.column = 0
		}

		fmt.Fprintf(self.Writer, "\\x%02x", data[offset])
		self.column++
	}

	return offset + 1, nil
}

func main() {
	fmt.Println("package main;")
	fmt.Print("\nvar zipbag_template = \"")
	tmpl, _ := ioutil.ReadAll(os.Stdin)
	tmpl = bytes.Replace(tmpl, []byte("package main\n"), []byte("package {{ .package }}\n"), 1)
	tmpl = bytes.Replace(tmpl, []byte("__DogbagBagType"), []byte("{{ .type }}"), -1)
	tmpl = bytes.Replace(tmpl, []byte("__DogbagBagFunction"), []byte("{{ .function }}"), -1)
	tmpl = bytes.Replace(tmpl, []byte("__DogbagBagData = ``"), []byte("{{ .type }}Data = string([]byte{ {{ .data }}\n})"), 1)
	tmpl = bytes.Replace(tmpl, []byte("__DogbagBagName"), []byte("\"{{ .name }}\""), -1)
	io.Copy(&StringWriter{Writer: os.Stdout}, bytes.NewReader(tmpl))
	fmt.Println("\"\n")
}
