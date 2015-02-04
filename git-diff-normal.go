package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(string(formatNormal(gitDiff(os.Args[1:]))))
}

var (
	re1   = regexp.MustCompile(`(?ms)\A.*?@@`)
	repl1 = []byte("@@")
	re2   = regexp.MustCompile(`(?m)^\+`)
	repl2 = []byte("> ")
	re3   = regexp.MustCompile(`(?m)^-`)
	repl3 = []byte("< ")
	re4   = regexp.MustCompile(`(?m)^(<.*)(\r|\n|\r\n)(>)`)
	repl4 = []byte("$1$2---$2$3")
	re5   = regexp.MustCompile(`(?m)^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@.*`)

	add = []byte("a")
	del = []byte("d")
	mod = []byte("c")
)

func gitDiff(args []string) []byte {
	base := strings.Fields("diff --no-index --no-color --no-ext-diff --unified=0")
	params := append(base, args...)
	out, err := exec.Command("git", params...).Output()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatal(err)
	}
	return out
}

func formatNormal(unifiedDiff []byte) []byte {
	diff := unifiedDiff

	diff = re1.ReplaceAll(diff, repl1)
	diff = re2.ReplaceAll(diff, repl2)
	diff = re3.ReplaceAll(diff, repl3)
	diff = re4.ReplaceAll(diff, repl4)

	return re5.ReplaceAllFunc(diff, func(src []byte) []byte {
		m := re5.FindSubmatch(src)
		before := startEnd(m[1], m[2])
		after := startEnd(m[3], m[4])
		action := actionType(m[2], m[4])
		return append(append(before, action...), after...)
	})
}

func startEnd(start, size []byte) []byte {
	if len(size) == 0 || isZero(size) {
		return start
	}
	end := num(start) + num(size) - 1
	return []byte(fmt.Sprintf("%s,%d", start, end))
}

func actionType(beforeSize, afterSize []byte) []byte {
	if isZero(beforeSize) {
		return add
	}
	if isZero(afterSize) {
		return del
	}
	return mod
}

func isZero(b []byte) bool {
	return bytes.Equal(b, []byte("0"))
}

func num(str []byte) int64 {
	n, err := strconv.ParseInt(string(str), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return n
}
