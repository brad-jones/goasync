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

func TestStream(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		assert.Equal(t, "foo did some work", actual.At(1).String())
		assert.Equal(t, "bar did some work", actual.At(2).String())

		c, err := actual.Count()
		assert.Nil(t, err)
		assert.Contains(t, actual.At(c-2).String(), "END 2.0")
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
