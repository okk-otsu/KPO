package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		args      []string
		wantError bool
	}{
		{args: nil, wantError: false},
		{args: []string{"--demo"}, wantError: false},
		{args: []string{"--help"}, wantError: false},
		{args: []string{"--unknown"}, wantError: true},
		{args: []string{"extra"}, wantError: true},
	}

	for _, testCase := range cases {
		var output bytes.Buffer
		err := run(testCase.args, strings.NewReader("0\n"), &output)
		gotError := err != nil
		if gotError != testCase.wantError {
			t.Fatalf("аргументы %v: ожидалась ошибка=%v, получена ошибка=%v; err=%v",
				testCase.args, testCase.wantError, gotError, err)
		}
	}
}
