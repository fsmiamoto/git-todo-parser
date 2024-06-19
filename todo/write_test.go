package todo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteWithDefaultCommentChar(t *testing.T) {
	testWrite(t, '#')
}

func TestWriteWithCustomCommentChar(t *testing.T) {
	testWrite(t, ';')
}

func testWrite(t *testing.T, commentChar byte) {
	tests := []struct {
		name     string
		todos    []Todo
		expected string
	}{
		{
			"todo1",
			[]Todo{
				{Command: Pick, Commit: "deadbeef", Msg: "My commit msg"},
				{Command: Pick, Commit: "beefdead", Msg: "Another awesome commit"},
				{Command: Reset, Label: "somecommit"},
				{Command: Comment, Comment: " comment"},
				{Command: Exec, ExecCommand: "cd subdir; make test"},
				{Command: Label, Label: "awesomecommit"},
				{Command: UpdateRef, Ref: "refs/heads/my-branch"},
				{Command: Merge, Commit: "6f5e4d", Flag: "-C", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
				{Command: Fixup, Commit: "abbaceef"},
				{Command: Break},
				{Command: NoOp},
			},
			fmt.Sprintf(`pick deadbeef My commit msg
pick beefdead Another awesome commit
reset somecommit
%c comment
exec cd subdir; make test
label awesomecommit
update-ref refs/heads/my-branch
merge -C 6f5e4d report-a-bug # Merge 'report-a-bug'
fixup abbaceef
break
noop
`, commentChar),
		},
		{
			"example from git website",
			[]Todo{
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
				{Command: Merge, Commit: "6f5e4d", Flag: "-c", Label: "report-a-bug", Msg: "Merge 'report-a-bug'"},
			},
			fmt.Sprintf(`label onto
%c Branch: refactor-button
reset onto
pick 123456 Extract a generic Button class from the DownloadButton one
pick 654321 Use the Button class for all buttons
label refactor-button
%c Branch: report-a-bug
reset refactor-button
pick abcdef Add the feedback button
label report-a-bug
reset onto
merge -C a1b2c3 refactor-button # Merge 'refactor-button'
merge -c 6f5e4d report-a-bug # Merge 'report-a-bug'
`, commentChar, commentChar),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			err := Write(f, tt.todos, commentChar)
			require.NoError(t, err)

			require.Equal(t, tt.expected, f.String())
		})
	}
}
