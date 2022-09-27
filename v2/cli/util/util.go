package util

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/cobra"
)

func Log(message string) {
	if os.Getenv("LOG") == "true" {
		fmt.Println(message)
	}
}

func LogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func PanicError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func GetFlagOrDefaultString(cmd *cobra.Command, flagName string) string {
	var name string
	if cmd.Flag(flagName).Changed {
		var err error
		name, err = cmd.Flags().GetString(flagName)
		LogError(err)
	} else {
		name = cmd.Flag(flagName).DefValue
	}
	return name

}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getPublicKeys(privateKeyFile string, password string) (*ssh.PublicKeys, error) {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		return nil, err
	}

	// Try to resolve password if a file path is provided as the password param
	file, err := os.Stat(password)
	if err == nil && !file.IsDir() {
		dat, err := os.ReadFile(password)
		if dat != nil && err == nil {
			password = string(dat)
		}
	}

	// Clone the given repository to the given directory
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		return nil, err
	}

	return publicKeys, nil
}

func CloneRepo(url, dest, sshKeyPath, sshPassword string) error {
	cloneOptions := &git.CloneOptions{
		URL: url,
	}
	publicKeys, err := getPublicKeys(sshKeyPath, sshPassword)
	if err == nil {
		cloneOptions.Auth = publicKeys
	}

	_, err = git.PlainClone(dest, false, cloneOptions)
	if err != nil {
		return err
	}
	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func UnzipSource(source, destination string) error {
	// Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
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

func UntarSource(source, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}
}

func TarSource(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})

	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
