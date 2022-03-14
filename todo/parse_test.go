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

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		expect      []todo.Todo
		expectError error
	}{
		{name: "basic", inputPath: "./fixtures/todo1", expect: []todo.Todo{
			{Command: todo.Pick, Commit: "deadbeef", Msg: "My commit msg"},
			{Command: todo.Pick, Commit: "beefdead", Msg: "Another awesome commit"},
			{Command: todo.Reset, Label: "somecommit"},
			{Command: todo.Comment, Comment: " comment"},
			{Command: todo.Exec, ExecCommand: "cd subdir; make test"},
			{Command: todo.Label, Label: "awesomecommit"},
			{Command: todo.Merge, Commit: "6f5e4d", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
			{Command: todo.Fixup, Commit: "abbaceef"},
			{Command: todo.Break},
		}},
		{name: "missing exec cmd", inputPath: "./fixtures/missing_exec_cmd", expectError: todo.ErrMissingExecCmd},
		{name: "missing label", inputPath: "./fixtures/missing_label", expectError: todo.ErrMissingLabel},
		{name: "example from git website", inputPath: "./fixtures/git_example", expect: []todo.Todo{
			{Command: todo.Label, Label: "onto"},
			{Command: todo.Comment, Comment: " Branch: refactor-button"},
			{Command: todo.Reset, Label: "onto"},
			{Command: todo.Pick, Commit: "123456", Msg: "Extract a generic Button class from the DownloadButton one"},
			{Command: todo.Pick, Commit: "654321", Msg: "Use the Button class for all buttons"},
			{Command: todo.Label, Label: "refactor-button"},
			{Command: todo.Comment, Comment: " Branch: report-a-bug"},
			{Command: todo.Reset, Label: "refactor-button"},
			{Command: todo.Pick, Commit: "abcdef", Msg: "Add the feedback button"},
			{Command: todo.Label, Label: "report-a-bug"},
			{Command: todo.Reset, Label: "onto"},
			{Command: todo.Merge, Commit: "a1b2c3", Label: "refactor-button", Msg: "Merge 'refactor-button'"},
			{Command: todo.Merge, Commit: "6f5e4d", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputPath)
			defer f.Close()

			require.NoError(t, err)

			result, err := todo.Parse(f)

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
				return
			}

			require.NoError(t, err)

			if !reflect.DeepEqual(result, tt.expect) {
				t.Fatalf("Parse(%v) = %v; want %v", tt.inputPath, result, tt.expect)
			}
		})
	}
}
