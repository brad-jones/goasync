package main_test

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
)

func TestSimple(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		c, err := actual.Filter(func(v string) bool { return v == "doing work" }).Count()
		assert.Nil(t, err)
		assert.Equal(t, 2, c)

		c, err = actual.Count()
		assert.Nil(t, err)
		assert.Contains(t, actual.At(c-2).String(), "END 1.0")
	}
}

func normaliseCmdOutput(in []byte) stream.Stream {
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

	return koazee.StreamOf(strings.Split(out, "\n"))
}
