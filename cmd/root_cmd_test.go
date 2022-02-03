package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	is := is.New(t)

	err := rootCmd.Execute()

	is.NoErr(err)
}

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

func TestSubCmd(t *testing.T) {
	is := is.New(t)

	testcases := []struct {
		args      []string
		out_check bool
		err       error
		out       string
	}{
		{[]string{}, false, nil, ""},
		{[]string{"wrong"}, true, errors.New("unknown command \"wrong\" for \"tksadmin\""), ""},
		{[]string{"wrong", "cmd"}, true, errors.New("unknown command \"wrong\" for \"tksadmin\""), ""},

		{[]string{"completion"}, false, nil, ""},
		{[]string{"completion", "bash"}, false, nil, ""},
		{[]string{"completion", "fish"}, false, nil, ""},
		{[]string{"completion", "powershell"}, false, nil, ""},
		{[]string{"completion", "zsh"}, false, nil, ""},

		{
			args:      []string{"--config"},
			out_check: true,
			err:       errors.New("flag needs an argument: --config"),
			out:       "",
		},
		{[]string{"-t"}, false, nil, ""},

		{[]string{"contract"}, true, nil, ""},
		// {[]string{"contract", "create"}, false, nil, ""},	# It generates Exit(1) and the FAILed test
		{[]string{"contract", "create", "cli-unit-test"}, true, nil, ""},
	}

	for _, tc := range testcases {
		out, err := execute(t, rootCmd, tc.args...)

		is.Equal(tc.err, err)

		if tc.err == nil && tc.out_check {
			is.Equal(tc.out, out)
		}
	}
}
