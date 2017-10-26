package clt

import (
	"bytes"
	"testing"
	"time"

	"github.com/BTBurke/snapshot"
)

func TestProgressSpinner(t *testing.T) {
	out := bytes.NewBuffer(nil)

	p := NewProgressSpinner("Testing a successful result")
	p.output = out
	p.Start()
	time.Sleep(1 * time.Second)
	p.Success()
	snapshot.Assert(t, out.Bytes())
}

func TestProgressBar(t *testing.T) {
	out := bytes.NewBuffer(nil)

	p := NewProgressBar("Testing a successful result")
	p.output = out
	p.Start()
	p.Update(0.5)
	p.Success()
	snapshot.Assert(t, out.Bytes())
}

func TestLoading(t *testing.T) {
	out := bytes.NewBuffer(nil)

	p := NewLoadingMessage("Testing a successful result", Dots, time.Duration(0))
	p.output = out
	p.Start()
	time.Sleep(1 * time.Second)
	p.Success()
	snapshot.Assert(t, out.Bytes())
}
