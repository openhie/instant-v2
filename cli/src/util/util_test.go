package util

import (
	"archive/tar"
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

func Test_untarSource(t *testing.T) {
	type testArgs struct {
		source      string
		destination string
		contentFile string
	}

	testCases := []testArgs{
		// untar file into current directory
		{
			source:      "test_tar.tar",
			destination: "test.txt",
			contentFile: "test.txt",
		},
		// return error from not specifying source file
		{},
		// untar file into nested directory,
		{
			source:      "test_tar.tar",
			destination: "./testDir/test.txt",
			contentFile: "test.txt",
		},
	}

	for _, tc := range testCases {
		var tarFile *os.File
		if tc.source != "" {
			tarFile = createTestTarFile(t, tc.source, tc.contentFile)

		}
		// Ensure removal on panic
		defer os.RemoveAll(tc.source)
		defer os.RemoveAll(tc.destination)

		err := UntarSource(tc.source, tc.destination)
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

		tarFile.Close()

		// Ensure removal per test case
		os.RemoveAll(tc.source)
		os.RemoveAll(tc.destination)
	}

	os.RemoveAll("testDir")
}

func createTestTarFile(t *testing.T, tarFileName, contentFileName string) *os.File {
	tarFile, err := os.OpenFile(tarFileName, os.O_CREATE|os.O_RDWR, 0777)
	jtest.RequireNil(t, err)
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	fileData := []byte("test data")
	header := &tar.Header{
		Name: contentFileName,
		Size: int64(len(fileData)),
		Mode: 0777,
	}
	err = tarWriter.WriteHeader(header)
	jtest.RequireNil(t, err)

	_, err = tarWriter.Write(fileData)
	jtest.RequireNil(t, err)

	err = tarWriter.Flush()
	jtest.RequireNil(t, err)

	return tarFile
}
