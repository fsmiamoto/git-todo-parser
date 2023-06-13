# git-todo-parser

[![Test](https://github.com/fsmiamoto/git-todo-parser/actions/workflows/test.yml/badge.svg)](https://github.com/fsmiamoto/git-todo-parser/actions/workflows/test.yml)

Small parser for git todo files.

## Grammar
```
p, pick <commit> = use commit
r, reword <commit> = use commit, but edit the commit message
e, edit <commit> = use commit, but stop for amending
s, squash <commit> = use commit, but meld into previous commit
f, fixup [-C | -c] <commit> = like "squash" but keep only the previous
                   commit's log message, unless -C is used, in which case
                   keep only this commit's message; -c is same as -C but
                   opens the editor
x, exec <command> = run command (the rest of the line) using shell
b, break = stop here (continue rebase later with 'git rebase --continue')
d, drop <commit> = remove commit
l, label <label> = label current HEAD with a name
t, reset <label> = reset HEAD to a label
m, merge [-C <commit> | -c <commit>] <label> [# <oneline>]
.       create a merge commit using the original merge commit's
.       message (or the oneline, if no original merge commit was
.       specified); use -c <commit> to reword the commit message
u, update-ref <ref> = track a placeholder for the <ref> to be updated
                      to this position in the new commits. The <ref> is
                      updated at the end of the rebase
```

## Example

```go
package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/fsmiamoto/git-todo-parser/todo"
)

func main() {
	reader := bytes.NewBufferString(`
pick 33bf560 Add third description heading
pick 979e6c4 Create link to blog site
edit b499fc4 Insert section to explain feature
exec cd src; make build
pick 5bd6691 Update numbered list to include more talking points
exec make test
    `)

	commentChar := byte('#')
	if output, err := exec.Command("git", "config", "--null", "core.commentChar").Output(); err == nil {
		if len(output) == 2 { // ends with the null character
			commentChar = output[0]
		}
	}

	todos, _ := todo.Parse(reader, commentChar)

	for _, t := range todos {
		if t.Command == todo.Pick {
			fmt.Println("Picking commit", t.Commit)
		} else if t.Command == todo.Exec {
			fmt.Println("Will exec ", t.ExecCommand)
		}
	}

	writer := &bytes.Buffer{}

	if err := todo.Write(writer, todos, commentChar); err == nil {
		fmt.Print("Original:\n", writer)
	}
}
```
