package file

import (
	"archive/tar"
	"bufio"
	"io/fs"
	"os"
	"testing"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

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

func Test_untarSource(t *testing.T) {
	type cases struct {
		source          string
		destination     string
		contentFileName string
	}

	testCases := []cases{
		// case: untar file into current directory
		{
			source:          "test_tar.tar",
			destination:     "test.txt",
			contentFileName: "test.txt",
		},
		// case: return error from not specifying source file
		{},
		// case: return error from not specifying destination file
		{
			source:          "test_tar.tar",
			contentFileName: "test.txt",
		},
		// case: untar file into nested directory,
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

			require.Equal(t, errors.New(expectedErr.Error()).Error(), err.Error())
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

func Test_tarSource(t *testing.T) {
	type cases struct {
		source          string
		contentFileName string
	}

	testCases := []cases{
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
			expectedErr := fs.PathError{
				Op:  "stat",
				Err: unix.ENOENT,
			}

			require.Equal(t, errors.New(expectedErr.Error()).Error(), err.Error())
		} else {
			tarFile, err := os.OpenFile(tc.source, os.O_CREATE|os.O_RDWR, 0777)
			jtest.RequireNil(t, err)
			defer tarFile.Close()

			_, err = tarFile.ReadFrom(re)
			jtest.RequireNil(t, err)

			untarFile := createTestFile(t, "untarred")
			err = UntarSource(tc.source, "untarred")
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
