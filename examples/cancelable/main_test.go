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

func TestCancelable(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		cancelableAsync := actual.Filter(func(v string) bool { return strings.HasPrefix(v, "cancelableAsync:") })
		c, err := cancelableAsync.Count()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, c, 5)
		assert.Contains(t, cancelableAsync.Last().String(), "I stopped cooperatively")

		chainedCancelableAsync := actual.Filter(func(v string) bool { return strings.HasPrefix(v, "chainedCancelableAsync:") })
		c, err = chainedCancelableAsync.Count()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, c, 5)
		assert.Contains(t, chainedCancelableAsync.Last().String(), "I stopped cooperatively")
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
