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

func TestChained(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		hello := actual.Filter(func(v string) bool { return strings.HasPrefix(v, "hello:") })
		c, err := hello.Count()
		assert.Nil(t, err)
		assert.Equal(t, 3, c)

		v1 := strings.Split(hello.At(0).String(), ": ")[1]
		v2 := strings.Split(hello.At(1).String(), ": ")[1]
		v3 := strings.Split(hello.At(2).String(), ": ")[1]
		assert.NotEqual(t, v1, v3)
		assert.Equal(t, v2, v3)

		c, err = actual.Count()
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
