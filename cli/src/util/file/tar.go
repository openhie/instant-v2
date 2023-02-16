package file

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/luno/jettison/errors"
)

func TarSource(path string) (io.Reader, error) {
	dstInfo := archive.CopyInfo{
		Exists: true,
		IsDir:  true,
	}

	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	srcInfo, err := archive.CopyInfoSourcePath(path, false)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	_, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer preparedArchive.Close()

	return preparedArchive, nil
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
