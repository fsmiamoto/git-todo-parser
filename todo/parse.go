package todo

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(f io.Reader) ([]Todo, error) {
	var result []Todo

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		cmd, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line: %w", err)
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

	for i := TodoCommand(Pick); i < Comment; i++ {
		if isCommand(i, fields[0]) {
			todo.Command = TodoCommand(i)
			fields = fields[1:]
			break
		}
	}

	if todo.Command == 0 {
		// unexpected command
		return todo, fmt.Errorf("unexpected command: %q", trimmed)
	}

	if todo.Command == Break {
		return todo, nil
	}

	if len(fields) == 0 {
		return todo, fmt.Errorf("missing commit id: %q", trimmed)
	}

	if todo.Command == Label || todo.Command == Reset {
		todo.Label = strings.Join(fields, " ")
		return todo, nil
	}

	if todo.Command == Exec {
		todo.ExecCommand = strings.Join(fields, " ")
		return todo, nil
	}

	todo.Commit = fields[0]

	return todo, nil
}

func isCommand(i TodoCommand, s string) bool {
	if i < 0 || i > Comment {
		return false
	}
	return len(s) > 0 &&
		(todoCommandInfo[i].cmd == s || todoCommandInfo[i].nickname == s)
}
