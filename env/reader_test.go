package env

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
)

func TestReadFromEnvFile(t *testing.T) {
	type testData struct {
		filePath  string
		expectMap map[string]string
		isErr     bool
		errString string
	}

	testDir := "./testData"

	testCases := map[string]testData{
		"default .env file": {
			filePath: filepath.Join(testDir, ".env"),
			expectMap: map[string]string{
				"name":  "hatobus",
				"age":   "24",
				"place": "Tokyo",
			},
		},
		"invalid .env format": {
			filePath:  filepath.Join(testDir, ".env.invalid"),
			expectMap: nil,
			isErr:     true,
			errString: "read data failed invalid data format",
		},
		"read from current dir": {
			filePath:  "",
			expectMap: nil,
			isErr:     true,
			errString: "no such file, check your input file name",
		},
		"read from has a empty line": {
			filePath: filepath.Join(testDir, ".env.hasline"),
			expectMap: map[string]string{
				"name":        "hato=bus",
				"age":         "24",
				"description": "this case has line",
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			values, err := ReadFromDotEnvFile(tc.filePath)
			if err != nil {
				if !tc.isErr {
					t.Fatalf("Unexpected error, error not required but got: %v", err)
				} else if err.Error() != tc.errString {
					t.Fatalf("Unexpected error, want %v but %v", tc.errString, err.Error())
				}
			}

			if diff := cmp.Diff(tc.expectMap, values); diff != "" {
				t.Fatalf("unexpected output, diff: %v", diff)
			}
		})
	}
}

func TestReadFromFileEncryptBase64(t *testing.T) {
	type testData struct {
		filePath  string
		expectOut string
		isErr     bool
		errString string
	}

	testDir := "./testData"

	testCases := map[string]testData{
		"default .env file": {
			filePath:  filepath.Join(testDir, ".env"),
			expectOut: "bmFtZT1oYXRvYnVzCmFnZT0yNApwbGFjZT1Ub2t5bwo=",
		},
		"file not found": {
			filePath:  filepath.Join(testDir, "file_not_found.txt"),
			isErr:     true,
			errString: "no such file, check your input file name",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			values, err := ReadFromFileEncryptBase64(tc.filePath)
			if err != nil {
				if !tc.isErr {
					t.Fatalf("Unexpected error, error not required but got: %v", err)
				} else if err.Error() != tc.errString {
					t.Fatalf("Unexpected error, want %v but %v", tc.errString, err.Error())
				}
			}

			if diff := cmp.Diff(tc.expectOut, values); diff != "" {
				t.Fatalf("unexpected output, diff: %v", diff)
			}
		})
	}
}
