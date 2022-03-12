package todo

type TodoCommand int

const (
	Pick TodoCommand = iota + 1
	Revert
	Edit
	Reword
	Fixup
	Squash

	Exec
	Break
	Label
	Reset
	Merge

	NoOp
	Drop

	Comment
)

const CommentChar = "#"

type Todo struct {
	Command     TodoCommand
	Commit      string
	ExecCommand string
	Label       string
	Msg         string
}

var todoCommandInfo = [Comment]struct {
	nickname string
	cmd      string
}{
	{"", ""}, // dummy value since we're using 1-based indexing
	{"p", "pick"},
	{"", "revert"},
	{"e", "edit"},
	{"r", "reword"},
	{"f", "fixup"},
	{"s", "squash"},
	{"x", "exec"},
	{"b", "break"},
	{"l", "label"},
	{"t", "reset"},
	{"m", "merge"},
	{"", "noop"},
	{"d", "drop"},
}
