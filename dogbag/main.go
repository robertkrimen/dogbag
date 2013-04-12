// Command dogbag is way to bundle data with your Go executable.
//
// (Take your data to Go)
//
// To dogbag a file and use it:
//
//        # Create template.tmpl_dogbag.go from the file template.tmpl:
//        $ dogbag template.tmpl .
//        ...
//        // var template []byte
//        template := template_tmpl()
//
// To dogbag a directory and use it
//
//        # Create assets_dogbag.go from the directory ./assets:
//        $ dogbag ./assets .
//        ...
//        // dogbag is a Bag
//        dogbag, err := assets()
//
// Install
//
//      go get github.com/robertkrimen/dogbag/dogbag
//
// Usage
//  
//      dogbag <input> <output> [...]
//
//          <input> Can be a file, directory, or - (stdin)
//          If a file, then a filebag will be built
//          If a directory, then a zipbag will be built
//          If omitted, the default is stdin (and a filebag will be built)
//
//          <output> Can be a file, directory, or - (stdout)
//          If a directory, then an <input>.go file will be put in the directory
//          If omitted, the default is stdout
//      
//          -function=""      The name of the function returning the dogbag
//          -package=""       The package to put the dogbag .go file in
//          -empty=false      Make an empty (zip) dogbag (for development)
//          -fmt=true         Postprocess through gofmt
//          -usage=false      More help: Bag interface description, etc.
//
package main

