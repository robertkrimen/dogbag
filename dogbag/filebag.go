package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

func filebagByteCompress(input io.Reader, output io.Writer, pkg, fn string) error {
	if pkg == "" {
		return errors.New("invalid package (\"\")")
	}
	if fn == "" {
		return errors.New("invalid function (\"\")")
	}

	fmt.Fprintf(output, kilt.GraveTrim(`
package %s

import (
    "bytes"
    "compress/gzip"
    "io"
)

// %s returns raw, uncompressed file data.
func %s() []byte {

    `), pkg, fn, fn)

	fmt.Fprintf(output, `deflate, err := gzip.NewReader(bytes.NewBuffer([]byte{`)
	{
		inflate := gzip.NewWriter(&byteWriter{Writer: output})
		io.Copy(inflate, input)
		inflate.Close()
	}
	fmt.Fprintf(output, `}))`)

	fmt.Fprint(output, kilt.GraveTrim(`

    if err != nil {
        panic("dogbag: Decompress failed: " + err.Error())
    }

    var data bytes.Buffer
    io.Copy(&data, deflate)
    deflate.Close()

    return data.Bytes()
}
    `))

	return nil
}

func _filebag(input io.Reader, bagPackage, bagFunction string, output io.Writer) error {
	bagInput := input
	bagOutput := output
	return filebagByteCompress(bagInput, bagOutput, bagPackage, bagFunction)
}
