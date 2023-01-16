package file

import (
	"archive/zip"
	"io/fs"
	"os"
	"testing"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestUnzipSource(t *testing.T) {
	type cases struct {
		source          string
		destination     string
		contentFileName string
	}

	testCases := []cases{
		// unzip file into current directory
		{
			source:          "test_zip.zip",
			destination:     "test.txt",
			contentFileName: "test.txt",
		},
		// return error from not specifying source file
		{},
		// unzip file into nested directory,
		{
			source:          "test_zip.zip",
			destination:     "./testDir/test.txt",
			contentFileName: "test.txt",
		},
	}
	for _, tc := range testCases {
		var zipFile *os.File
		if tc.source != "" {
			zipFile = createTestZipFile(t, tc.source, tc.contentFileName)
		}
		defer os.RemoveAll(tc.source)
		defer os.RemoveAll(tc.destination)
		defer os.RemoveAll("testDir")

		err := UnzipSource(tc.source, tc.destination)
		if err != nil {
			expectedErr := fs.PathError{
				Op:   "open",
				Path: tc.source,
				Err:  unix.ENOENT,
			}

			require.Equal(t, errors.New(expectedErr.Error()).Error(), err.Error())
		} else {
			_, err = os.Stat(tc.destination)
			jtest.RequireNil(t, err)
		}

		zipFile.Close()
		os.RemoveAll(tc.source)
		os.RemoveAll(tc.destination)
	}

	os.RemoveAll("testDir")
}

func createTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	jtest.RequireNil(t, err)

	return file
}

func createTestZipFile(t *testing.T, zipFileName string, contentFileName string) *os.File {
	zipFile := createTestFile(t, zipFileName)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	writer, err := zipWriter.Create(contentFileName)
	jtest.RequireNil(t, err)

	_, err = writer.Write([]byte("test data"))
	jtest.RequireNil(t, err)

	err = zipWriter.Flush()
	jtest.RequireNil(t, err)

	return zipFile
}
