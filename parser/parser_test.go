package parser_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/fsmiamoto/git-rebase-todo-parser/parser"
	"github.com/stretchr/testify/require"
)

func readFixture(name string) []byte {
	bytes, _ := ioutil.ReadFile("./fixtures/" + name)
	return bytes
}

func TestParser(t *testing.T) {
	tests := []struct {
		name      string
		inputPath string
		expect    []parser.Todo
	}{
		{name: "basic", inputPath: "./fixtures/todo1", expect: []parser.Todo{
			{Command: parser.Pick},
			{Command: parser.Pick},
			{Command: parser.Exec},
			{Command: parser.Break},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputPath)
			defer f.Close()

			require.NoError(t, err)

			result, err := parser.Parse(f)

			require.NoError(t, err)

			if !reflect.DeepEqual(result, tt.expect) {
				t.Fatalf("Parser(%v) = %v; want %v", tt.inputPath, result, tt.expect)
			}
		})
	}
}
