package file

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/luno/jettison/errors"
)

func unzipFile(f *zip.File, destination string) error {
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return errors.Wrap(errors.New("invalid file path: "+filePath), "")
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return errors.Wrap(err, "")
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return errors.Wrap(err, "")
	}

	// Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func UnzipSource(source, destination string) error {
	// Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer reader.Close()

	// Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return errors.Wrap(err, "")
	}

	// Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}
