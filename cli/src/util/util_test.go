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
		source          string
		destination     string
		contentFileName string
	}

	testCases := []testArgs{
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

func createTestZipFile(t *testing.T, zipFileName string, contentFileName string) *os.File {
	zipFile := createTestFile(t, zipFileName)
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
		source          string
		destination     string
		contentFileName string
	}

	testCases := []testArgs{
		// untar file into current directory
		{
			source:          "test_tar.tar",
			destination:     "test.txt",
			contentFileName: "test.txt",
		},
		// return error from not specifying source file
		{},
		// return error from not specifying destination file
		{
			source:          "test_tar.tar",
			contentFileName: "test.txt",
		},
		// untar file into nested directory,
		{
			source:          "test_tar.tar",
			destination:     "./testDir/test.txt",
			contentFileName: "test.txt",
		},
	}

	for _, tc := range testCases {
		// Ensure removal on panic
		defer os.RemoveAll(tc.source)
		defer os.RemoveAll(tc.destination)

		var contentFile, tarFile *os.File
		if tc.source != "" {
			contentFile = createTestFile(t, tc.contentFileName)
			tarFile = createTestTarFile(t, tc.source, contentFile)
		}

		err := UntarSource(tc.source, tc.destination)
		if err != nil {
			expectedErr := fs.PathError{
				Op:  "open",
				Err: unix.ENOENT,
			}

			if !assert.Equal(t, errors.New(expectedErr.Error()).Error(), err.Error()) {
				t.FailNow()
			}
		} else {
			_, err = os.Stat(tc.destination)
			jtest.RequireNil(t, err)
		}
		contentFile.Close()
		tarFile.Close()

		// Ensure removal per test case
		os.RemoveAll(tc.source)
		os.RemoveAll(tc.destination)
	}

	os.RemoveAll("testDir")
}

func createTestTarFile(t *testing.T, tarFileName string, contentFile *os.File) *os.File {
	tarFile := createTestFile(t, tarFileName)
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	fileData := []byte("test data")
	header := &tar.Header{
		Name: contentFile.Name(),
		Size: int64(len(fileData)),
		Mode: 0777,
	}
	err := tarWriter.WriteHeader(header)
	jtest.RequireNil(t, err)

	_, err = tarWriter.Write(fileData)
	jtest.RequireNil(t, err)

	err = tarWriter.Flush()
	jtest.RequireNil(t, err)

	return tarFile
}

func createTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	jtest.RequireNil(t, err)

	return file
}
