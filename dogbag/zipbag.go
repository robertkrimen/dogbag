package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	pth "path"
	"path/filepath"
	"text/template"
)

func _zipbag(input, bagIdentifier, bagPackage, bagFunction string, output io.Writer) error {
	tmpl := template.Must(template.New("").Parse(zipbag_template))
	archive := &bytes.Buffer{}
	zipWriter := zip.NewWriter(archive)

	if input != "" {
		err := walk("", input, zipWriter)
		if err != nil {
			return err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return err
	}

	var data bytes.Buffer
	encoder := &byteWriter{Writer: &data}
	_, err = io.Copy(encoder, archive)
	if err != nil {
		return err
	}

	var bagOutput io.Writer
	bagOutput = output

	{
		form := map[string]string{}
		form["package"] = bagPackage
		form["function"] = bagFunction
		form["type"] = "__Dogbag_" + bagIdentifier
		form["name"] = bagIdentifier
		form["data"] = data.String()

		tmpl.Execute(bagOutput, form)
	}

	return nil
}

func walk(base, startPath string, zipWriter *zip.Writer) error {
	return filepath.Walk(startPath, func(path string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(startPath, path)
		if err != nil {
			return err
		}

		filePath := path
		if stat.Mode()&os.ModeSymlink != 0 {
			filePath, err = filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			stat, err = os.Stat(filePath)
			if err != nil {
				return err
			}
			if stat.IsDir() {
				if err != nil {
					return err
				}
				return walk(pth.Join(base, ".", relativePath), filePath, zipWriter)
			}
		}

		name := pth.Join(base, filepath.ToSlash(relativePath))

		dbg("+", filePath, "=>", name)

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(name)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})
}
