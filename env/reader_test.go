package env

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	type testData struct {
		filePath string
		expectMap map[string]string
		isErr bool
		errString string
	}

	testDir := "./testData"

	testCases := map[string]testData{
		"default .env file": {
			filePath: filepath.Join(testDir, ".env"),
			expectMap: map[string]string{
				"name": "hatobus",
				"age": "24",
				"place": "Tokyo",
			},
		},
		"invalid .env format": {
			filePath: filepath.Join(testDir, ".env.invalid"),
			expectMap: nil,
			isErr: true,
			errString: "read data failed invalid data format",
		},
		"read from current dir": {
			filePath: "",
			expectMap: map[string]string{
				"name": "bus-hato",
				"age": "24",
				"description": "hello=world",
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T){
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
