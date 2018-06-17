package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_global(t *testing.T) {
	tests := []struct {
		desc           string
		arg            string
		expectedStatus int
		expectedSubOut string
		expectedSubErr string
	}{
		{
			desc:           "undefined flag",
			arg:            "lstf --undefined",
			expectedStatus: exitCodeErr,
			expectedSubErr: "flag provided but not defined",
		},
		{
			desc:           "help",
			arg:            "lstf --help",
			expectedStatus: exitCodeErr,
			expectedSubErr: "Usage: lstf",
		},
		{
			desc:           "version",
			arg:            "lstf --version",
			expectedStatus: exitCodeOK,
			expectedSubErr: "lstf version",
		},
		{
			desc:           "credits",
			arg:            "lstf --credits",
			expectedStatus: exitCodeOK,
			expectedSubOut: "= lstf licensed under: =",
		},
		{
			desc:           "normal",
			arg:            "lstf",
			expectedStatus: exitCodeOK,
			expectedSubOut: "Local Address:Port",
		},
		{
			desc:           "--numeric",
			arg:            "lstf -n",
			expectedStatus: exitCodeOK,
			expectedSubOut: "Local Address:Port",
		},
		{
			desc:           "--json",
			arg:            "lstf --json",
			expectedStatus: exitCodeOK,
			expectedSubOut: "{\"direction\":",
		},
	}
	for _, tc := range tests {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split(tc.arg, " ")

		status := cli.Run(args)
		if status != tc.expectedStatus {
			t.Errorf("desc: %q, status should be %v, not %v", tc.desc, tc.expectedStatus, status)
		}

		if !strings.Contains(outStream.String(), tc.expectedSubOut) {
			t.Errorf("desc: %q, subout should contain %q, got %q", tc.desc, tc.expectedSubOut, outStream.String())
		}
		if !strings.Contains(errStream.String(), tc.expectedSubErr) {
			t.Errorf("desc: %q, subout should contain %q, got %q", tc.desc, tc.expectedSubErr, errStream.String())
		}
	}
}
