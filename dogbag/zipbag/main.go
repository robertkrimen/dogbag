package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	dogbag := __DogbagBagFunction()

	err := func() error {
		fmt.Println(dogbag.digest())

		do := func() error {
			fmt.Println("---")
			fmt.Println(dogbag.Path(""))
			file, err := dogbag.Open("abc.")
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, file)
			return err
		}

		err := do()
		if err != nil {
			return err
		}

		err = dogbag.Deploy("")
		if err != nil {
			return err
		}

		return do()
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
