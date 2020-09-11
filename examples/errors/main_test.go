package main_test

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		assert.Equal(t,
			[]string{
				"something went wrong",
				"",
				"main.crash3.func1:/main.go:32",
				"\tt.Reject(goerr.New(\"something went wrong\"))",
				"main.crash2.func1:/main.go:23",
				"\tt.Reject(e)",
				"main.crash1.func1:/main.go:12",
				"\tt.Reject(e)",
				"",
				"",
			},
			normaliseCmdOutput(out),
		)
	}
}

func normaliseCmdOutput(in []byte) []string {
	root := strings.ReplaceAll(runtime.GOROOT(), "\\", "/")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cwd = strings.ReplaceAll(cwd, "\\", "/")

	out := string(in)
	out = strings.ReplaceAll(out, "\r\n", "\n")
	out = strings.ReplaceAll(out, root, "")
	out = strings.ReplaceAll(out, cwd, "")

	return strings.Split(out, "\n")
}
