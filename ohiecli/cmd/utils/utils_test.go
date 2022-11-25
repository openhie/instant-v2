package utils

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
)

func Test_sliceContains(t *testing.T) {
	testCases := []struct {
		slice   []string
		element string
		result  bool
		name    string
	}{
		{
			name:    "SliceContain test - should return true when slice contains element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Optimus Prime",
			result:  true,
		},
		{
			name:    "SliceContain test - should return false when slice does not contain element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Megatron",
			result:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ans := SliceContains(tt.slice, tt.element)

			if ans != tt.result {
				t.Fatal("SliceContains should return " + fmt.Sprintf("%t", tt.result) + " but returned " + fmt.Sprintf("%t", ans))
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_tarSource(t *testing.T) {
	type testArgs struct {
		source          string
		contentFileName string
	}

	testCases := []testArgs{
		// untar file into current directory
		{
			source:          "test_tar.tar",
			contentFileName: "test.txt",
		},
		// return error from not specifying source file
		{},
	}

	for _, tc := range testCases {
		defer os.Remove(tc.contentFileName)
		defer os.Remove(tc.source)
		defer os.Remove("untarred")

		if tc.contentFileName != "" {
			testFile := createTestFile(t, tc.contentFileName)
			defer testFile.Close()

			_, err := testFile.Write([]byte("test data"))
			jtest.RequireNil(t, err)
		}

		re, err := TarSource(tc.contentFileName)
		if err != nil {
			require.Equal(t, strings.Contains(err.Error(), "lstat : no such file or directory"), true)
		} else {
			tarFile, err := os.OpenFile(tc.source, os.O_CREATE|os.O_RDWR, 0777)
			jtest.RequireNil(t, err)
			defer tarFile.Close()

			_, err = tarFile.ReadFrom(re)
			jtest.RequireNil(t, err)

			untarFile := createTestFile(t, "untarred")
			err = untarSource(tc.source, "untarred")
			jtest.RequireNil(t, err)

			scanner := bufio.NewScanner(untarFile)
			scanner.Scan()
			if scanner.Text() != "test data" {
				t.FailNow()
			}
		}

		os.Remove(tc.contentFileName)
		os.Remove("untarred")
		os.Remove(tc.source)
	}
}

func createTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	jtest.RequireNil(t, err)

	return file
}

func untarSource(source, destination string) error {
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
