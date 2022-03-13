# git-todo-parser

[![Test](https://github.com/fsmiamoto/git-todo-parser/actions/workflows/test.yml/badge.svg)](https://github.com/fsmiamoto/git-todo-parser/actions/workflows/test.yml)

Small parser for git todo files.

WIP

## Example

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/fsmiamoto/git-todo-parser/todo"
)

func main() {
	buf := bytes.NewBufferString(`
pick 33bf560 Add third description heading
pick 979e6c4 Create link to blog site
edit b499fc4 Insert section to explain feature
exec cd src; make build
pick 5bd6691 Update numbered list to include more talking points
exec make test
    `)

	todos, _ := todo.Parse(buf)

	for _, t := range todos {
		if t.Command == todo.Pick {
			fmt.Println("Picking commit", t.Commit)
		} else if t.Command == todo.Exec {
			fmt.Println("Will exec ", t.ExecCommand)
		}
	}

}
```
