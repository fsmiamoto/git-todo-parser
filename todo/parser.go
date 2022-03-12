package todo

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TodoCommand int

const (
	Pick TodoCommand = iota
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

var todoCommandInfo = [Comment + 1]struct {
	nickname string
	literal  string
}{
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
	{"", ""},
}

type Todo struct {
	Command     TodoCommand
	Commit      string
	ExecCommand string
	Label       string
	Msg         string
}

func Parse(f io.Reader) ([]Todo, error) {
	var result []Todo

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		cmd, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %q: %w", line, err)
		}

		result = append(result, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	return result, nil
}

func parseLine(l string) (Todo, error) {
	var todo Todo

	trimmed := strings.TrimSpace(l)

	if trimmed == "" || strings.HasPrefix(trimmed, CommentChar) {
		todo.Command = Comment
		return todo, nil
	}

	fields := strings.Fields(trimmed)

	for i := TodoCommand(0); i < Comment; i++ {
		if isCommand(i, fields[0]) {
			todo.Command = TodoCommand(i)
			fields = fields[1:]
			break
		}
	}

	if todo.Command == Break {
		return todo, nil
	}

	if todo.Command == Exec {
		todo.ExecCommand = strings.Join(fields, " ")
		return todo, nil
	}

	if len(fields) == 0 {
		return todo, fmt.Errorf("missing commit id: %s", trimmed)
	}

	todo.Commit = fields[0]

	return todo, nil
}

func isCommand(i TodoCommand, s string) bool {
	if i < 0 || i > Comment {
		return false
	}
	return len(s) > 0 &&
		(todoCommandInfo[i].literal == s || todoCommandInfo[i].nickname == s)
}
