package main

import (
	"io"
	"os"
	"os/exec"
)

func fmtPipe(input func(io.Writer) error, output io.Writer) error {

	inputOutput := output

	gofmt := false
	// First, see if gofmt is available
	if flag_fmt {
		err := exec.Command("gofmt").Run()
		if err == nil {
			gofmt = true
		}
	}

	if !gofmt {
		return input(output)
	}

	cmd := exec.Command("gofmt")
	cmdStdin, err := cmd.StdinPipe()
	if err == nil {
		cmd.Stderr = os.Stderr
		cmd.Stdout = output

		err = cmd.Start()
		if err == nil {
			inputOutput = cmdStdin
		} else {
			cmdStdin.Close()
			cmd = nil
		}
	}

	err = input(inputOutput)
	if cmd != nil {
		cmdStdin.Close()
		if err == nil {
			err = cmd.Wait()
		}
	}
	return err
}
