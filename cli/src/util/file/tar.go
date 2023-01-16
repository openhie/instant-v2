package file

import (
	"archive/tar"
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/luno/jettison/errors"
)

func TarSource(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "")
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.Wrap(err, "")
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return errors.Wrap(err, "")
		}

		f, err := os.Open(file)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return errors.Wrap(err, "")
		}

		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "")
		}

		return nil
	})

	if ok != nil {
		return nil, ok
	}

	err := tw.Close()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return bufio.NewReader(&buf), nil
}

func UntarSource(source, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "")
	}

	tr := tar.NewReader(file)
	if err != nil {
		return errors.Wrap(err, "")
	}

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return errors.Wrap(err, "")

		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(destination, 0755)
			if err != nil {
				return errors.Wrap(err, "")
			}

		case tar.TypeReg:
			dir, _ := filepath.Split(destination)
			if dir != "" {
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return errors.Wrap(err, "")
				}
			}

			f, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrap(err, "")
			}

			if _, err := io.Copy(f, tr); err != nil {
				return errors.Wrap(err, "")
			}

			f.Close()
		}
	}
}
