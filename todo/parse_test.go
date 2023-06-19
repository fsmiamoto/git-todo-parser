package todo

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		commentChar byte
		expect      []Todo
		expectError error
	}{
		{name: "basic", inputPath: "./fixtures/todo1", commentChar: '#', expect: []Todo{
			{Command: Pick, Commit: "deadbeef", Msg: "My commit msg"},
			{Command: Pick, Commit: "beefdead", Msg: "Another awesome commit"},
			{Command: Reset, Label: "somecommit"},
			{Command: Reset, Label: "[new root]"},
			{Command: Reset, Label: "[new root]"},
			{Command: Comment, Comment: " comment"},
			{Command: Exec, ExecCommand: "cd subdir; make test"},
			{Command: Label, Label: "awesomecommit"},
			{Command: Label, Label: "theRef"},
			{Command: UpdateRef, Ref: "refs/heads/my-branch"},
			{Command: Merge, Commit: "6f5e4d", Flag: "-C", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
			{Command: Fixup, Commit: "abbaceef", Flag: "-C"},
			{Command: Break},
		}},
		{name: "missing exec cmd", inputPath: "./fixtures/missing_exec_cmd", commentChar: '#', expectError: ErrMissingExecCmd},
		{name: "missing label", inputPath: "./fixtures/missing_label", commentChar: '#', expectError: ErrMissingLabel},
		{name: "example from git website", inputPath: "./fixtures/git_example", commentChar: '#', expect: []Todo{
			{Command: Label, Label: "onto"},
			{Command: Comment, Comment: " Branch: refactor-button"},
			{Command: Reset, Label: "onto"},
			{Command: Pick, Commit: "123456", Msg: "Extract a generic Button class from the DownloadButton one"},
			{Command: Pick, Commit: "654321", Msg: "Use the Button class for all buttons"},
			{Command: Label, Label: "refactor-button"},
			{Command: Comment, Comment: " Branch: report-a-bug"},
			{Command: Reset, Label: "refactor-button"},
			{Command: Pick, Commit: "abcdef", Msg: "Add the feedback button"},
			{Command: Label, Label: "report-a-bug"},
			{Command: Reset, Label: "onto"},
			{Command: Merge, Commit: "a1b2c3", Flag: "-C", Label: "refactor-button", Msg: "Merge 'refactor-button'"},
			{Command: Merge, Commit: "6f5e4d", Flag: "-C", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
		}},
		{name: "custom comment char", inputPath: "./fixtures/custom_comment_char", commentChar: ';', expect: []Todo{
			{Command: Label, Label: "onto"},
			{Command: Comment, Comment: " Branch dev"},
			{Command: Reset, Label: "onto"},
			{Command: Pick, Commit: "086d35c", Msg: "one"},
			{Command: Label, Label: "dev"},
			{Command: Reset, Label: "onto"},
			{Command: Pick, Commit: "ad56c2e", Msg: "two"},
			{Command: Merge, Commit: "1c87252", Flag: "-C", Label: "dev", Msg: "Merge branch 'dev'"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputPath)
			require.NoError(t, err)

			defer f.Close()

			result, err := Parse(f, tt.commentChar)

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
