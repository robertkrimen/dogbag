/*
This is a stub command for generating: test_dogbag.go
This is the bare minimum necessary to test: zipbag/
*/
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	flag_fmt = true
)

func main() {
	err := func() error {
		data, err := _zipbagWalk(filepath.FromSlash("test/assets"))
		if err != nil {
			return err
		}

		return fmtPipe(func(output io.Writer) error {
			fmt.Fprintf(output, `
package main

var __DogbagBagName = "test"

var __DogbagBagTypeData = string([]byte{%s})
            `, data.String())
			return nil
		}, os.Stdout)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "stub: %s\n", err)
		os.Exit(1)
	}
}
