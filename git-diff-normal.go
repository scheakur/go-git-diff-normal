package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(formatNormal(gitDiff(os.Args[1:])))
}

func gitDiff(args []string) string {
	base := strings.Fields("diff --no-index --no-color --no-ext-diff --unified=0")
	params := append(base, args...)
	out, err := exec.Command("git", params...).Output()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatal(err)
	}
	return string(out)
}

func formatNormal(unifiedDiff string) string {
	diff := unifiedDiff

	diff = regexp.MustCompile(`(?ms)\A.*?@@`).ReplaceAllString(diff, "@@")
	diff = regexp.MustCompile(`(?m)^\+`).ReplaceAllString(diff, "> ")
	diff = regexp.MustCompile(`(?m)^-`).ReplaceAllString(diff, "< ")
	diff = regexp.MustCompile(`(?m)^(<.*)(\r|\n|\r\n)(>)`).ReplaceAllString(diff, "$1$2---$2$3")

	re := regexp.MustCompile(`(?m)^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@.*`)
	return re.ReplaceAllStringFunc(diff, func(src string) string {
		m := re.FindStringSubmatch(src)
		before := startEnd(m[1], m[2])
		after := startEnd(m[3], m[4])
		action := actionType(m[2], m[4])
		return before + action + after
	})
}

func startEnd(start, size string) string {
	if size == "" || size == "0" {
		return start
	}
	end := num(start) + num(size) - 1
	return fmt.Sprintf("%s,%d", start, end)
}

func actionType(beforeSize, afterSize string) string {
	if beforeSize == "0" {
		return "a"
	}
	if afterSize == "0" {
		return "d"
	}
	return "c"
}

func num(str string) int64 {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return n
}
