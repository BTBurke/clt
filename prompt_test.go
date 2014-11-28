package clt

import (
	"testing"
)

func stringBuildTester(b *StringBuilder, expected string, t *testing.T) {
	got := b.Finalize()
	if got != expected {
		t.Error("Got:\n%s\n-------------Expected:\n%s", got, expected)
	}
}

func TestSay(t *testing.T) {
	b := buildSay("This is a test")
	expected := "This is a test\n"
	stringBuildTester(b, expected, t)
}