import (
	"flag"
	"fmt"
	"github.com/robertkrimen/isatty"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func usageUsage() {
	fmt.Println(kilt.GraveTrim(`

    # To dogbag a directory:
    $ dogbag ./assets

    // assets.go:
    dogbag, err := assets()
    // dogbag is a Bag, with the following interface:

 SetBase(path string) error
    SetBase will set the base directory that the bag is to serve assets from.

    The bag will ignore what it has in memory and serve from disk instead.

    The format of the path argument is platform-specific (filepath).

 Deploy(path string) error
    Deploy will deploy the bag to disk, triggering SetBase in the process.

    The directory that the bag deploys to is guarded by the SHA1 digest of bag.
    That is, a different version of the same program deploying to the
    same path will not collide with a previous version.

    The format of the path argument is platform-specific (filepath).

 Path(path ...string) string
    Path will transform a bag path into a file path.

        path := bag.Path("tmpl", "main.tmpl")
        // path => ""

        bag.Deploy("/home/example")
        path := bag.Path("tmpl", "main.tmpl")
        // path => /home/example/bag.f7d9.../tmpl/main.tmpl

        path = bag.Path("data/image.png")
        // path => /home/example/bag.f7d9.../data/image.png

        path = bag.Path()
        // path => /home/example/bag.f7d9.../

    The path arguments for this method are always UNIX-style (forward-slash).
    The return value is a platform specific path (via filepath).

    This method is ONLY useful after either SetBase or Deploy has run.
    Otherwise, since the bag does not exist on disk, this method will always return ""
    (You can use this behavior to test if a bag has been deployed, if necessary)

 Open(path ...string) (io.ReadCloser, error)
    Open will open the asset at the given bag path for reading,
    whether the asset is in memory or on disk.

        file, err := bag.Open("tmpl/main.tmpl")
        if err != nil {
            return err
        }
        defer file.Close()
    `))
	os.Exit(0)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	if true {
		flag.PrintDefaults()
	}
	fmt.Fprintf(os.Stderr, kilt.GraveTrim(`

    # To dogbag a file:
    $ dogbag template.tmpl

    // template.tmpl.go:
    // var template []byte
    template := template_tmpl()

    ---

    # To dogbag a directory:
    $ dogbag ./assets 

    // assets.go:
    dogbag, err := assets()
    // dogbag is a Bag (See: dogbag --usage)

    // var file io.ReadCloser
    // Read from the zip in memory
    file, err := dogbag.Open("style.css")

    // Or deploy first
    err = dogbag.Deploy(".")

    // Read from the file on disk
    file, err := dogbag.Open("style.css")

    // var path string
    path := dogbag.Path("style.css")
    ... := dogbag.Path("index.js")

    `))
}

var (
	flag_empty    = flag.Bool("empty", false, "Make an empty (zip) dogbag (for development)")
	flag_input    = flag.String("input", "", "The input for the dogbag. Can either be a file or a directory")
	flag_output   = flag.String("output", "", "The name of the dogbag .go file to create/overwrite. Emit to stdout with -")
	flag_package  = flag.String("package", "", "The package to put the dogbag .go file in")
	flag_function = flag.String("function", "", "The name of the function returning the dogbag")
	flag_fmt      = flag.Bool("fmt", true, "Postprocess through gofmt")
	flag_usage    = flag.Bool("usage", false, "More help: Bag interface description, etc.")
)

var (
	global = struct {
		input         string
		inputPath     string
		output        string
		bagName       string
		bagIdentifier string
		bagPackage    string
		bagFunction   string
	}{}
)

func flagParse() {
	flag.Parse()

	if *flag_usage {
		usageUsage()
	}

	index := 0 // flag.Arg(...)

	global.input = *flag_input
	if global.input == "" {
		if !isatty.Check(os.Stdin.Fd()) {
			global.input = "-"
		} else {
			global.input = flag.Arg(index)
			if global.input == "" {
				if isatty.Check(os.Stdin.Fd()) {
					usage()
					os.Exit(2)
				}
				global.input = "-"
			} else {
				// Consume argument
				index += 1
			}
		}
	}

	if global.input != "-" {
		global.inputPath = global.input
		global.bagName = filepath.Base(global.inputPath)
		global.bagIdentifier = regexp.MustCompile(`[^\w+]`).ReplaceAllString(global.bagName, "_")
	}

	global.bagPackage = *flag_package
	if global.bagPackage == "" {
		pkg, err := build.ImportDir(".", 0)
		if err == nil {
			global.bagPackage = pkg.Name
			dbgf(`No package name specified: Using "%s"`, global.bagPackage)
		} else {
			dbgf(`No package name specified: Using "main"`)
			global.bagPackage = "main"
		}
	}

	global.bagFunction = *flag_function
	if global.bagFunction == "" {
		global.bagFunction = regexp.MustCompile(`^(\d)`).ReplaceAllString(global.bagIdentifier, "_$1")
		dbgf(`No function name specified: Using "%s"`, global.bagFunction)
	}

	global.output = *flag_output
	if global.output == "" {
		global.output = flag.Arg(index)
		if global.output == "" {
			global.output = "-"
		} else {
			// Consume argument
			index += 1
		}
	}
	if global.output != "-" {
		file, _ := os.Stat(global.output)
		if file != nil && file.IsDir() {
			// If we're here, bagName should ALWAYS != "", but just in case
			if global.bagName == "" {
				dbgf(`%/fatal//Cannot output file to: %s: already exists and is a directory`, global.output)
			} else {
				global.output = filepath.Join(global.output, global.bagName+".go")
				dbgf(`Using output file: %s`, global.output)
			}
		}
	}
}

func filebag() error {
	var bagInput io.Reader
	if global.input == "-" {
		bagInput = os.Stdin
		if global.bagFunction == "" {
			fmt.Fprintf(os.Stderr, "dogbag: File from stdin without -function\n")
			os.Exit(1)
		}
	} else {
		file, err := os.Open(global.input)
		if err != nil {
			return err
		}
		defer file.Close()
		bagInput = file
	}

	var bagOutput io.Writer
	if global.output == "-" {
		bagOutput = os.Stdout
	} else {
		file, err := os.Create(global.output)
		if err != nil {
			return err
		}
		defer file.Close()
		bagOutput = file
	}

	return fmtPipe(func(output io.Writer) error {
		return _filebag(bagInput, global.bagPackage, global.bagFunction, output)
	}, bagOutput)
}

func fmtPipe(input func(io.Writer) error, output io.Writer) error {

	inputOutput := output

	gofmt := false
	// First, see if gofmt is available
	if *flag_fmt {
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

func zipbag() error {

	bagInput := global.inputPath
	if *flag_empty {
		bagInput = ""
	}

	var bagOutput io.Writer
	if global.output == "-" {
		bagOutput = os.Stdout
	} else {
		file, err := os.Create(global.output)
		if err != nil {
			return err
		}
		bagOutput = file
	}

	return fmtPipe(func(output io.Writer) error {
		return _zipbag(bagInput, global.bagIdentifier, global.bagPackage, global.bagFunction, output)
	}, bagOutput)
}

func main() {
	flag.Usage = usage
	flagParse()

	err := func() error {

		if *flag_empty {
			return zipbag()
		}

		if global.inputPath != "" {
			stat, err := os.Stat(global.inputPath)
			if err != nil {
				return err
			}

			if stat.IsDir() {
				return zipbag()
			}
		}

		return filebag()
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "dogbag: error: %s\n", err)
		os.Exit(1)
	}
}
