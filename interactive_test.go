package clt

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/BTBurke/snapshot"
)

func WithInput(input string) (*InteractiveSession, *bytes.Buffer) {
	var out bytes.Buffer
	return &InteractiveSession{
		input:  bufio.NewReader(bytes.NewBufferString(input)),
		output: &out,
	}, &out
}

func TestSay(t *testing.T) {

	tt := []struct {
		Name       string
		Method     string
		Prompt     string
		Default    string
		ValHint    string
		Options    map[string]string
		Input      string
		Resp       string
		Validators []ValidationFunc
	}{
		{Name: "simple ask", Method: "ask", Prompt: "Did this work", Input: "yes\n", Resp: "yes"},
		{Name: "simple yn", Method: "yn", Prompt: "Do you want this to work", Input: "y", Default: "n", Resp: "y"},
		{Name: "retry yn", Method: "yn", Prompt: "Do you want this to work", Input: "q\ny", Default: "n", Resp: "y"},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			sess, buf := WithInput(tc.Input)

			var got string
			switch tc.Method {
			case "ask":
				got = sess.ask(tc.Prompt, tc.Default, tc.ValHint, tc.Validators...)
			case "yn":
				got = sess.AskYesNo(tc.Prompt, tc.Default)
			default:
				t.Errorf("unexpected test: %s", tc.Method)
			}
			if tc.Resp != got {
				t.Errorf("interactive session returned bad response, expected %s, got %s", tc.Resp, got)
			}
			snapshot.Assert(t, buf.Bytes())
		})

	}
}
