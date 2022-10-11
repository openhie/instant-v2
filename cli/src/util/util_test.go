package util

import (
	"archive/zip"
	"io/fs"
	"os"
	"testing"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

func TestUnzipSource(t *testing.T) {
	type testArgs struct {
		source      string
		destination string
		contentFile string
	}

	testCases := []testArgs{
		// unzip file into current directory
		{
			source:      "test_zip.zip",
			destination: "test.txt",
			contentFile: "test.txt",
		},
		// return error from not specifying source file
		{},
		// unzip file into nested directory,
		{
			source:      "test_zip.zip",
			destination: "./testDir/test.txt",
			contentFile: "test.txt",
		},
	}
	for _, tc := range testCases {
		var zipFile *os.File
		if tc.source != "" {
			zipFile = createTestZipFile(tc.source, tc.contentFile)

		}
		defer os.RemoveAll(tc.source)
		defer os.RemoveAll(tc.destination)

		err := UnzipSource(tc.source, tc.destination)
		if err != nil {
			expectedErr := fs.PathError{
				Op:   "open",
				Path: tc.source,
				Err:  unix.ENOENT,
			}

			assert.Equal(t, errors.New(expectedErr.Error()).Error(), err.Error())
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

func createTestZipFile(zipFileName, contentFileName string) *os.File {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	writer, err := zipWriter.Create(contentFileName)
	if err != nil {
		panic(err)
	}

	_, err = writer.Write([]byte("test data"))
	if err != nil {
		panic(err)
	}

	err = zipWriter.Flush()
	if err != nil {
		panic(err)
	}

	return zipFile
}
