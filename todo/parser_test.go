package todo_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/fsmiamoto/git-todo-parser/todo"
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
		expect    []todo.Todo
	}{
		{name: "basic", inputPath: "./fixtures/todo1", expect: []todo.Todo{
			{Command: todo.Pick, Commit: "deadbeef"},
			{Command: todo.Pick, Commit: "beefdead"},
			{Command: todo.Comment},
			{Command: todo.Exec, ExecCommand: "cd subdir; make test"},
			{Command: todo.Break},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputPath)
			defer f.Close()

			require.NoError(t, err)

			result, err := todo.Parse(f)

			require.NoError(t, err)

			if !reflect.DeepEqual(result, tt.expect) {
				t.Fatalf("Parser(%v) = %v; want %v", tt.inputPath, result, tt.expect)
			}
		})
	}
}
